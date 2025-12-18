package handlers

import (
	"log"
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"nutri-api/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Name           string  `json:"name" binding:"required"`
	Surname        string  `json:"surname"`
	Email          string  `json:"email" binding:"required,email"`
	Phone          string  `json:"phone" binding:"required"`
	Password       string  `json:"password" binding:"required,min=6"`
	Role           string  `json:"role"` // user or nutritionist
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
	InvitationToken string `json:"invitation_token"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

// Login maneja el inicio de sesión
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Check if nutritionist profile exists
	if user.Role == "nutritionist" {
		var count int64
		database.DB.Model(&models.Nutritionist{}).Where("user_id = ?", user.ID).Count(&count)
		user.HasProfile = count > 0
	} else {
		user.HasProfile = true // Users don't need extra profile for now
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Signup maneja el registro de nuevos usuarios
func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// AGREGAR LOG AQUÍ PARA VER EL ERROR
		log.Printf("❌ Signup validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:           req.Name,
		Surname:        req.Surname,
		Email:          req.Email,
		Phone:          req.Phone,
		Role:           req.Role,
		Height:         req.Height,
		Weight:         req.Weight,
	}

	// Determine role based on invitation
	if req.InvitationToken != "" {
		user.Role = "user"
	} else {
		user.Role = "nutritionist"
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Process Invitation if present
	if req.InvitationToken != "" {
		var invitation models.Invitation
		if err := database.DB.Where("token = ? AND status = 'pending'", req.InvitationToken).First(&invitation).Error; err == nil {
			// Check expiration
			if time.Now().Before(invitation.ExpiresAt) {
				// Create relationship
				if err := LinkPatientToNutritionist(user.ID, invitation.NutritionistID); err != nil {
					log.Printf("Failed to link patient to nutritionist: %v", err)
				}

				// Update invitation status
				invitation.Status = "accepted"
				database.DB.Save(&invitation)
			}
		}
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}