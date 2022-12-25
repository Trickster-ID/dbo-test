package model

import "time"

type Product struct {
	ProductID   int       `db:"product_id" json:"product_id" form:"product_id"`
	Name        string    `db:"name" json:"name" form:"name"`
	Price       float64   `db:"price" json:"price" form:"price"`
	Description string    `db:"description" json:"description" form:"description"`
	ImageURL    string    `db:"image_url" json:"image_url" form:"image_url"`
	Quantity    int       `db:"quantity" json:"quantity" form:"quantity"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" form:"created_at"`
}
