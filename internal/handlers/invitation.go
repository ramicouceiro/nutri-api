package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"`
}

type ValidateInvitationResponse struct {
	Valid        bool   `json:"valid"`
	Nutritionist struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	} `json:"nutritionist"`
	PatientEmail string `json:"patient_email"`
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateInvitation genera una nueva invitación para un paciente
func CreateInvitation(c *gin.Context) {
	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener nutricionista autenticado
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	nutritionist := userCtx.(models.User)

	// Generar token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	invitation := models.Invitation{
		NutritionistID: nutritionist.ID,
		PatientEmail:   req.Email,
		Token:          token,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour),
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	c.JSON(http.StatusCreated, invitation)
}

// ValidateInvitation verifica si un token es válido
func ValidateInvitation(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	var invitation models.Invitation
	if err := database.DB.Where("token = ? AND status = 'pending'", token).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired invitation"})
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation expired"})
		return
	}

	// Obtener datos del nutricionista
	var nutritionist models.User
	if err := database.DB.First(&nutritionist, invitation.NutritionistID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nutritionist not found"})
		return
	}

	response := ValidateInvitationResponse{
		Valid: true,
		Nutritionist: struct {
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}{
			Name:    nutritionist.Name,
			Surname: nutritionist.Surname,
		},
		PatientEmail: invitation.PatientEmail,
	}

	c.JSON(http.StatusOK, response)
}
