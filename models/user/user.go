package user

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleStudent   Role = "student"
	RoleAssistant Role = "assistant"
	RoleLecturer  Role = "lecturer"
	RoleAdmin     Role = "admin"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Username  string    `gorm:"uniqueIndex"`
	Name      string    `gorm:"not null"`
	Role      Role      `gorm:"not null"` // Assistant, SubDev, Student, Lecturer
	ProfileID string    `gorm:"not null"`
	CreatedAt time.Time
}
