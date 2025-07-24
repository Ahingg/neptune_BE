package testCaseModel

import (
	"github.com/google/uuid"
	"time"
)

type TestCase struct {
	CaseID    uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Number    int       `gorm:"primaryKey;type:int;"`
	InputUrl  string    `gorm:"type:varchar(255);not null"`
	OutputUrl string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
