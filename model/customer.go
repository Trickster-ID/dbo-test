package model

import "time"

type User struct {
	UserID    int       `db:"user_id" json:"user_id" form:"user_id"`
	Username  string    `db:"username" json:"username" form:"username"`
	Password  string    `db:"password" json:"password" form:"password"`
	Email     string    `db:"email" json:"email" form:"email"`
	FirstName string    `db:"first_name" json:"first_name" form:"first_name"`
	LastName  string    `db:"last_name" json:"last_name" form:"last_name"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin" form:"is_admin"`
	CreatedAt time.Time `db:"created_at" json:"created_at" form:"created_at"`
}

type UserPaginate struct {
	UserID   int
	Username string
	Fullname string
}
