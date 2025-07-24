package contestModel

import (
	"github.com/google/uuid"
	"time"
)

type ClassContest struct {
	// Foreign key to Class (using ClassTransactionID as PK of Class)
	ClassTransactionID uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	// Foreign key to Contest
	ContestID uuid.UUID `gorm:"primaryKey;type:uuid;not null"`

	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations for GORM to understand the relationships
	Contest Contest `gorm:"foreignKey:ContestID;references:ID"`
}
