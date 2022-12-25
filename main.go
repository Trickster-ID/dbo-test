package main

import (
	"database/sql"

	"github.com/Trickster-ID/dbo/config"
	"github.com/Trickster-ID/dbo/controller"
	"github.com/Trickster-ID/dbo/middleware"
	"github.com/Trickster-ID/dbo/service"
	"github.com/gin-gonic/gin"
)

var (
	db        *sql.DB                 = config.SetUpDatabaseConnection()
	jwtsvc    service.JWTService      = service.NewJWTService()
	authctrl  controller.Authentikasi = controller.NewAuthentikasi(db, jwtsvc)
	custctrl  controller.Customer     = controller.NewCustomer(db)
	orderctrl controller.Order        = controller.NewOrder(db)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	apiGroup := r.Group("/api")
	{
		authModuleGroup := apiGroup.Group("/auth")
		{
			authModuleGroup.GET("", middleware.AuthJWT(jwtsvc), authctrl.GetLoginData)
			authModuleGroup.POST("", authctrl.InsertLoginData)
		}
		customerModuleGroup := apiGroup.Group("/customer", middleware.AuthJWT(jwtsvc))
		{
			customerModuleGroup.GET("/", custctrl.GetWithPaginate)
			customerModuleGroup.GET("/:id", custctrl.GetDetail)
			customerModuleGroup.POST("/", custctrl.Insert)
			customerModuleGroup.PUT("/:id", custctrl.Update)
			customerModuleGroup.DELETE("/:id", custctrl.Delete)
		}
		orderModuleGroup := apiGroup.Group("/order", middleware.AuthJWT(jwtsvc))
		{
			orderModuleGroup.GET("/", orderctrl.GetWithPaginate)
			orderModuleGroup.GET("/:id", orderctrl.GetDetail)
			orderModuleGroup.POST("/", orderctrl.Insert)
			orderModuleGroup.PUT("/:id", orderctrl.Update)
			orderModuleGroup.DELETE("/:id", orderctrl.Delete)
		}
	}
	r.Run(":8080")
}
