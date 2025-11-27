package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // No enviar en JSON
	Name      string         `gorm:"not null" json:"name"`
	Phone     string         `gorm:"not null" json:"phone"`
	Role      string         `gorm:"not null;default:'user'" json:"role"` // user, nutritionist, admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword encripta la contraseña antes de guardar
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword verifica si la contraseña es correcta
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}