package requests

import (
	"github.com/google/uuid"
	"time"
)

type AssignContestToClassRequest struct {
	ContestID uuid.UUID `json:"contest_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}
