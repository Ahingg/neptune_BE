package semester

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	model "neptune/backend/models/semester"
)

type semesterRepository struct {
	db *gorm.DB
}

func (s *semesterRepository) Save(ctx context.Context, semester *model.Semester) error {
	// Use Upsert logic
	err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // Conflict on external ID
		DoUpdates: clause.Assignments(map[string]interface{}{ // Update these fields on conflict
			"description": semester.Description,
			"start":       semester.Start,
			"end":         semester.End,
			// gorm.Model's UpdatedAt will be handled automatically
		}),
	}).Create(semester).Error // Use Create method with Clauses
	if err != nil {
		return fmt.Errorf("failed to save or update internal_semester %s (%s): %w", semester.Description, semester.ID, err)
	}
	return nil
}

func (s *semesterRepository) FindAll(ctx context.Context) ([]model.Semester, error) {
	var semesters []model.Semester
	result := s.db.WithContext(ctx).Find(&semesters)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve all semesters: %w", result.Error)
	}
	return semesters, nil
}

func (s *semesterRepository) GetSemesterByID(ctx context.Context, semesterID string) (model.Semester, error) {
	var semester model.Semester
	result := s.db.WithContext(ctx).Where("id = ?", semesterID).First(&semester)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Semester{}, fmt.Errorf("internal_semester with ID %s not found", semesterID)
		}
		return model.Semester{}, fmt.Errorf("failed to retrieve internal_semester by ID %s: %w", semesterID, result.Error)
	}
	return semester, nil
}

func NewSemesterRepository(db *gorm.DB) SemesterRepository {
	return &semesterRepository{db: db}
}
