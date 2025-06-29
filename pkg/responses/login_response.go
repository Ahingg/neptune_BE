package responses

import "github.com/google/uuid"

type LoginResponse struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	ProfileID string    `json:"profile_id"`
}
