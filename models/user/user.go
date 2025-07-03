package user

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleStudent   Role = "Student"
	RoleAssistant Role = "Assistant"
	RoleAdmin     Role = "Admin"
)

func (r Role) String() string {
	switch r {
	case RoleStudent:
		return "Student"
	case RoleAssistant:
		return "Assistant"
	case RoleAdmin:
		return "Admin"
	}

	return "Unknown"
}

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Username  string    `gorm:"uniqueIndex"`
	Name      string    `gorm:"not null"`
	Role      Role      `gorm:"not null"` // Assistant, SubDev, Student, Lecturer
	CreatedAt time.Time
}
