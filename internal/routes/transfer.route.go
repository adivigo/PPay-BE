package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func TransferRoute(router *gin.Engine) {
	trasnferGroup := router.Group("/transfer")
	{
		trasnferGroup.POST("/:id", middlewares.ValidateToken(), controllers.Transfer)
	}
}
