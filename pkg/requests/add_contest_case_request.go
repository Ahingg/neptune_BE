package requests

import "github.com/google/uuid"

type AddCasesToContestRequest struct {
	Problems []struct {
		CaseID      uuid.UUID `json:"case_id" binding:"required"`
		ProblemCode string    `json:"problem_code" binding:"required"` // e.g., "A", "B", "C"
	} `json:"problems" binding:"required,min=1"`
}
