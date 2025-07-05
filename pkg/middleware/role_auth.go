package middleware

import (
	"github.com/gin-gonic/gin"
	"neptune/backend/models/user"
	"net/http"
)

func RequireRole(requiredRoles ...user.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString("role")
		if roleStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Role not found in token"})
			c.Abort()
			return
		}

		role := user.Role(roleStr)
		for _, allowed := range requiredRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role"})
		c.Abort()
	}
}
