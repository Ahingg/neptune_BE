package submissionModel

import "github.com/google/uuid"

type SubmissionResult struct {
	SubmissionID   uuid.UUID        `gorm:"primaryKey;type:uuid;"`
	TestcaseNumber int              `gorm:"primaryKey"`
	Status         SubmissionStatus `gorm:"type:varchar(50);not null"`
	TimeSeconds    float64
	MemoryKB       int
	Input          string `gorm:"type:text"`
	ExpectedOutput string `gorm:"type:text"`
	ActualOutput   string `gorm:"type:text"`
}
