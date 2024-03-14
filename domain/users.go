package domain

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UserAuthResponse struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type UserSellerData struct {
	Name             string        `json:"name"`
	ProductSoldTotal string        `json:"productSoldTotal"`
	BankAccounts     []BankAccount `json:"bankAccounts"`
}
