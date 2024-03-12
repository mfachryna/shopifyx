package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type JWTToken struct {
	Token    string
	Claim    Claim
	ExpireAt time.Time
	Scheme   string
}

func SignedToken(claim Claim) (string, error) {
	exp := time.Now().Add(2 * time.Minute)
	expAt := exp.Unix()
	iat := time.Now().Unix()

	claim.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expAt,
		IssuedAt:  iat,
	}
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

type JWTError string

func (e JWTError) Error() string {
	return string(e)
}
