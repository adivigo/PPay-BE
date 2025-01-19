package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ppay/docs"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backend P-Pay
// @version         1.0
// @description     This is example for server P-Pay Aplication

// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in						   header
// @name					   Authorization

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowMethods:    []string{"POST", "GET", "PATCH", "DELETE"},
	}))
	r.Static("/public/images", "public/images")
	routes.UserRoutes(r)
	routes.TopupRoute(r)
	routes.TransferRoute(r)
	routes.AuthRoutes(r)
	routes.TransactionRoute(r)
	routes.ChangePasswordRoutes(r)
	r.Run()
}
