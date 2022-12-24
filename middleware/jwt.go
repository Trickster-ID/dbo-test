package middleware

import (
	"net/http"

	"github.com/Trickster-ID/dbo/helper"
	"github.com/Trickster-ID/dbo/service"
	"github.com/gin-gonic/gin"
)

// func GetJwtToken(creds model.Credentials, w http.ResponseWriter) error {
// 	expirationTime := time.Now().Add(20 * time.Minute)
// 	claims := &model.Claims{
// 		Username: creds.Username,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		return err
// 	}
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Value:   tokenString,
// 		Expires: expirationTime,
// 	})
// 	return nil
// }

func AuthJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := c.Cookie("token")
		if err != nil {
			response := helper.BuildErrorResponse("Fail when get cookie", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		username, err := jwtService.ValidateToken(tkn)
		if err != nil {
			response := helper.BuildErrorResponse("Fail when validate token", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		if username == "" {
			response := helper.BuildErrorResponse("Token not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
	}
}
