package repository

import "context"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	Insert(ctx context.Context, usr *UserRegister) error
	FindOneById(ctx context.Context, id int) (*User, error)
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOneCredentialByEmail(ctx context.Context, email string) (*User, error)
	UpdateTokenIdById(ctx context.Context, token string, id int) error
	DeleteTokenIdById(ctx context.Context, id int) error
}
