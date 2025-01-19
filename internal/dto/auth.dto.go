package dto

type RegisterDTO struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type LoginDTO struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type PinDTO struct {
	Pin string `json:"pin" form:"pin"`
}

type PasswordDTO struct {
	Password string `json:"password" form:"password"`
}
