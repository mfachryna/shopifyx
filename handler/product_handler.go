package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Croazt/shopifyx/domain"
	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	apisuccess "github.com/Croazt/shopifyx/utils/response/success"
	"github.com/Croazt/shopifyx/utils/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
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

func (ph *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.ServerError())
		return
	}

	var filter domain.ProductFilter
	if err := schema.NewDecoder().Decode(&filter, r.Form); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := ph.validate.Struct(filter); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fmt.Println(err.Error())
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}
	*filter.Offset = *filter.Limit * (*filter.Offset)
	sql, sqlTotal, err := getFilteredSql(r, filter)
	if err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomError(http.StatusForbidden, err.Error()))
		return
	}

	var count int64
	if err := ph.db.QueryRow(sqlTotal).Scan(&count); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}
	rows, err := ph.db.Query(sql)
	if err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	data := make([]domain.ProductData, 0)
	for rows.Next() {
		var product domain.ProductData
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, pq.Array(&product.Tags), &product.IsPurchasable, &product.PurchaseCount)
		if err != nil {
			fmt.Println(err.Error())
			response.Error(w, apierror.CustomServerError("Error scanning row:"+err.Error()))
			return
		}
		data = append(data, product)
	}
	rows.Close()

	response.SuccessMeta(w, apisuccess.IndexResponse(
		http.StatusOK,
		"ok",
		data,
		domain.Meta{
			Limit:  *filter.Limit,
			Offset: *filter.Offset,
			Total:  count,
		},
	))
}

func (ph *ProductHandler) Show(w http.ResponseWriter, r *http.Request) {
	var (
		productData domain.ProductDetail
		sellerId    string
	)

	productId := chi.URLParam(r, "productId")
	if productId == "" {
		fmt.Println("Product id not found")
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := ph.db.QueryRow(
		"SELECT id, name, price, image_url, stock, condition, tags, is_purchasable, purchase_count, user_id  FROM products WHERE products.id = $1",
		productId).
		Scan(&productData.Product.ID, &productData.Product.Name, &productData.Product.Price, &productData.Product.ImageUrl, &productData.Product.Stock, &productData.Product.Condition, pq.Array(&productData.Product.Tags), &productData.Product.IsPurchasable, &productData.Product.PurchaseCount, &sellerId); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if err := ph.db.QueryRow("SELECT name, product_sold_total FROM users WHERE id = $1", sellerId).Scan(&productData.Seller.Name, &productData.Seller.ProductSoldTotal); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	rows, err := ph.db.Query("SELECT id, bank_name, bank_account_name, bank_account_number FROM bank_accounts WHERE user_id = $1", sellerId)

	if err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}
	defer rows.Close()

	productData.Seller.BankAccounts = make([]domain.BankAccount, 0)
	for rows.Next() {
		var bankAccount domain.BankAccount
		err := rows.Scan(&bankAccount.ID, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber)
		if err != nil {
			fmt.Println(err.Error())
			response.Error(w, apierror.CustomServerError("Error scanning row:"+err.Error()))
			return
		}
		productData.Seller.BankAccounts = append(productData.Seller.BankAccounts, bankAccount)
	}

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
		"ok",
		productData,
	))
}

