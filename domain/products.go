package domain

type Product struct {
	ID            string   `json:"id"`
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         *int64   `json:"price" validate:"required,min=0"`
	ImageUrl      string   `json:"imageUrl" validate:"required,url"`
	Stock         *int64   `json:"stock" validate:"required,numeric,min=0"`
	Condition     string   `json:"condition" validate:"required,eq=new|eq=second"`
	Tags          []string `json:"tags" validate:"required,min=0,dive,min=0"`
	IsPurchasable bool     `json:"isPurchasable" validate:"isBool"`
}
type ProductData struct {
	ID            string   `json:"productId"`
	Name          string   `json:"name"`
	Price         *int64   `json:"price"`
	ImageUrl      string   `json:"imageUrl"`
	Stock         *int64   `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount int64    `json:"purchaseCount"`
}

type ProductFilter struct {
	UserOnly       bool     `json:"userOnly" schema:"userOnly"`
	Limit          *int64   `json:"limit" validate:"required,numeric" schema:"limit"`
	Offset         *int64   `json:"offset" validate:"required,numeric" schema:"offset"`
	Tags           []string `json:"tags" validate:"min=0,dive,min=0" schema:"tags"`
	Condition      string   `json:"condition" validate:"omitempty,eq=new|eq=second" schema:"condition"`
	ShowEmptyStock bool     `json:"showEmptyStock" schema:"showEmptyStock"`
	MaxPrice       *int64   `json:"maxPrice" validate:"required_with=MinPrice,omitempty,numeric,min=0" schema:"maxPrice"`
	MinPrice       *int64   `json:"minPrice" validate:"required_with=MaxPrice,omitempty,numeric,min=0" schema:"minPrice"`
	SortBy         string   `json:"sortBy" validate:"omitempty,eq=new|eq=second" schema:"sortBy"`
	OrderBy        string   `json:"orderBy" validate:"omitempty,eq=asc|eq=desc" schema:"orderBy"`
	Search         string   `json:"search" validate:"omitempty,min=3" schema:"search"`
}
