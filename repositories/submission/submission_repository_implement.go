package submissionRepo

import (
	"context"
	"gorm.io/gorm"
	submissionModel "neptune/backend/models/submission"
)

type submissionRepository struct {
	db *gorm.DB
}

func (r submissionRepository) Save(ctx context.Context, submission *submissionModel.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r submissionRepository) FindByID(ctx context.Context, id string) (*submissionModel.Submission, error) {
	var submission submissionModel.Submission
	err := r.db.WithContext(ctx).Preload("SubmissionResults").First(&submission, "id = ?", id).Error
	return &submission, err
}

func (r submissionRepository) Update(ctx context.Context, submission *submissionModel.Submission) error {
	return r.db.WithContext(ctx).Save(submission).Error
}

func (r submissionRepository) SaveResultsBatch(ctx context.Context, results []submissionModel.SubmissionResult) error {
	if len(results) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&results).Error
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}
