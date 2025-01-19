package models

import (
	"time"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;unique;constraint:OnDelete:CASCADE"`
	Balance   float64   `gorm:"type:money;default:0.00"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
