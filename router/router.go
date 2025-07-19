package router

import (
	caseHandler "neptune/backend/handlers/case"
	"neptune/backend/handlers/class"
	contestHandler "neptune/backend/handlers/contest"
	"neptune/backend/handlers/semester"
	testCaseHand "neptune/backend/handlers/test_case"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/models/user"
	"neptune/backend/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *userHand.UserHandler,
	semesterHandler *semester.SemesterHandler,
	classHandler *class.ClassHandler,
	contestHandler *contestHandler.ContestHandler,
	caseHandler *caseHandler.CaseHandler,
	testCaseHandler *testCaseHand.TestCaseHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Static("/public/cases", "./public/case_file") // Serve static files from the public directory
	// Serve private test cases

	// Public auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", userHandler.LoginHandler)
		authGroup.POST("logout", middleware.RequireAuth(), userHandler.LogOutHandler)
		authGroup.GET("/me", middleware.RequireAuth(), userHandler.MeHandler)
	}

	authRestrictedGroup := r.Group("/api")
	authRestrictedGroup.Use(middleware.RequireAuth())
	{
		authRestrictedGroup.GET("/semesters", semesterHandler.GetSemestersHandler)

		// Class routes (existing)
		authRestrictedGroup.GET("/debug-semesters", semesterHandler.DebugSemestersHandler)
		authRestrictedGroup.GET("/classes", classHandler.GetClassesBySemesterAndCourseHandler)
		authRestrictedGroup.GET("/classes/detail", classHandler.GetClassDetailBySemesterAndCourseHandler)               // Specific detail
		authRestrictedGroup.GET("/class-detail/:classTransactionId", classHandler.GetClassDetailByTransactionIDHandler) // General class detail by ID

		// Contest routes (general access if contests are viewable by all authenticated users)
		authRestrictedGroup.GET("/contests", contestHandler.GetAllContests)
		authRestrictedGroup.GET("/contests/:contestId", contestHandler.GetContestByID)
		authRestrictedGroup.GET("/classes/:classTransactionId/contests", contestHandler.GetContestsForClass) // Get contests assigned to a class

		// Case routes (general access if problems are viewable by all authenticated users)
		authRestrictedGroup.GET("/cases", caseHandler.GetAllCases)
		authRestrictedGroup.GET("/cases/:caseId", caseHandler.GetCaseByID)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireAuth(), middleware.RequireRole(user.RoleAdmin))
	{
		// TODO: Implement admin-specific routes
		adminGroup.POST("/sync-semester", semesterHandler.SyncSemestersHandler)
		adminGroup.POST("/sync-classes", classHandler.SyncClassesHandler)
		adminGroup.POST("/sync-class-students", classHandler.SyncClassStudentsHandler)
		adminGroup.POST("/sync-class-assistants", classHandler.SyncClassAssistantsHandler)

		adminGroup.POST("/contests", contestHandler.CreateContest)
		adminGroup.PUT("/contests/:contestId", contestHandler.UpdateContest)
		adminGroup.DELETE("/contests/:contestId", contestHandler.DeleteContest)
		adminGroup.POST("/contests/:contestId/cases", contestHandler.AddCasesToContest)

		adminGroup.POST("/classes/:classTransactionId/assign-contest", contestHandler.AssignContestToClass)
		adminGroup.DELETE("/classes/:classTransactionId/contests/:contestId", contestHandler.RemoveContestFromClass)

		adminGroup.POST("/cases", caseHandler.CreateCase)
		adminGroup.PUT("/cases/:caseId", caseHandler.UpdateCase)
		adminGroup.DELETE("/cases/:caseId", caseHandler.DeleteCase)

		adminGroup.POST("/cases/:case_id/test-cases", testCaseHandler.UploadTestCasesHandler)
		adminGroup.GET("/cases/:case_id/test-cases", testCaseHandler.GetTestCasesByCaseIDHandler)
		adminGroup.Static("/private/test_case", "./private/test_case")
	}

	return r
}
