package service

import (
	"errors"
	"os"
	"time"

	"github.com/Trickster-ID/dbo/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(username string, cx *gin.Context) error
	ValidateToken(token string) (string, error)
	Logout(cx *gin.Context)
}

type jwtService struct {
	secretKey string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: os.Getenv("jwt-secret-key"),
	}
}

func (js *jwtService) GenerateToken(username string, cx *gin.Context) error {
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(js.secretKey))
	if err != nil {
		return err
	}
	cx.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	return nil
}

func (js *jwtService) ValidateToken(tknStr string) (string, error) {
	claims := &model.Claims{}
	tok, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if tok.Valid {
		return claims.Username, nil
	}
	return "", errors.New("token not valid")
}

func (js *jwtService) Logout(cx *gin.Context) {
	cx.SetCookie("token", "", -1, "/", "localhost", false, true)
}
