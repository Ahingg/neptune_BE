package requests

import (
	"time"

	"github.com/google/uuid"
)

type CreateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	// Optional class assignment data
	ClassTransactionID *uuid.UUID `json:"class_transaction_id"`
	StartTime          *time.Time `json:"start_time"`
	EndTime            *time.Time `json:"end_time"`
}

type UpdateContestRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
