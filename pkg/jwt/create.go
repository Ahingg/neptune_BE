package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	model "neptune/backend/models/user"
	"time"
)

func CreateJWT(userID, username string, name string, role model.Role, expire time.Time) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Name:     name,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
