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

func (ph *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		data domain.Product
		id   string
	)

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
	productId := chi.URLParam(r, "productId")
	if productId == "" {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := ph.db.QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		response.Error(w, apierror.ClientForbidden())
		return
	}

	_, err = ph.db.Exec(
		`UPDATE products SET name = $1, price = $2, image_url = $3, stock = $4, condition = $5, tags = $6, is_purchasable = $7 WHERE id = $8`,
		data.Name, data.Price, data.ImageUrl, data.Stock, data.Condition, pq.Array(data.Tags), data.IsPurchasable, productId,
	)
	if err != nil {
		response.Error(w, apierror.CustomServerError("failed to update product"))
		return
	}

	data.ID = productId

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"product updated successfully",
		data,
	))
}

func (ph *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var id string
	userId := r.Context().Value("user_id").(string)

	productId := chi.URLParam(r, "productId")
	if productId == "" {
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := ph.db.QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		response.Error(w, apierror.ClientForbidden())
		return
	}

	_, err = ph.db.Exec(`DELETE FROM products WHERE id = $1`, productId)
	if err != nil {
		response.Error(w, apierror.CustomServerError("failed to delete product"))
		return
	}

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"product deleted successfully",
		struct {
			ID string `json:"id"`
		}{
			ID: productId,
		},
	))
}
