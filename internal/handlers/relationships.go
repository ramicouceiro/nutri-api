package handlers

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"

	"github.com/gin-gonic/gin"
)

// LinkPatientToNutritionist vincula un paciente con un nutricionista
func LinkPatientToNutritionist(patientID, nutritionistID uint) error {
	var patient models.User
	if err := database.DB.First(&patient, patientID).Error; err != nil {
		return err
	}

	var nutritionist models.User
	if err := database.DB.First(&nutritionist, nutritionistID).Error; err != nil {
		return err
	}

	// Verificar si ya existe la relación
	var found []models.User
	if err := database.DB.Model(&nutritionist).Association("Patients").Find(&found, "id = ?", patientID); err != nil {
		return err
	}
	if len(found) > 0 {
		return nil // Ya existe, no hacemos nada
	}

	// Crear relación
	if err := database.DB.Model(&nutritionist).Association("Patients").Append(&patient); err != nil {
		return err
	}

	return nil
}

// GetMyPatients obtiene los pacientes del nutricionista autenticado
func GetMyPatients(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	nutritionist := userCtx.(models.User)

	var patients []models.User
	// Usamos Preload o Association para traer los pacientes
	// Como definimos "Patients" en el modelo User con many2many:
	if err := database.DB.Model(&nutritionist).Association("Patients").Find(&patients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// GetMyNutritionists obtiene los nutricionistas del paciente autenticado
func GetMyNutritionists(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	patient := userCtx.(models.User)

	var nutritionists []models.User
	if err := database.DB.Model(&patient).Association("Nutritionists").Find(&nutritionists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nutritionists"})
		return
	}

	c.JSON(http.StatusOK, nutritionists)
}
