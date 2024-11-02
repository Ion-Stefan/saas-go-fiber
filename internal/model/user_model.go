package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email"`
	ID        uint           `gorm:"primaryKey" json:"id"`
	IsAdmin   bool           `gorm:"default:false" json:"isAdmin"`
}
