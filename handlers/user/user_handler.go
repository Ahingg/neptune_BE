package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
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
	secure := false
	c.SetCookie("token", accessToken, accessTokenMaxAge, "/", domain, secure, true)
	c.JSON(http.StatusOK, gin.H{"user": loginResp})
}

func (handler *UserHandler) MeHandler(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	name := c.GetString("name")
	fmt.Println(username, userID, role, name)
	if username == "" || userID == "" || role == "" || name == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"username": username,
			"name":     name,
			"role":     role,
		},
	})
}

func (handler *UserHandler) LogOutHandler(c *gin.Context) {
	domain := os.Getenv("FRONTEND_URL")

	// set secure to false on prod
	secure := false
	c.SetCookie("token", "", -1, "/", domain, secure, true) // Set cookie with empty value and negative max age to delete it
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
