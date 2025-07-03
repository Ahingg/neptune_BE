package log_on

import (
	"github.com/google/uuid"
	"time"
)

type studentResponse struct {
	UserID   uuid.UUID `json:"UserId"` // UUID string
	UserName string    `json:"UserName"`
	Name     string    `json:"Name"`
}

type tokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"` // Seconds until access token expires
}

type LogOnStudentResponse struct {
	Student studentResponse `json:"User"`
	Token   tokenResponse   `json:"Token"`
}
