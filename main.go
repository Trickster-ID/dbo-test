package main

import (
	"database/sql"

	"github.com/Trickster-ID/dbo/config"
	"github.com/Trickster-ID/dbo/controller"
	docs "github.com/Trickster-ID/dbo/docs"
	"github.com/Trickster-ID/dbo/middleware"
	"github.com/Trickster-ID/dbo/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	db        *sql.DB                 = config.SetUpDatabaseConnection()
	jwtsvc    service.JWTService      = service.NewJWTService()
	authctrl  controller.Authentikasi = controller.NewAuthentikasi(db, jwtsvc)
	custctrl  controller.Customer     = controller.NewCustomer(db)
	orderctrl controller.Order        = controller.NewOrder(db)
)

// @title           Swagger DBO Assessment API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Pikri
// @contact.email  pikritaufanaziz@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//group to /api
	apiGroup := r.Group("/api")
	{
		//group by module auth
		authModuleGroup := apiGroup.Group("/auth")
		{
			authModuleGroup.POST("", authctrl.InsertLoginData)
			authModuleGroup.GET("", middleware.AuthJWT(jwtsvc), authctrl.GetLoginData)
		}
		//group by module customer
		customerModuleGroup := apiGroup.Group("/customer", middleware.AuthJWT(jwtsvc))
		{
			customerModuleGroup.GET("/", custctrl.GetWithPaginate)
			customerModuleGroup.GET("/:id", custctrl.GetDetail)
			customerModuleGroup.POST("/", custctrl.Insert)
			customerModuleGroup.PUT("/:id", custctrl.Update)
			customerModuleGroup.DELETE("/:id", custctrl.Delete)
		}
		//group by module order
		orderModuleGroup := apiGroup.Group("/order", middleware.AuthJWT(jwtsvc))
		{
			orderModuleGroup.GET("/", orderctrl.GetWithPaginate)
			orderModuleGroup.GET("/:id", orderctrl.GetDetail)
			orderModuleGroup.POST("/", orderctrl.Insert)
			orderModuleGroup.PUT("/:id", orderctrl.Update)
			orderModuleGroup.DELETE("/:id", orderctrl.Delete)
		}
	}
	//run and listening on port 8080
	r.Run(":8080")
}
