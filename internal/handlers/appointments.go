package handlers

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentResponse struct {
	ID             		uint   `json:"id"`
	UserID         		uint   `json:"user_id"`
	NutritionistID 		uint   `json:"nutritionist_id"`
	NutritionistName 	string `json:"nutritionist_name"`
	UserName			string `json:"user_name"` // Nuevo campo para el nombre del paciente
	Location 	  		string `json:"location"`
	ScheduledAt    		string `json:"scheduled_at"`
	Notes          		string `json:"notes"`
}

// Creo una ruta de prueba simple, que crea un appointment fijo y lo devuelve
// Permite user_id y nutritionist_id por query params
func CreateTestAppointment(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "6"))
	nutritionistID, _ := strconv.Atoi(c.DefaultQuery("nutritionist_id", "5"))

	appointment := models.Appointment{
		UserID:        uint(userID),
		NutritionistID: uint(nutritionistID),
		ScheduledAt:  time.Now().Add(24 * time.Hour),
		Notes:         "This is a test appointment",
	}
	if err := database.DB.Create(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test appointment: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func GetAppointments(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user := userInterface.(models.User)

	now := time.Now()
	var appointments []models.Appointment
	
	query := database.DB.Preload("Nutritionist").Preload("User").Where("scheduled_at >= ?", now).Order("scheduled_at ASC")

	if user.Role == "nutritionist" {
		// Find the nutritionist profile for this user
		var nutritionist models.Nutritionist
		if err := database.DB.Where("user_id = ?", user.ID).First(&nutritionist).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Nutritionist profile not found"})
			return
		}
		query = query.Where("nutritionist_id = ?", nutritionist.ID)
	} else {
		// Regular user
		query = query.Where("user_id = ?", user.ID)
	}

	if err := query.Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	var response []AppointmentResponse
	for _, appt := range appointments {
		nutritionistName := ""
		if appt.Nutritionist != nil {
			nutritionistName = appt.Nutritionist.Name
		}

		userName := ""
		if appt.User != nil {
			userName = appt.User.Name + " " + appt.User.Surname
		}
		
		response = append(response, AppointmentResponse{
			ID:               appt.ID,
			UserID:           appt.UserID,
			NutritionistID:   appt.NutritionistID,
			NutritionistName: nutritionistName,
			UserName:         userName,
			Location:         appt.Location,
			ScheduledAt:      appt.ScheduledAt.Format(time.RFC3339),
			Notes:            appt.Notes,
		})
	}
	c.JSON(http.StatusOK, response)
}