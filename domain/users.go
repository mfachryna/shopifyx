package domain

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserAuthResponse struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}
