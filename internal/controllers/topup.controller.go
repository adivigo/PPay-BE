package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

func Topup(c *gin.Context) {
	var input struct {
		Amount          float64 `json:"amount" binding:"required"`
		PaymentMethodID uint    `json:"payment_method_id" binding:"required"`
	}
	
	response := lib.NewResponse(c)
	userId, exists := c.Get("UserId")
	if !exists {

		// fmt.Println(userId)
		// fmt.Println(exists)
		response.Unauthorized("Unauthorized access", nil)
		return
	}

	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Validate input
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Begin a database transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin database transaction"})
		return
	}

	// Create the main transaction record
	transaction := models.Transaction{
		UserID:          uint(id),
		Amount:          input.Amount,
		TransactionType: "top_up",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "details": err.Error()})
		return
	}

	// Create the top-up transaction record
	topUp := models.TopupTransaction{
		TransactionID:   transaction.ID,
		PaymentMethodID: input.PaymentMethodID,
	}

	if err := tx.Create(&topUp).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create top-up transaction", "details": err.Error()})
		return
	}

	var wallet models.Wallet
	wallet.Balance += input.Amount
	wallet.UserID = uint(id)
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to update wallet balance", err.Error())
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit database transaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Top-up transaction created successfully",
		"transaction": map[string]interface{}{
			"id":          transaction.ID,
			"amount":      transaction.Amount,
			"type":        transaction.TransactionType,
			"created_at":  transaction.CreatedAt,
			"top_up": map[string]interface{}{
				"id":              topUp.ID,
				"payment_method":  topUp.PaymentMethodID,
				"created_at":      topUp.CreatedAt,
			},
		},
	})
}
