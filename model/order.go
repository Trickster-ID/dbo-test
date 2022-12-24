package model

import "time"

type Order struct {
	OrderID         int       `db:"order_id" json:"order_id" form:"order_id"`
	UserID          int       `db:"user_id" json:"user_id" form:"user_id"`
	ProductID       int       `db:"product_id" json:"product_id" form:"product_id"`
	Quantity        int       `db:"quantity" json:"quantity" form:"quantity"`
	Price           float64   `db:"price" json:"price" form:"price"`
	ShippingAddress string    `db:"shipping_address" json:"shipping_address" form:"shipping_address"`
	Status          string    `db:"status" json:"status" form:"status"`
	CreatedAt       time.Time `db:"created_at" json:"created_at" form:"created_at"`
}
