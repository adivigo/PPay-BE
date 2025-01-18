package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

func Transfer(c *gin.Context) {
	var input struct {
		Amount          float64 `form:"amount" json:"amount" binding:"required"`
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
	
	targetId, _ := strconv.Atoi(c.Param("id"))

	var userSummary models.User
	if err := initializers.DB.Model(&models.User{}).
		Select("image, fullname, phone").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&userSummary).Error; err != nil {
		response.NotFound("User not found", nil)
		return
	}

	// Create the top-up transaction record
	transfer := models.TransferTransaction{
		TransactionID:   transaction.ID,
		TargetUserID: uint(targetId),
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transfer transaction", "details": err.Error()})
		return
	}

	var wallet models.Wallet
	wallet.Balance -= input.Amount
	wallet.UserID = uint(id)
	
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to update wallet balance", err.Error())
		return
	}

	var walletTarget models.Wallet
	walletTarget.Balance += input.Amount
	walletTarget.UserID = uint(targetId)
	
	if err := tx.Save(&walletTarget).Error; err != nil {
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
			"id":         transaction.ID,
			"userId": userId,
			"amount":     transaction.Amount,
			"type":       transaction.TransactionType,
			"created_at": transaction.CreatedAt,
			"top_up": map[string]interface{}{
				"id":             transfer.ID,
				"targetUserId": targetId,
				"created_at":     transfer.CreatedAt,
			},
		},
	})
}
