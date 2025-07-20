package models

import (
	"github.com/google/uuid"
	contestModel "neptune/backend/models/contest"
	"neptune/backend/models/user"
)

type Class struct {
	ClassTransactionID uuid.UUID                   `gorm:"primaryKey;not null"`
	ClassCode          string                      `gorm:"not null"`
	CourseOutlineID    uuid.UUID                   `gorm:"not null"`
	SemesterID         uuid.UUID                   `gorm:"not null"`
	Students           []ClassStudent              `gorm:"foreignKey:ClassTransactionID;references:ClassTransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Assistants         []ClassAssistant            `gorm:"foreignKey:ClassTransactionID;references:ClassTransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Contests           []contestModel.ClassContest `gorm:"foreignKey:ClassTransactionID;references:ClassTransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ClassStudent struct {
	ClassTransactionID uuid.UUID `gorm:"not null;uniqueIndex:idx_class_student_assignment,priority:1"`           // Foreign key to the Class table
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_class_student_assignment,priority:2"` // Foreign key to the User table
	User               user.User `gorm:"foreignKey:UserID"`
	Class              Class     `gorm:"foreignKey:ClassTransactionID"`
}

type ClassAssistant struct {
	ClassTransactionID uuid.UUID `gorm:"not null;uniqueIndex:idx_class_assistant_assignment,priority:1"`           // Foreign key to the Class table
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_class_assistant_assignment,priority:2"` // Foreign key to the User table
	User               user.User `gorm:"foreignKey:UserID"`
}
