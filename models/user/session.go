package user

import "time"

type MessierToken struct {
	UserID              string    `json:"id" gorm:"primaryKey;type:uuid"`
	MessierAccessToken  string    `json:"token" gorm:"uniqueIndex;not null"`
	MessierTokenExpires time.Time `json:"expires_at" gorm:"not null"`
}
