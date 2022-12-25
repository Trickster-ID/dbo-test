package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Username string `json:"username" example:"admin1"`
	Password string `json:"password" example:"password1"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
