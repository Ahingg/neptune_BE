package main

import (
	"neptune/backend/models/user"
	"neptune/backend/pkg/container"
	"neptune/backend/pkg/database"
	"neptune/backend/pkg/utils"
	"neptune/backend/router"
	"net/http"
	"os"

	"neptune/backend/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()

	if err != nil {
		utils.CheckPanic(err)
	}
}

func main() {
	// Auto migrate schemas
	db := database.Connect()

	if err := db.AutoMigrate(
		&user.User{},
		&models.Case{},
		&models.TestCase{},
		&models.Submission{}); err != nil {
		utils.CheckPanic(err)
	}

	handlerContainer := container.NewHandlerContainer(db)
	r := router.NewRouter(&(handlerContainer.UserHandler))

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}
	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if err := r.Run(":" + port); err != nil {
		panic("Failed to start server: " + err.Error())
	}

	if err := server.ListenAndServe(); err != nil {
		panic("Failed to listen and serve: " + err.Error())
	}
}
