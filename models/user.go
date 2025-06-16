package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null" json:"role"` // Assistant, SubDev, Student, Lecturer
	Class     string    `json:"class"`                // For students
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (u *User) ValidateCredentials(password string) bool {
	// TODO: Implement proper password hashing and validation
	return u.Password == password
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	if user.Role != "Assistant" && user.Role != "SubDev" && user.Role != "Student" && user.Role != "Lecturer" {
		return errors.New("invalid role")
	}
	return db.Create(user).Error
}

func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}
