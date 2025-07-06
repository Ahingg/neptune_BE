package model

import "time"

type Semester struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	Description string     `gorm:"not null"`
	Start       time.Time  `gorm:"not null"` // Assuming Start is always present and valid
	End         *time.Time `gorm:""`
}
