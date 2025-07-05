package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	jwtPkg "neptune/backend/pkg/jwt"
	"net/http"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token cookie provided"})
			c.Abort()
			return
		}

		claims := &jwtPkg.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtPkg.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or expired token"})
			c.Abort()
			return
		}

		// Store the claims in Gin context for later use
		//fmt.Println("Authenticated user:", claims.Name, claims.UserID, claims.Role, claims.Username)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role.String())
		c.Set("username", claims.Username)
		c.Set("name", claims.Name)

		c.Next()
	}
}
