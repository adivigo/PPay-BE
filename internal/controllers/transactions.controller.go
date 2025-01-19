package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/lib"
)

type TransactionHistoryResponse struct {
	ID              uint    `json:"id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at"`
	FromUser        *string `json:"from_user,omitempty"`
	ToUser          *string `json:"to_user,omitempty"`
}

type IncomeResponse struct {
	Income string `json:"income"`
}

type ExpenseResponse struct {
	Expense string `json:"expense"`
}

// GetTransactionHistory retrieves a paginated list of transaction history for a user
func GetTransactionHistory(c *gin.Context) {
	response := lib.NewResponse(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		response.BadRequest("Invalid page number", nil)
		return
	}
	if limit < 1 {
		response.BadRequest("Invalid limit number", nil)
		return
	}

	offset := (page - 1) * limit

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to fetch transaction history with user details
	var transactions []TransactionHistoryResponse
	query := `
        SELECT t.id, 
               amount, 
               t.transaction_type, 
               t.notes, 
               to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
               fu.fullname AS from_user,
               tu.fullname AS to_user
        FROM transactions t
        LEFT JOIN transfer_transactions tt ON t.id = tt.transaction_id
        LEFT JOIN users fu ON tt.target_user_id = fu.id
        LEFT JOIN users tu ON t.user_id = tu.id
        WHERE t.user_id = ? AND t.is_deleted = false
        ORDER BY t.created_at DESC
        LIMIT ? OFFSET ?
    `
	if err := initializers.DB.Raw(query, id, limit, offset).Scan(&transactions).Error; err != nil {
		response.InternalServerError("Failed to retrieve transaction history", err.Error())
		return
	}

	// Count total transactions for pagination
	var totalCount int64
	countQuery := `
        SELECT COUNT(*) 
        FROM transactions t
        WHERE t.user_id = ? AND t.is_deleted = false
    `
	if err := initializers.DB.Raw(countQuery, id).Scan(&totalCount).Error; err != nil {
		response.InternalServerError("Failed to count transactions", err.Error())
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Build pagination info
	pageInfo := &lib.PageInfo{
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		TotalPage:   totalPages,
		TotalData:   int(totalCount),
	}
	if page >= totalPages {
		pageInfo.NextPage = 0
	}
	if page <= 1 {
		pageInfo.PrevPage = 0
	}

	// Return success response
	response.GetAllSuccess("Success get transaction history", transactions, pageInfo)
}

func GetUserIncome(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to get income transactions
	var income IncomeResponse
	query := `
        select 
		sum(amount) 
		as income from 
		transactions t where user_id = ? and transaction_type ='top_up'
    `
	if err := initializers.DB.Raw(query, id).Scan(&income).Error; err != nil {
		response.InternalServerError("Failed to retrieve income transactions", err.Error())
		return
	}

	response.Success("Success get user income", income)
}

func GetUserExpenses(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to get expense transactions
	var expenses ExpenseResponse
	query := `
        select 
		sum(amount) 
		as expense from 
		transactions t where user_id = ? and transaction_type ='transfer'
    `
	if err := initializers.DB.Raw(query, id).Scan(&expenses).Error; err != nil {
		response.InternalServerError("Failed to retrieve expense transactions", err.Error())
		return
	}

	response.Success("Success get user expenses", expenses)
}
