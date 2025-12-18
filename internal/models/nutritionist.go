package models

import (
	"time"

	"gorm.io/gorm"
)

type Nutritionist struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	User           User           `json:"user"`
	Name           string         `gorm:"not null" json:"name"`
	Surname        string         `json:"surname"`
	Email          string         `gorm:"unique;not null" json:"email"`
	Phone          string         `json:"phone"`
	Specialty      string         `json:"specialty"`
	Rating         float64        `json:"rating"`
	Image          string         `json:"image"`
	Description    string         `json:"description"`
	OfficeLocation string         `json:"office_location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}