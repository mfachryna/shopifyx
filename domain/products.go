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
