package controller

import (
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/Trickster-ID/dbo/helper"
	"github.com/Trickster-ID/dbo/model"
	"github.com/gin-gonic/gin"
)

type Customer interface {
	GetWithPaginate(cx *gin.Context)
	GetDetail(cx *gin.Context)
	Insert(cx *gin.Context)
	Update(cx *gin.Context)
	Delete(cx *gin.Context)
}

type customer struct {
	conn *sql.DB
}

func NewCustomer(db *sql.DB) Customer {
	return &customer{
		conn: db,
	}
}

func (db *customer) GetWithPaginate(cx *gin.Context) {
	page, err := strconv.Atoi(cx.Query("page"))
	if err != nil {
		page = 1
	}
	var totalCount int
	errScan := db.conn.QueryRow(`SELECT COUNT(*) FROM tbl_user;`).Scan(&totalCount)
	if errScan != nil {
		response := helper.BuildErrorResponse("Fail when scan", errScan.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	perPage := 10
	offset := (page - 1) * perPage
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	var users []model.UserPaginate
	rows, err := db.conn.Query(`SELECT user_id, username, first_name || ' ' || last_name as fullname FROM tbl_user LIMIT $1 OFFSET $2;`, perPage, offset)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when execute query", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	for rows.Next() {
		var user model.UserPaginate
		err := rows.Scan(&user.UserID, &user.Username, &user.Fullname)
		if err != nil {
			response := helper.BuildErrorResponse("Fail when scan rows", err.Error(), helper.EmptyObject{})
			cx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		users = append(users, user)
	}
	pagination := model.Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	}
	responseSuccess := helper.BuildResponse(gin.H{
		"items":      users,
		"pagination": pagination,
	})
	cx.JSON(http.StatusOK, responseSuccess)
}
func (db *customer) GetDetail(cx *gin.Context) {
	idParam := cx.Param("id")
	var user model.User
	sqlstatement := `SELECT user_id, username, "password", email, first_name, last_name, is_admin, created_at FROM tbl_user WHERE user_id = $1 LIMIT 1;`
	err := db.conn.QueryRow(sqlstatement, idParam).Scan(&user.UserID, &user.Username, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when execute query", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(user)
	cx.JSON(http.StatusOK, responseSuccess)
}
func (db *customer) Insert(cx *gin.Context) {
	var user model.User
	err := cx.ShouldBind(&user)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when binding", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	resid := 0
	sqlstatement := `INSERT INTO tbl_user (username, "password", email, first_name, last_name, is_admin) VALUES ($1,$2,$3,$4,$5,$6) RETURNING user_id;`
	errExec := db.conn.QueryRow(sqlstatement, user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.IsAdmin).Scan(&resid)
	if errExec != nil {
		response := helper.BuildErrorResponse("Fail when execute query", errExec.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(resid)
	cx.JSON(http.StatusOK, responseSuccess)
}
func (db *customer) Update(cx *gin.Context) {
	var existingData model.User
	idParam := cx.Param("id")
	sqlstatement := `SELECT username, "password", email, first_name, last_name, is_admin FROM tbl_user WHERE user_id = $1 LIMIT 1;`
	err := db.conn.QueryRow(sqlstatement, idParam).Scan(&existingData.Username, &existingData.Password, &existingData.Email, &existingData.FirstName, &existingData.LastName, &existingData.IsAdmin)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when get", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	var newData model.User
	errSB := cx.ShouldBind(&newData)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when binding", errSB.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	newData.Username = helper.Ifelse(newData.Username, existingData.Username).(string)
	newData.Password = helper.Ifelse(newData.Password, existingData.Password).(string)
	newData.Email = helper.Ifelse(newData.Email, existingData.Email).(string)
	newData.FirstName = helper.Ifelse(newData.FirstName, existingData.FirstName).(string)
	newData.LastName = helper.Ifelse(newData.LastName, existingData.LastName).(string)
	if newData.IsAdmin != existingData.IsAdmin {
		newData.IsAdmin = existingData.IsAdmin
	}
	sqlstatementupdate := `UPDATE tbl_user SET username = $1, "password" = $2, email = $3, first_name = $4, last_name = $5, is_admin = $6 WHERE user_id = $7;`
	res, err := db.conn.Exec(sqlstatementupdate, newData.Username, newData.Password, newData.Email, newData.FirstName, newData.LastName, newData.IsAdmin, idParam)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when update", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		response := helper.BuildErrorResponse("Fail when update", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if affected == 0 {
		response := helper.BuildErrorResponse("Fail when update", "no data affected", helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(idParam)
	cx.JSON(http.StatusOK, responseSuccess)
}
func (db *customer) Delete(cx *gin.Context) {
	idParam := cx.Param("id")
	sqlstatementupdate := `DELETE FROM tbl_user WHERE user_id = $1;`
	res, err := db.conn.Exec(sqlstatementupdate, idParam)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when update", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		response := helper.BuildErrorResponse("Fail when update", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if affected == 0 {
		response := helper.BuildErrorResponse("Fail when update", "no data affected", helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(idParam)
	cx.JSON(http.StatusOK, responseSuccess)
}
