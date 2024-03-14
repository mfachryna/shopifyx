package domain

type Payments struct {
	ID                   string `json:"id"`
	BankAccountId        string `json:"bankAccountId" validate:"required"`
	PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,url"`
	Quantity             int64  `json:"quantity" validate:"required,min=1"`
	ProductId            string `json:"product_id"`
	UserId               string `json:"user_id"`
}
