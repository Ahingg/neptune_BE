package caseRepository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	caseModel "neptune/backend/models/contest"
	"time"
)

type caseRepositoryImpl struct {
	db *gorm.DB
}

func NewCaseRepository(db *gorm.DB) CaseRepository {
	return &caseRepositoryImpl{db: db}
}

// SaveCase creates or updates a Case.
func (r *caseRepositoryImpl) SaveCase(ctx context.Context, problemCase *caseModel.Case) error {
	if problemCase.ID == uuid.Nil {
		problemCase.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // Conflict on primary key (ID)
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":            problemCase.Name,
			"description":     problemCase.Description,
			"pdf_file_url":    problemCase.PDFFileUrl,
			"time_limit_ms":   problemCase.TimeLimitMs,
			"memory_limit_mb": problemCase.MemoryLimitMb,
			"updated_at":      time.Now(),
		}),
	}).Create(problemCase).Error
}

// FindCaseByID retrieves a Case.
func (r *caseRepositoryImpl) FindCaseByID(ctx context.Context, caseID uuid.UUID) (*caseModel.Case, error) {
	var problemCase caseModel.Case
	result := r.db.WithContext(ctx).Where("id = ?", caseID).First(&problemCase)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find case by ID %s: %w", caseID.String(), result.Error)
	}
	return &problemCase, nil
}

// FindAllCases retrieves all Cases.
func (r *caseRepositoryImpl) FindAllCases(ctx context.Context) ([]caseModel.Case, error) {
	var cases []caseModel.Case
	result := r.db.WithContext(ctx).Find(&cases)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all cases: %w", result.Error)
	}
	return cases, nil
}

// DeleteCase soft deletes a case.
func (r *caseRepositoryImpl) DeleteCase(ctx context.Context, caseID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&caseModel.Case{}, caseID).Error
}
