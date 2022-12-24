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
	db       *sql.DB                 = config.SetUpDatabaseConnection()
	jwtsvc   service.JWTService      = service.NewJWTService()
	authctrl controller.Authentikasi = controller.NewAuthentikasi(db, jwtsvc)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	apiGroup := r.Group("/api")
	{
		moduleGroup := apiGroup.Group("/auth")
		{
			moduleGroup.GET("/logindata", middleware.AuthJWT(jwtsvc), authctrl.GetLoginData)
			moduleGroup.POST("/logindata", authctrl.InsertLoginData)
		}
	}
	r.Run(":8080")
}
