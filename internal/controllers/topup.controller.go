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
		Amount          float64 `json:"amount" form:"amount" binding:"required"`
		PaymentMethodID int     `json:"payment_method_id" form:"payment_method_id" binding:"required"`
	}

	response := lib.NewResponse(c)
	userId, exists := c.Get("UserId")
	if !exists {
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

	// Fetch the user's existing wallet to update its balance
	var wallet models.Wallet
	if err := tx.Where("user_id = ?", id).First(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to retrieve wallet", err.Error())
		return
	}

	// Update the wallet balance
	wallet.Balance += input.Amount
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

	type ResponseTopup struct {
		Amount float64
		Type   string
	}

	response.Success("Success Top up", ResponseTopup{
		Amount: transaction.Amount,
		Type:   transaction.TransactionType,
	})
}
