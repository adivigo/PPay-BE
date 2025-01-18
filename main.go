package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()

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
	r.Run()
}