func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data domain.Product

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err.Error())
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	if err := ph.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fmt.Println(e)
			response.Error(w, apierror.CustomError(http.StatusBadRequest, validation.CustomError(e)))
			return
		}
	}

	userId := r.Context().Value("user_id").(string)
	if userId == "" {
		fmt.Println("userId not found in context")
		response.Error(w, apierror.CustomError(http.StatusForbidden, "userId not found in context"))
		return
	}
	uuid := uuid.New()

	if _, err := ph.db.Exec(
		`INSERT INTO products (id,name,price,image_url,stock,condition,is_purchasable,tags,user_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		uuid, data.Name, data.Price, data.ImageUrl, data.Stock, data.Condition, data.IsPurchasable, pq.Array(data.Tags), userId,
	); err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError("failed to insert data"))
		return
	}

	data.ID = uuid.String()

	response.Success(w, apisuccess.CustomResponse(
		http.StatusOK,
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
	productId := chi.URLParam(r, "productId")
	if productId == "" {
		fmt.Println("userId not found in context")
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := ph.db.QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		err := apierror.ClientForbidden()
		fmt.Println(err.Message)
		response.Error(w, err)
		return
	}

	_, err = ph.db.Exec(
		`UPDATE products SET name = $1, price = $2, image_url = $3, stock = $4, condition = $5, tags = $6, is_purchasable = $7 WHERE id = $8`,
		data.Name, data.Price, data.ImageUrl, data.Stock, data.Condition, pq.Array(data.Tags), data.IsPurchasable, productId,
	)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println("ProductID not found in context")
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := ph.db.QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		err := apierror.ClientForbidden()
		fmt.Println(err.Message)
		response.Error(w, err)
		return
	}

	_, err = ph.db.Exec(`DELETE FROM products WHERE id = $1`, productId)
	if err != nil {
		fmt.Println(err.Error())
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

func (ph *ProductHandler) Stock(w http.ResponseWriter, r *http.Request) {
	var (
		id    string
		stock int64
	)
	userId := r.Context().Value("user_id").(string)

	productId := chi.URLParam(r, "productId")
	if productId == "" {
		fmt.Println("ProductID not found in context")
		response.Error(w, apierror.ClientBadRequest())
		return
	}

	err := ph.db.QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	if id != userId {
		err := apierror.ClientForbidden()
		fmt.Println(err.Message)
		response.Error(w, err)
		return
	}

	err = ph.db.QueryRow("SELECT stock FROM products WHERE id = $1", productId).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			response.Error(w, apierror.ClientNotFound("product"))
			return
		}

		fmt.Println(err.Error())
		response.Error(w, apierror.CustomServerError(err.Error()))
		return
	}

	type Stock struct {
		Stock int64 `json:"stock"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Stock{Stock: stock})
}

func getFilteredSql(r *http.Request, filter domain.ProductFilter) (string, string, error) {
	sort := "id"
	if !(filter.SortBy == "") {
		sort = " ORDER BY " + filter.SortBy
	}

	order := "asc"
	if !(filter.OrderBy == "") {
		order = filter.OrderBy
	}
	where := ""
	if filter.UserOnly {
		userId := r.Context().Value("user_id")
		if userId == nil {
			return "", "", fmt.Errorf("userOnly filter can be used if you logged in")
		}
		where = fmt.Sprintf("user_id = '%s'", userId)
	}

	if len(filter.Tags) > 0 {
		jsonTag, err := json.Marshal([]string(filter.Tags))
		if err == nil {
			replacer := strings.NewReplacer("[", "{", "]", "}")
			stringTag := replacer.Replace(string(jsonTag))
			if where != "" {
				where += fmt.Sprintf(" AND tags && '%s'", stringTag)
			} else {
				where = fmt.Sprintf("tags && '%s'", stringTag)
			}
		}
	}

	if filter.Condition != "" {
		if where != "" {
			where += fmt.Sprintf(" AND condition = '%s'", filter.Condition)
		} else {
			where = fmt.Sprintf("condition = '%s'", filter.Condition)
		}
	}

	if *filter.MaxPrice > -1 || *filter.MinPrice > -1 {
		if where != "" {
			where += fmt.Sprintf(" AND price >= '%d' AND price <= '%d'", *filter.MinPrice, *filter.MaxPrice)
		} else {
			where = fmt.Sprintf("price >= '%d' AND price <= '%d'", *filter.MinPrice, *filter.MaxPrice)
		}
	}
	if filter.Search != "" {
		if where != "" {
			where += " AND name LIKE '%" + filter.Search + "%'"
		} else {
			where = "name LIKE '%" + filter.Search + "%'"
		}
	}

	sql := fmt.Sprintf("SELECT id,name,price,image_url,stock,condition,tags,is_purchasable,purchase_count FROM products WHERE %s ORDER BY %s %s LIMIT %d OFFSET %d", where, sort, order, *filter.Limit, *filter.Offset)
	sqlTotal := fmt.Sprintf("SELECT count(id) FROM products WHERE %s", where)
	return sql, sqlTotal, nil
}
