package main

import (
	"log"
	"os"

	"neptune/backend/handlers"
	"neptune/backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schemas
	db.AutoMigrate(&models.User{}, &models.Case{}, &models.TestCase{}, &models.Submission{})

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.HandleLogin)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(handlers.AuthMiddleware())
		{
			// Case routes
			cases := protected.Group("/cases")
			{
				cases.GET("/", handlers.GetCases)
				cases.POST("/", handlers.CreateCase)
				cases.GET("/:id", handlers.GetCase)
				cases.PUT("/:id", handlers.UpdateCase)
				cases.DELETE("/:id", handlers.DeleteCase)
				cases.POST("/:id/test-cases", handlers.AddTestCase)
				cases.DELETE("/:id/test-cases/:testCaseId", handlers.DeleteTestCase)
			}

			// Submission routes
			submissions := protected.Group("/submissions")
			{
				submissions.POST("/", handlers.CreateSubmission)
				submissions.GET("/", handlers.GetSubmissions)
				submissions.GET("/:id", handlers.GetSubmission)
			}

			// Leaderboard routes
			leaderboard := protected.Group("/leaderboard")
			{
				leaderboard.GET("/class/:classId", handlers.GetClassLeaderboard)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
