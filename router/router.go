package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"neptune/backend/handlers/semester"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/models/user"
	"neptune/backend/pkg/middleware"
)

func NewRouter(userHandler *userHand.UserHandler, semesterHandler *semester.SemesterHandler) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
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

	// Public auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", userHandler.LoginHandler)
		authGroup.POST("logout", middleware.RequireAuth(), userHandler.LogOutHandler)
		authGroup.GET("/me", middleware.RequireAuth(), userHandler.MeHandler)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireAuth(), middleware.RequireRole(user.RoleAdmin))
	{
		adminGroup.POST("/sync-semester", semesterHandler.SyncSemestersHandler)
		adminGroup.GET("/semesters", semesterHandler.GetSemestersHandler)
	}

	assistantGroup := r.Group("/assistant")
	assistantGroup.Use(middleware.RequireAuth(), middleware.RequireRole(user.RoleAssistant))
	{
		// TODO: Implement assistant-specific routes
	}

	studentGroup := r.Group("/student")
	studentGroup.Use(middleware.RequireAuth(), middleware.RequireRole(user.RoleStudent))
	{
		// TODO: Implement student-specific routes
	}

	return r
}
