package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func ChangePasswordRoutes(router *gin.Engine) {
	changePasswordRoutes := router.Group("/change-password")
	{
		changePasswordRoutes.POST("", middlewares.ValidateToken(), controllers.CheckPassword)
	}
}
