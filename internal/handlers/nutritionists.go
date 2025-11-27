package handlers

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"github.com/gin-gonic/gin"
)

type NutritionistResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func GetNutritionists(c *gin.Context) {

    var nutritionists []models.Nutritionist
    if err := database.DB.Find(&nutritionists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nutritionists"})
        return
    }
    var response []NutritionistResponse
    for _, nutritionist := range nutritionists {
        response = append(response, NutritionistResponse{
            ID:    nutritionist.ID,
            Name:  nutritionist.Name,
            Email: nutritionist.Email,
            Phone: nutritionist.Phone,
        })
    }
    c.JSON(http.StatusOK, response)
}