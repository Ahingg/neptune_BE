package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	"neptune/backend/pkg/utils"
	"neptune/backend/services/user"
	"net/http"
	"os"
	"regexp"
	"time"
)

type UserHandler struct {
	service user.UserService
}

func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

var nimRegex = regexp.MustCompile(`^\d{10}$`)

func (handler *UserHandler) LoginHandler(c *gin.Context) {
	var (
		req requests.LoginRequest
		err error
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var (
		loginResp   *responses.LoginResponse
		accessToken string
		expires     time.Time
	)
	if nimRegex.MatchString(req.Username) {
		loginResp, accessToken, expires, err = handler.service.LoginStudent(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid credentials %s", err)})
			return
		}
	} else {
		loginResp, accessToken, expires, err = handler.service.LoginAssistant(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	accessTokenMaxAge := int(time.Until(expires).Seconds())
	if accessTokenMaxAge <= 0 {
		accessTokenMaxAge = 3600
	}

	domain := os.Getenv("FRONTEND_URL")

	// set secure to false on prod
	secure := os.Getenv("APP_ENV") == "production"
	c.SetCookie("token", accessToken, accessTokenMaxAge, "/", domain, secure, true)
	c.JSON(http.StatusOK, gin.H{"user": loginResp})
}

func (handler *UserHandler) MeHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id") // Assuming your middleware sets "userID"
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error: Invalid user ID format in context"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	meResp, err := handler.service.GetDetailedUserProfile(ctx, userID)
	if err != nil {
		log.Printf("Error getting detailed user profile for %s: %v", userID.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve user profile: %v", err.Error())})
		return
	}
	if meResp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": meResp})
}

func (handler *UserHandler) LogOutHandler(c *gin.Context) {
	domain := os.Getenv("FRONTEND_URL")

	// delete user external token if exists
	err := handler.service.DeleteUserAccessToken(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed on Log Out: %s", err)})
		utils.CheckPanic(fmt.Errorf("failed deleting access token: %s", err))
	}

	// set secure to false on prod
	secure := os.Getenv("APP_ENV") == "production"
	c.SetCookie("token", "", -1, "/", domain, secure, true) // Set cookie with empty value and negative max age to delete it
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
