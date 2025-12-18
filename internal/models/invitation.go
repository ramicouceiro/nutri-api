package models

import (
	"time"

	"gorm.io/gorm"
)

type Invitation struct {
	gorm.Model
	NutritionistID uint      `gorm:"not null" json:"nutritionist_id"`
	PatientEmail   string    `gorm:"not null" json:"patient_email"`
	Token          string    `gorm:"unique;index;not null" json:"token"`
	Status         string    `gorm:"default:'pending'" json:"status"` // pending, accepted
	ExpiresAt      time.Time `gorm:"not null" json:"expires_at"`
}
