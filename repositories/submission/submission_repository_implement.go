package submissionRepo

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	submissionModel "neptune/backend/models/submission"
	"time"
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

func (r submissionRepository) FindAllForContest(ctx context.Context, caseIDs []uuid.UUID, userIDs []uuid.UUID, contestStartTime time.Time) ([]submissionModel.Submission, error) {
	var submissions []submissionModel.Submission
	err := r.db.WithContext(ctx).
		Where("case_id IN ?", caseIDs).
		Where("user_id IN ?", userIDs).
		Where("created_at >= ?", contestStartTime).
		Order("created_at asc"). // IMPORTANT: Sort by time to process chronologically
		Find(&submissions).Error
	return submissions, err
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}
