package middleware

import (
	"net/http"
	"nutri-api/internal/database"
	"nutri-api/internal/models"
	"nutri-api/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el JWT en cada request
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extraer token (formato: "Bearer TOKEN")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Guardar user en el contexto
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if nutritionist profile exists
		if user.Role == "nutritionist" {
			var count int64
			database.DB.Model(&models.Nutritionist{}).Where("user_id = ?", user.ID).Count(&count)
			user.HasProfile = count > 0
		} else {
			user.HasProfile = true
		}

		c.Set("user", user)
		c.Next()
	}
}