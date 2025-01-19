// routes/index.route.go
package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	UserRoutes(router)
	TopupRoute(router)
	TransferRoute(router)
	AuthRoutes(router)
	TransactionRoute(router)
	ChangePasswordRoutes(router)
}
