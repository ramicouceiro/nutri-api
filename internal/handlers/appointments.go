package handlers

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentResponse struct {
	ID             		uint   `json:"id"`
	UserID         		uint   `json:"user_id"`
	NutritionistID 		uint   `json:"nutritionist_id"`
	NutritionistName 	string `json:"nutritionist_name"` // Nuevo campo para el nombre del nutricionista
	Location 	  		string `json:"location"`
	ScheduledAt    		string `json:"scheduled_at"`
	Notes          		string `json:"notes"`
}

// Creo una ruta de prueba simple, que crea un appointment fijo y lo devuelve
func CreateTestAppointment(c *gin.Context) {
	appointment := models.Appointment{
		UserID:        1,
		NutritionistID: 1,
		ScheduledAt:  time.Now().Add(24 * time.Hour),
		Notes:         "This is a test appointment",
	}
	if err := database.DB.Create(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test appointment"})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func GetAppointments(c *gin.Context) {
	userID := c.Query("user_id")
	now := time.Now()
	var appointments []models.Appointment
	if err := database.DB.
		Preload("Nutritionist"). // Esto carga la relación
		Where("user_id = ?", userID).
		Or("nutritionist_id = ?", userID).
		Where("scheduled_at >= ?", now).
		Order("scheduled_at ASC").
		Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	var response []AppointmentResponse
	for _, appt := range appointments {
		nutritionistName := ""
		if appt.Nutritionist != nil {
			nutritionistName = appt.Nutritionist.Name
		}
		
		response = append(response, AppointmentResponse{
			ID:               appt.ID,
			UserID:           appt.UserID,
			NutritionistID:   appt.NutritionistID,
			NutritionistName: nutritionistName,
			Location:         appt.Location,
			ScheduledAt:      appt.ScheduledAt.Format(time.RFC3339),
			Notes:            appt.Notes,
		})
	}
	c.JSON(http.StatusOK, response)
}