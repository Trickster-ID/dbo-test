package controller

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Order interface {
	GetWithPaginate(cx *gin.Context)
	GetDetail(cx *gin.Context)
	Insert(cx *gin.Context)
	Update(cx *gin.Context)
	Delete(cx *gin.Context)
}

type order struct {
	conn *sql.DB
}

func NewOrder(db *sql.DB) Order {
	return &order{
		conn: db,
	}
}

func (db *order) GetWithPaginate(cx *gin.Context) {

}
func (db *order) GetDetail(cx *gin.Context) {

}
func (db *order) Insert(cx *gin.Context) {

}
func (db *order) Update(cx *gin.Context) {

}
func (db *order) Delete(cx *gin.Context) {

}
