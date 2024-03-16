package domain

type BankAccount struct {
	ID                string `json:"bankAccountId"`
	BankName          string `json:"bankName" validate:"required,min=5,max=15"`
	BankAccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
	BankAccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=60"`
}
