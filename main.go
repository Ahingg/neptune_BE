package main

import (
	models "neptune/backend/models/class"
	contestModel "neptune/backend/models/contest"
	semester "neptune/backend/models/semester"
	submissionModel "neptune/backend/models/submission"
	testCaseModel "neptune/backend/models/test_case"
	"neptune/backend/models/user"
	"neptune/backend/pkg/container"
	"neptune/backend/pkg/database"
	"neptune/backend/pkg/utils"
	"neptune/backend/router"
	"net/http"
	"os"

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
		&semester.Semester{},
		&user.MessierToken{},
		&models.Class{},
		&contestModel.Contest{},
		&contestModel.Case{},
		&testCaseModel.TestCase{},
		&models.ClassStudent{},
		&models.ClassAssistant{},
		&contestModel.ContestCase{}, // NEW: Migrate ContestCase (FKs to Contest and Case)
		&contestModel.ClassContest{},
		&submissionModel.Submission{},
		&submissionModel.SubmissionResult{},
		&contestModel.GlobalContestDetail{},
	); err != nil {
		utils.CheckPanic(err)
	}

	handlerContainer := container.NewHandlerContainer(db)

	r := router.NewRouter(
		&(handlerContainer.UserHandler),
		&(handlerContainer.InternalSemesterHandler),
		&(handlerContainer.ClassHandler),
		&(handlerContainer.ContestHandler),
		&(handlerContainer.CaseHandler),
		&(handlerContainer.TestCaseHandler),
		&(handlerContainer.WebSocketHandler),
		&(handlerContainer.SubmissionHandler),
		&(handlerContainer.LanguageHandler),
		&(handlerContainer.LeaderboardHandler),
	)

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
