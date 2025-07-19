package requests

import "github.com/google/uuid"

type SubmitCodeRequest struct {
	SourceCode         string     `json:"source_code" binding:"required"`
	CaseID             uuid.UUID  `json:"case_id" binding:"required"`
	LanguageID         int        `json:"language_id" binding:"required"`
	ContestID          *uuid.UUID `json:"contest_id"`           // Optional
	ClassTransactionID *uuid.UUID `json:"class_transaction_id"` // Optional
}
