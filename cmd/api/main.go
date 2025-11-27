package main

import (
	"log"
	"nutri-api/internal/database"
	"nutri-api/internal/handlers"
	"nutri-api/internal/middleware"
	"nutri-api/internal/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found")
	} else {
		log.Println("✅ .env loaded")
	}

	// Conectar a base de datos
	log.Println("🔌 Connecting to database...")
	database.Connect()
	log.Println("✅ Database connected!")

	// Auto migrate models
	log.Println("🔄 Running migrations...")
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Appointment{})
	database.DB.AutoMigrate(&models.Nutritionist{})
	log.Println("✅ Migrations completed!")

	// Setup Gin
	log.Println("🚀 Setting up Gin...")
	r := gin.Default()

	// Configurar trusted proxies para ngrok
	r.SetTrustedProxies(nil)

	// Logging middleware
	r.Use(func(c *gin.Context) {
		log.Printf("📨 Request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
		log.Printf("✅ Response: %d", c.Writer.Status())
	})

	// CORS middleware - SIMPLIFICADO
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/login", handlers.Login)
		api.POST("/signup", handlers.Signup)
		api.GET("/appointments/create", handlers.CreateTestAppointment)
		api.GET("/appointments/get", handlers.GetAppointments)
		api.GET("/nutritionists/get", handlers.GetNutritionists)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/me", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(200, user)
			})
		}
	}

	// Health check
	r.GET("/up", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Println("📱 Ready to accept connections!")
	r.Run("0.0.0.0:" + port)
}