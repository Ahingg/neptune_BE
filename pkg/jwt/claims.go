package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"neptune/backend/models/user"
	"os"
)

type Claims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Role     user.Role `json:"role"`
	jwt.RegisteredClaims
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))
