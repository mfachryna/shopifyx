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
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductHandler struct {
	db       *sql.DB
	validate *validator.Validate
}

func NewProductHandler(db *sql.DB, validate *validator.Validate) *ProductHandler {
	return &ProductHandler{
		db:       db,
		validate: validate,
	}
}

func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Product

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := ph.validate.Struct(data); err != nil {
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

	if _, err := ph.db.Exec(
		`INSERT INTO products (id,name,price,image_url,stock,condition,is_purchasable,tags,user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		uuid, data.Name, data.Price, data.ImageUrl, data.Stock, data.Condition, data.IsPurchasable, pq.Array(data.Tags), userId,
	); err != nil {
		fmt.Println(err)
		response.Error(w, apierror.CustomServerError("failed to insert data"))
		return
	}

	data.ID = uuid.String()

	response.Success(w, apisuccess.CustomResponse(
		http.StatusCreated,
		"product added successfully",
		data,
	))
}
