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

// @BasePath /api

// @Tags Order
// @Summary Get Orders with Paginate
// @Router /order [get]
// @Description Get with Paginated, if you just execute without query, by default it will show result page 1 and show 10 datas, you can also add query param for just input int as page will show.
// @Param page query int false "page select by page [page]"
// @Success 200 {array} helper.Response
// @Failure 400
func (db *order) GetWithPaginate(cx *gin.Context) {
	page, err := strconv.Atoi(cx.Query("page"))
	if err != nil {
		page = 1
	}
	var totalCount int
	errScan := db.conn.QueryRow(`SELECT COUNT(*) FROM tbl_order;`).Scan(&totalCount)
	if errScan != nil {
		response := helper.BuildErrorResponse("Fail when scan", errScan.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	perPage := 10
	offset := (page - 1) * perPage
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	var orders []model.OrderPaginate
	rows, err := db.conn.Query(`SELECT order_id, user_id, product_id FROM tbl_order LIMIT $1 OFFSET $2;`, perPage, offset)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when execute query", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	for rows.Next() {
		var order model.OrderPaginate
		err := rows.Scan(&order.OrderID, &order.UserID, &order.ProductID)
		if err != nil {
			response := helper.BuildErrorResponse("Fail when scan rows", err.Error(), helper.EmptyObject{})
			cx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		orders = append(orders, order)
	}
	pagination := model.Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	}
	responseSuccess := helper.BuildResponse(gin.H{
		"items":      orders,
		"pagination": pagination,
	})
	cx.JSON(http.StatusOK, responseSuccess)
}

// @Tags Order
// @Summary Get Detail of Order
// @Router /order/{id} [get]
// @Description Get Detail of order will show joined table by relational talbe data by input id of order_id.
// @Param id path int true "request id path"
// @Success 200 {array} helper.Response
// @Failure 400
func (db *order) GetDetail(cx *gin.Context) {
	idParam := cx.Param("id")
	var order model.OrderDetail
	sqlstatement := `SELECT o.order_id, o.quantity, o.price, o.shipping_address, o.status, o.created_at, u.user_id, u.username, u."password", u.email, u.first_name, u.last_name, u.is_admin, u.created_at, p.product_id, p."name", p.price, p.description, p.image_url, p.quantity, p.created_at FROM tbl_order o
	JOIN tbl_user u ON o.user_id = u.user_id
	JOIN tbl_product p ON o.product_id = p.product_id
	WHERE o.order_id = $1 LIMIT 1;`
	err := db.conn.QueryRow(sqlstatement, idParam).Scan(&order.OrderID, &order.Quantity, &order.Price, &order.ShippingAddress, &order.Status, &order.CreatedAt, &order.User.UserID, &order.User.Username, &order.User.Password, &order.User.Email, &order.User.FirstName, &order.User.LastName, &order.User.IsAdmin, &order.User.CreatedAt, &order.Product.ProductID, &order.Product.Name, &order.Product.Price, &order.Product.Description, &order.Product.ImageURL, &order.Product.Quantity, &order.Product.CreatedAt)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when execute query", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(order)
	cx.JSON(http.StatusOK, responseSuccess)
}

// @Tags Order
// @Summary Insert Order
// @Router /order [post]
// @Description Just regular Insert order.
// @Param request body model.OrderDto true "Payload Body [RAW]"
// @Success 200 {array} helper.Response
// @Failure 400
func (db *order) Insert(cx *gin.Context) {
	var order model.OrderDto
	err := cx.ShouldBind(&order)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when binding", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	resid := 0
	sqlstatement := `INSERT INTO tbl_order (user_id, product_id, quantity, price, shipping_address, status) VALUES ($1,$2,$3,$4,$5,$6) RETURNING order_id;`
	errExec := db.conn.QueryRow(sqlstatement, order.UserID, order.ProductID, order.Quantity, order.Price, order.ShippingAddress, order.Status).Scan(&resid)
	if errExec != nil {
		response := helper.BuildErrorResponse("Fail when execute query", errExec.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(resid)
	cx.JSON(http.StatusOK, responseSuccess)
}

// @Tags Order
// @Summary Insert Order
// @Router /order/{id} [put]
// @Description Just regular Update order, just change the value before execute, and you can check by get detail api.
// @Param id path int true "order_id param to be update"
// @Param request body model.OrderDto true "Payload Body [RAW]"
// @Success 200 {array} helper.Response
// @Failure 400
func (db *order) Update(cx *gin.Context) {
	var existingData model.OrderDto
	idParam := cx.Param("id")
	sqlstatement := `SELECT user_id, product_id, quantity, price, shipping_address, status FROM tbl_order WHERE order_id = $1 LIMIT 1;`
	err := db.conn.QueryRow(sqlstatement, idParam).Scan(&existingData.UserID, &existingData.ProductID, &existingData.Quantity, &existingData.Price, &existingData.ShippingAddress, &existingData.Status)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when get", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	var newData model.OrderDto
	errSB := cx.ShouldBind(&newData)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when binding", errSB.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	newData.UserID = helper.Ifelse(newData.UserID, existingData.UserID).(int)
	newData.ProductID = helper.Ifelse(newData.ProductID, existingData.ProductID).(int)
	newData.Quantity = helper.Ifelse(newData.Quantity, existingData.Quantity).(int)
	newData.Price = helper.Ifelse(newData.Price, existingData.Price).(float64)
	newData.ShippingAddress = helper.Ifelse(newData.ShippingAddress, existingData.ShippingAddress).(string)
	newData.Status = helper.Ifelse(newData.Status, existingData.Status).(string)
	sqlstatementupdate := `UPDATE tbl_order SET user_id = $1, product_id = $2, quantity = $3, price = $4, shipping_address = $5, status = $6 WHERE order_id = $7;`
	res, err := db.conn.Exec(sqlstatementupdate, newData.UserID, newData.ProductID, newData.Quantity, newData.Price, newData.ShippingAddress, newData.Status, idParam)
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

// @Tags Order
// @Summary Delete Order
// @Router /order/{id} [delete]
// @Description Just regular delete data by parsing id as param.
// @Param id path int true "request id path"
// @Success 200 {array} helper.Response
// @Failure 400
func (db *order) Delete(cx *gin.Context) {
	idParam := cx.Param("id")
	sqlstatementupdate := `DELETE FROM tbl_order WHERE order_id = $1;`
	res, err := db.conn.Exec(sqlstatementupdate, idParam)
	if err != nil {
		response := helper.BuildErrorResponse("Fail when exec db", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		response := helper.BuildErrorResponse("Fail when get rows affected", err.Error(), helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if affected == 0 {
		response := helper.BuildErrorResponse("Fail when delete", "no data affected", helper.EmptyObject{})
		cx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseSuccess := helper.BuildResponse(idParam)
	cx.JSON(http.StatusOK, responseSuccess)
}
