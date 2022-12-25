package controller

import (
	"database/sql"
	"net/http"

	"github.com/Trickster-ID/dbo/helper"
	"github.com/Trickster-ID/dbo/model"
	"github.com/Trickster-ID/dbo/service"
	"github.com/gin-gonic/gin"
)

type Authentikasi interface {
	InsertLoginData(cx *gin.Context)
	GetLoginData(cx *gin.Context)
}

type authentikasi struct {
	con    *sql.DB
	jwtsvc service.JWTService
}

func NewAuthentikasi(db *sql.DB, jwtSvc service.JWTService) Authentikasi {
	return &authentikasi{
		con:    db,
		jwtsvc: jwtSvc,
	}
}

// @BasePath /api

// Get Detail Login Data
// @Tags AUTH
// @Summary Post Login Data
// @Router /auth [post]
// @Description Post to get jwt token that save in cookies
// @Param request body model.Credentials true "Payload Body [RAW]"
// @Success 200 {array} helper.Response
// @Failure 400
func (auth *authentikasi) InsertLoginData(cx *gin.Context) {
	var creds model.Credentials
	err := cx.ShouldBind(&creds)
	if err != nil {
		auth.jwtsvc.Logout(cx)
		response := helper.BuildErrorResponse("Fail", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	userid := 0
	sqlStatement := `SELECT user_id FROM tbl_user WHERE username = $1 AND "password" = $2 LIMIT 1;`
	errExec := auth.con.QueryRow(sqlStatement, creds.Username, creds.Password).Scan(&userid)
	if errExec != nil {
		auth.jwtsvc.Logout(cx)
		response := helper.BuildErrorResponse("Fail when execute database", errExec.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	if userid == 0 {
		auth.jwtsvc.Logout(cx)
		response := helper.BuildErrorResponse("Fail", "You are not authorized", helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	errJwt := auth.jwtsvc.GenerateToken(creds.Username, cx)
	if errJwt != nil {
		auth.jwtsvc.Logout(cx)
		response := helper.BuildErrorResponse("Fail", errJwt.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	responseSuccess := helper.BuildResponse("")
	cx.JSON(http.StatusOK, responseSuccess)
}

// Get Detail Login Data
// @Tags AUTH
// @Summary Get Detail Login Data
// @Router /auth [get]
// @Description Get the detail of user by token
// @Success 200 {array} helper.Response
// @Failure 400
func (auth *authentikasi) GetLoginData(cx *gin.Context) {
	token, err := cx.Cookie("token")
	if err != nil {
		response := helper.BuildErrorResponse("Fail when get cookie", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	username, err := auth.jwtsvc.ValidateToken(token)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when validate token", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	sqlStatement := `SELECT * FROM tbl_user WHERE username = $1 LIMIT 1;`
	rows, err := auth.con.Query(sqlStatement, username)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when execute db", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	var r model.User
	for rows.Next() {
		err := rows.Scan(&r.UserID, &r.Username, &r.Password, &r.Email, &r.FirstName, &r.LastName, &r.IsAdmin, &r.CreatedAt)
		if err != nil {
			response := helper.BuildErrorResponse("Fail when scan rows", err.Error(), helper.EmptyObject{})
			cx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}
	responseSuccess := helper.BuildResponse(r)
	cx.JSON(http.StatusOK, responseSuccess)
}
