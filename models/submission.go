package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	CaseID    uint      `json:"case_id"`
	Code      string    `gorm:"not null" json:"code"`
	Language  string    `gorm:"not null" json:"language"` // C or Python
	Status    string    `json:"status"`                   // Pending, Accepted, Wrong Answer, etc.
	ErrorMsg  string    `json:"error_msg"`                // For failed test cases
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubmissionRequest struct {
	CaseID   uint   `json:"case_id" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required,oneof=C Python"`
}

func CreateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmission(db *gorm.DB, id uint) (*Submission, error) {
	var submission Submission
	if err := db.First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func GetUserSubmissions(db *gorm.DB, userID uint) ([]Submission, error) {
	var submissions []Submission
	if err := db.Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func GetCaseSubmissions(db *gorm.DB, caseID uint) ([]Submission, error) {
	var submissions []Submission
	if err := db.Where("case_id = ?", caseID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func UpdateSubmissionStatus(db *gorm.DB, id uint, status string, errorMsg string) error {
	return db.Model(&Submission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"error_msg": errorMsg,
	}).Error
}

func GetClassLeaderboard(db *gorm.DB, classID string) ([]struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Solved    int64  `json:"solved"`
	Submitted int64  `json:"submitted"`
}, error) {
	var results []struct {
		UserID    uint   `json:"user_id"`
		Username  string `json:"username"`
		Solved    int64  `json:"solved"`
		Submitted int64  `json:"submitted"`
	}

	// Get all users in the class
	var users []User
	if err := db.Where("class = ?", classID).Find(&users).Error; err != nil {
		return nil, err
	}

	// Get submissions for each user
	for _, user := range users {
		var solved, submitted int64
		if err := db.Model(&Submission{}).Where("user_id = ? AND status = ?", user.ID, "Accepted").Count(&solved).Error; err != nil {
			return nil, err
		}
		if err := db.Model(&Submission{}).Where("user_id = ?", user.ID).Count(&submitted).Error; err != nil {
			return nil, err
		}

		results = append(results, struct {
			UserID    uint   `json:"user_id"`
			Username  string `json:"username"`
			Solved    int64  `json:"solved"`
			Submitted int64  `json:"submitted"`
		}{
			UserID:    user.ID,
			Username:  user.Username,
			Solved:    solved,
			Submitted: submitted,
		})
	}

	return results, nil
}
