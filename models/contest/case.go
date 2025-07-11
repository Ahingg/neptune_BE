package contestModel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Case struct { // Renamed from "Problem" to "Case" as per your terminology
	ID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name          string    `gorm:"not null"`
	Description   string    `gorm:"type:text"`         // Problem description (e.g., Markdown)
	PDFFileUrl    string    `gorm:"type:varchar(255)"` // URL to the problem statement PDF
	TimeLimitMs   int       `gorm:"not null"`          // Time limit in milliseconds
	MemoryLimitMb int       `gorm:"not null"`          // Memory limit in megabytes
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"` // For soft deletes
}
