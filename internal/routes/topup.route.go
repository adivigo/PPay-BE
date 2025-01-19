package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func TopupRoute(router *gin.Engine) {
	topupGroup := router.Group("/topup")
	{
		topupGroup.POST("", middlewares.ValidateToken(), controllers.Topup)
	}
}
