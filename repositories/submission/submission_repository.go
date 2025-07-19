package submissionRepo

import (
	"context"
	submissionModel "neptune/backend/models/submission"
)

type SubmissionRepository interface {
	Save(ctx context.Context, submission *submissionModel.Submission) error
	FindByID(ctx context.Context, id string) (*submissionModel.Submission, error)
	Update(ctx context.Context, submission *submissionModel.Submission) error
	SaveResultsBatch(ctx context.Context, results []*submissionModel.SubmissionResult) error
}
