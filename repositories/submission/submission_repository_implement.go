package submissionRepo

import (
	"context"
	"fmt"
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

func (r submissionRepository) FindByUserInContest(ctx context.Context, contestID uuid.UUID, userID uuid.UUID, classID *uuid.UUID) ([]submissionModel.Submission, error) {
	var submissions []submissionModel.Submission
	classQuery := "class_transaction_id IS NOT NULL"
	if classID == nil {
		classQuery = "class_transaction_id IS NULL"
	}
	fmt.Println(contestID, userID, classID)
	err := r.db.WithContext(ctx).
		Where("contest_id = ?", contestID).
		Where("user_id = ?", userID).Where(classQuery).Find(&submissions).Error
	return submissions, err
}

func (r submissionRepository) FindClassSubmissions(ctx context.Context, classID uuid.UUID, contestID uuid.UUID) ([]submissionModel.Submission, error) {
	var submissions []submissionModel.Submission
	err := r.db.WithContext(ctx).
		Where("class_transaction_id = ?", classID).
		Where("contest_id = ?", contestID).
		Order("created_at asc"). // Sort by time to process chronologically
		Find(&submissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find class submissions for class %s and contest %s: %w", classID.String(), contestID.String(), err)

	}
	return submissions, nil
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}
