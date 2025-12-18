package models

import (
	"time"

	"gorm.io/gorm"
)

	type Appointment struct {
		ID             uint           `gorm:"primarykey" json:"id"`
		UserID         uint           `gorm:"not null" json:"user_id"`
		NutritionistID uint           `gorm:"not null" json:"nutritionist_id"`
		ScheduledAt    time.Time      `gorm:"not null" json:"scheduled_at"`
		Location       string         `json:"location"`
		Notes          string         `json:"notes"`
		CreatedAt      time.Time      `json:"created_at"`
		UpdatedAt      time.Time      `json:"updated_at"`
		DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
		Nutritionist   *Nutritionist  `gorm:"foreignKey:NutritionistID" json:"nutritionist,omitempty"`
		User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}