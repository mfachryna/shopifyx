package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Croazt/shopifyx/domain"
	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/Croazt/shopifyx/utils/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BankAccountHandler struct {
	db       *sql.DB
	validate *validator.Validate
}

func NewBankAccountHandler(db *sql.DB, validate *validator.Validate) *BankAccountHandler {
	return &BankAccountHandler{
		db:       db,
		validate: validate,
	}
}

func (bah *BankAccountHandler) Index(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	if userId == "" {
		response.Error(w, apierror.CustomServerError("userId not found in context"))
		return
	}

	rows, err := bah.db.Query(`SELECT id,bank_name,bank_account_name, bank_account_number FROM bank_accounts WHERE user_id = $1`, userId)
	if err != nil {
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}
	defer rows.Close()

	data := make([]domain.BankAccount, 0)
	for rows.Next() {
		var bankAccount domain.BankAccount
		err := rows.Scan(&bankAccount.ID, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber)
		if err != nil {
			response.Error(w, apierror.CustomServerError("Error scanning row:"+err.Error()))
			continue
		}
		data = append(data, bankAccount)
	}

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"success",
		data,
	))
}

func (bah *BankAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.BankAccount

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := bah.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	userId := r.Context().Value("user_id").(string)
	if userId == "" {
		response.Error(w, apierror.CustomServerError("userId not found in context"))
		return
	}
	uuid := uuid.New()

	if _, err := bah.db.Exec(
		`INSERT INTO bank_accounts (id,bank_name,bank_account_name,bank_account_number,user_id) VALUES ($1,$2,$3,$4,$5)`,
		uuid, data.BankName, data.BankAccountName, data.BankAccountNumber, userId,
	); err != nil {
		fmt.Println(err)
		response.Error(w, apierror.CustomServerError("failed to insert data"))
		return
	}

	data.ID = uuid.String()

	response.Success(w, apisuccess.CustomResponse(
		http.StatusCreated,
		"Bank Account added successfully",
		data,
	))
}

func (bah *BankAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		data domain.BankAccount
		id   string
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := bah.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	userId := r.Context().Value("user_id").(string)
	bankAccountId := chi.URLParam(r, "bankAccountId")
	if bankAccountId == "" {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := bah.db.QueryRow("SELECT user_id FROM bank_accounts WHERE id = $1", bankAccountId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, apierror.ClientNotFound("Bank Account"))
			return
		}

		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		response.Error(w, apierror.ClientForbidden())
		return
	}

	_, err = bah.db.Exec(
		`UPDATE bank_accounts SET bank_name = $1, bank_account_name = $2, bank_account_number = $3 WHERE id = $4`,
		data.BankName, data.BankAccountName, data.BankAccountNumber, bankAccountId,
	)
	if err != nil {
		response.Error(w, apierror.CustomServerError("failed to update Bank Account"))
		return
	}

	data.ID = bankAccountId

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"Bank Account updated successfully",
		data,
	))
}

func (bah *BankAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var id string
	userId := r.Context().Value("user_id").(string)

	bankAccountId := chi.URLParam(r, "bankAccountId")
	if bankAccountId == "" {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := bah.db.QueryRow("SELECT user_id FROM bank_accounts WHERE id = $1", bankAccountId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, apierror.ClientNotFound("Bank Account"))
			return
		}

		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		response.Error(w, apierror.ClientForbidden())
		return
	}

	_, err = bah.db.Exec(`DELETE FROM bank_accounts WHERE id = $1`, bankAccountId)
	if err != nil {
		response.Error(w, apierror.CustomServerError("failed to delete Bank Account"))
		return
	}

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"Bank Account deleted successfully",
		struct {
			ID string `json:"id"`
		}{
			ID: bankAccountId,
		},
	))
}
