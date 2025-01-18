package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/lib"
)

func GetBalance(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}

	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query the wallet for the user
	var balance float64
	query := `SELECT balance::numeric FROM wallets WHERE user_id = $1`

	if err := initializers.DB.Raw(query, id).Scan(&balance).Error; err != nil {
		if err.Error() == "record not found" {
			response.NotFound("Wallet not found", nil)
			return
		}
		response.InternalServerError("Failed to retrieve wallet", err.Error())
		return
	}

	// Return the balance
	response.Success("Success get balance", gin.H{
		"balance": balance,
	})
}
