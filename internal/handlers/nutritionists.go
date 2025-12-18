package handlers

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"

	"github.com/gin-gonic/gin"
)

type NutritionistResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Surname        string  `json:"surname"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Specialty      string  `json:"specialty"`
	Rating         float64 `json:"rating"`
	Image          string  `json:"image"`
	Description    string  `json:"description"`
	OfficeLocation string  `json:"office_location"`
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
            ID:             nutritionist.ID,
            Name:           nutritionist.Name,
            Surname:        nutritionist.Surname,
            Email:          nutritionist.Email,
            Phone:          nutritionist.Phone,
            Specialty:      nutritionist.Specialty,
            Rating:         nutritionist.Rating,
            Image:          nutritionist.Image,
            Description:    nutritionist.Description,
            OfficeLocation: nutritionist.OfficeLocation,
        })
    }
    c.JSON(http.StatusOK, response)
}

type CreateProfileRequest struct {
	Specialty      string  `json:"specialty" binding:"required"`
	Description    string  `json:"description" binding:"required"`
	OfficeLocation string  `json:"office_location" binding:"required"`
	Rating         float64 `json:"rating"`
	Image          string  `json:"image"`
}

func CreateNutritionistProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nutritionist := models.Nutritionist{
		UserID:         currentUser.ID,
		Name:           currentUser.Name,
		Surname:        currentUser.Surname,
		Email:          currentUser.Email,
		Phone:          currentUser.Phone,
		Specialty:      req.Specialty,
		Description:    req.Description,
		OfficeLocation: req.OfficeLocation,
		Rating:         5.0, // Default rating
		Image:          req.Image,
	}

	if err := database.DB.Create(&nutritionist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, nutritionist)
}