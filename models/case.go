package models

import (
	"time"

	"gorm.io/gorm"
)

type Case struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Header      string     `gorm:"not null" json:"header"` // Group cases together
	IsPublic    bool       `gorm:"default:false" json:"is_public"`
	ClassID     string     `json:"class_id"` // If case is specific to a class
	TestCases   []TestCase `json:"test_cases"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TestCase struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	CaseID   uint   `json:"case_id"`
	Input    string `gorm:"not null" json:"input"`
	Output   string `gorm:"not null" json:"output"`
	IsHidden bool   `gorm:"default:false" json:"is_hidden"`
}

func CreateCase(db *gorm.DB, case_ *Case) error {
	return db.Create(case_).Error
}

func GetCase(db *gorm.DB, id uint) (*Case, error) {
	var case_ Case
	if err := db.Preload("TestCases").First(&case_, id).Error; err != nil {
		return nil, err
	}
	return &case_, nil
}

func GetCases(db *gorm.DB, classID string, isPublic bool) ([]Case, error) {
	var cases []Case
	query := db.Preload("TestCases")

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if isPublic {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Find(&cases).Error; err != nil {
		return nil, err
	}
	return cases, nil
}

func UpdateCase(db *gorm.DB, case_ *Case) error {
	return db.Save(case_).Error
}

func DeleteCase(db *gorm.DB, id uint) error {
	return db.Delete(&Case{}, id).Error
}

func AddTestCase(db *gorm.DB, testCase *TestCase) error {
	return db.Create(testCase).Error
}

func DeleteTestCase(db *gorm.DB, id uint) error {
	return db.Delete(&TestCase{}, id).Error
}
