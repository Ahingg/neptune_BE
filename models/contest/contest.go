package contestModel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`                 // Optional description
	Scope       string    `gorm:"type:varchar(50);not null"` // e.g., "public", "class"
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` // For soft deletes

	// Many-to-many relationship with Case via ContestCase
	GlobalContestDetail *GlobalContestDetail `gorm:"foreignKey:ContestID;references:ID"`
	ContestCases        []ContestCase        `gorm:"foreignKey:ContestID;references:ID"`
}

type GlobalContestDetail struct {
	ContestID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	StartTime time.Time `gorm:"not null"` // Start time of the contest
	EndTime   time.Time `gorm:"not null"` // End time of the contest
}
