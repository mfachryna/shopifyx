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

type PaymentHandler struct {
	db       *sql.DB
	validate *validator.Validate
}

func NewPaymentHandler(db *sql.DB, validate *validator.Validate) *PaymentHandler {
	return &PaymentHandler{
		db:       db,
		validate: validate,
	}
}

func (ph *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Payments

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := ph.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fmt.Println(err.Error())
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	userId := r.Context().Value("user_id").(string)
	if userId == "" {
		fmt.Println("userId not found in context")
		response.Error(w, apierror.CustomServerError("userId not found in context"))
		return
	}
	uuid := uuid.New()

	productId := chi.URLParam(r, "productId")
	var count int
	var sellerId string
	if err := ph.db.QueryRow(`SELECT COUNT(products.id), products.user_id FROM products JOIN bank_accounts ON products.user_id = bank_accounts.user_id WHERE bank_accounts.id = $1 AND products.id = $2 GROUP BY products.user_id`, data.BankAccountId, productId).Scan(&count, &sellerId); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if !(count > 0) {
		fmt.Println("Bank Id dan Product ID tidak sesuai")
		response.Error(w, apierror.CustomError(http.StatusBadRequest, "Bank Id dan Product ID tidak sesuai"))
		return
	}

	if _, err := ph.db.Exec(
		`INSERT INTO payments (id,bank_account_id,payment_proof_image_url,product_id,quantity,user_id) VALUES ($1,$2,$3,$4,$5,$6)`,
		uuid, data.BankAccountId, data.PaymentProofImageUrl, productId, data.Quantity, userId,
	); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return

	}
	if _, err := ph.db.Exec(`UPDATE users SET product_sold_total = product_sold_total::int + $1 WHERE id = $2`, data.Quantity, sellerId); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}
	if _, err := ph.db.Exec(`UPDATE products SET purchase_count = purchase_count::int + $1 WHERE id = $2`, data.Quantity, productId); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	data.ID = uuid.String()
	data.ProductId = productId
	data.UserId = userId

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"Payment processed successfully",
		data,
	))
}
