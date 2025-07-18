package responses

import (
	"time"

	"github.com/google/uuid"
)

type ContestResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClassContestAssignmentResponse struct {
	ClassTransactionID uuid.UUID        `json:"class_transaction_id"`
	ContestID          uuid.UUID        `json:"contest_id"`
	StartTime          time.Time        `json:"start_time"`
	EndTime            time.Time        `json:"end_time"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	Contest            *ContestResponse `json:"contest"`
}

type ContestCaseProblemResponse struct {
	CaseID        uuid.UUID `json:"case_id"`
	ProblemCode   string    `json:"problem_code"`
	Name          string    `json:"name"`
	TimeLimitMs   int       `json:"time_limit_ms"`
	MemoryLimitMb int       `json:"memory_limit_mb"`
	PDFFileUrl    string    `json:"pdf_file_url"`
}

type ContestDetailResponse struct {
	ID          uuid.UUID                    `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
	Cases       []ContestCaseProblemResponse `json:"cases"`
}
