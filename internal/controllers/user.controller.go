package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pilinux/argon2"
	"github.com/ppay/internal/dto"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

// Create User and Wallet
func CreateUser(c *gin.Context) {
	var input struct {
		Fullname string  `json:"fullname"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required,min=6"`
		Pin      *string `json:"pin"`
		Phone    *string `json:"phone"`
		Image    *string `json:"image"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != "" {
		input.Password, _ = argon2.CreateHash(input.Password, "", argon2.DefaultParams)
	}

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	user := models.User{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: input.Password,
		Pin:      input.Pin,
		Phone:    input.Phone,
		Image:    input.Image,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Buat Wallet untuk User
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaksi
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "User and wallet created successfully",
		"user":    user,
		"wallet":  wallet,
	})
}

// Get All Users
func GetUsers(c *gin.Context) {
	response := lib.NewResponse(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")

	// Validate input for page and limit
	if page < 1 {
		response.BadRequest("Invalid input", "Page must be 1 or greater")
		return
	}
	if limit < 1 {
		response.BadRequest("Invalid input", "Limit must be 1 or greater")
		return
	}

	// Validate sorting parameters
	if sort != "id" && sort != "fullname" && sort != "phone" {
		response.BadRequest("Invalid input", "Sort must be 'id', 'fullname', or 'phone'")
		return
	}
	if order != "asc" && order != "desc" {
		response.BadRequest("Invalid input", "Order must be 'asc' or 'desc'")
		return
	}

	offset := (page - 1) * limit

	// Query database with filters and pagination
	var users []dto.UserSummaryDTO
	query := initializers.DB.Model(&models.User{}).
		Select("id, fullname, phone, image").
		Where("is_deleted = ?", false)

	// Apply search filter if provided
	if search != "" {
		query = query.Where("fullname ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply sorting, limit, and offset
	if err := query.Order(sort + " " + order).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		response.InternalServerError("Failed to get users", nil)
		return
	}

	// Count total users matching the criteria
	var totalCount int64
	countQuery := initializers.DB.Model(&models.User{}).Where("is_deleted = ?", false)
	if search != "" {
		countQuery = countQuery.Where("fullname ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&totalCount)

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

	// Return response
	response.GetAllSuccess("Success get user", users, pageInfo)
}

func GetUserByID(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	log.Println(userId)
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Use DTO for selected fields
	var userSummary dto.UserSummaryDTO

	// Query only required fields
	if err := initializers.DB.Model(&models.User{}).
		Select("id, email, image, fullname, phone").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&userSummary).Error; err != nil {
		response.NotFound("User not found", nil)
		return
	}

	// Return response with filtered data
	response.Success("Success get user", userSummary)
}

func UpdateUser(c *gin.Context) {
	response := lib.NewResponse(c)

	// Ambil userId dari konteks
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

	// Cari data user berdasarkan ID
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}
	fmt.Println("Existing User:", user)

	// Bind input data
	var req dto.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	// Update data hanya jika ada input
	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		if !IsValidPassword(*req.Password) {
			response.BadRequest("Password must be at least 8 characters long", nil)
			return
		}
		hashedPassword := lib.GenerateHash(*req.Password)
		user.Password = hashedPassword
	}
	if req.Pin != nil {
		hashPin := lib.GenerateHash(*req.Pin)
		user.Pin = &hashPin
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Image != nil {
		user.Image = req.Image
	}

	// Perbarui waktu
	user.UpdatedAt = time.Now()

	// Simpan perubahan ke database
	if err := initializers.DB.Save(&user).Error; err != nil {
		response.InternalServerError("Failed to update user", err.Error())
		return
	}

	// Respon sukses
	response.Success("Update user success", dto.UserSummaryDTO{
		Id:       int(user.ID),
		Fullname: user.Fullname,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}

// Delete User
func DeleteUser(c *gin.Context) {
	response := lib.NewResponse(c)

	// Ambil userId dari konteks
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

	// Cari data user
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	// Hapus user (soft delete)
	user.IsDeleted = true
	user.UpdatedAt = time.Now()
	if err := initializers.DB.Save(&user).Error; err != nil {
		response.InternalServerError("Failed to delete user", err.Error())
		return
	}

	// Respon sukses
	response.Success(fmt.Sprintf("User with ID %d deleted successfully", id), nil)
}
