package submissionModel

import (
	"github.com/google/uuid"
	"time"
)

type SubmissionStatus string // Enum for submission status
const (
	SubmissionStatusJudging             SubmissionStatus = "Judging"
	SubmissionStatusAccepted            SubmissionStatus = "Accepted"
	SubmissionStatusWrongAnswer         SubmissionStatus = "Wrong Answer"
	SubmissionStatusTimeLimitExceeded   SubmissionStatus = "Time Limit Exceeded"
	SubmissionStatusMemoryLimitExceeded SubmissionStatus = "Memory Limit Exceeded"
	SubmissionStatusCompileError        SubmissionStatus = "Compile Error"
	SubmissionStatusRuntimeError        SubmissionStatus = "Runtime Error"
	SubmissionStatusInternalError       SubmissionStatus = "Internal Error"
)

func (s SubmissionStatus) String() string {
	return string(s)
}

type Submission struct {
	ID                 uuid.UUID        `gorm:"primaryKey;type:uuid;"`
	CaseID             uuid.UUID        `gorm:"type:uuid;not null"`
	UserID             uuid.UUID        `gorm:"type:uuid;not null"`
	LanguageID         int              `gorm:"type:int;not null"`
	Status             SubmissionStatus `gorm:"type:varchar(50);not null"`
	SourceCodePath     string           `gorm:"type:varchar(255);not null"`
	Score              int              `gorm:"default:0"`
	ContestID          *uuid.UUID       `gorm:"type:uuid"`
	ClassTransactionID *uuid.UUID       `gorm:"type:uuid"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	
	SubmissionResults []SubmissionResult `gorm:"foreignKey:SubmissionID"`
}
