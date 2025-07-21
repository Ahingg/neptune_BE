package responses

import (
	"github.com/google/uuid"
	"time"
)

type CaseResponse struct {
	ID            uuid.UUID `json:"case_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PDFFileUrl    string    `json:"pdf_file_url"`
	TimeLimitMs   int       `json:"time_limit_ms"`
	MemoryLimitMb int       `json:"memory_limit_mb"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
