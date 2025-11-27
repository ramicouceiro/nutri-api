package models

import (
	"time"

	"gorm.io/gorm"
)

type Nutritionist struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Phone     string         `json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}