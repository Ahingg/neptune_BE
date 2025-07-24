package requests

import "time"

type CreateContestRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Scope       string     `json:"scope" binding:"required"` // e.g., "public", "class"
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
}

type UpdateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"required"` // e.g., "public", "class"
}
