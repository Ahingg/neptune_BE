package contestModel

import (
	"github.com/google/uuid"
	"time"
)

type ContestCase struct {
	ContestID   uuid.UUID `gorm:"primaryKey;type:uuid;"` // Composite primary key part 1, FK to Contest
	CaseID      uuid.UUID `gorm:"primaryKey;type:uuid;"` // Composite primary key part 2, FK to Case
	ProblemCode string    `gorm:"not null"`              // e.g., "A", "B", "C" for this problem in this contest

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations for GORM to understand the relationships
	Case Case `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
