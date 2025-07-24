package requests

import "github.com/google/uuid"

type AddCasesToContestRequest struct {
	Problems []struct {
		CaseID uuid.UUID `json:"case_id" binding:"required"`
	} `json:"problems" binding:"required,min=1"`
}
