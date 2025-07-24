package submissionRepo

import (
	"context"
	"github.com/google/uuid"
	submissionModel "neptune/backend/models/submission"
	"time"
)

type SubmissionRepository interface {
	Save(ctx context.Context, submission *submissionModel.Submission) error
	FindByID(ctx context.Context, id string) (*submissionModel.Submission, error)
	Update(ctx context.Context, submission *submissionModel.Submission) error
	SaveResultsBatch(ctx context.Context, results []submissionModel.SubmissionResult) error
	FindAllForContest(ctx context.Context, caseIDs []uuid.UUID, userIDs []uuid.UUID, contestStartTime time.Time) ([]submissionModel.Submission, error)
	FindByUserInContest(ctx context.Context, contestID uuid.UUID, userID uuid.UUID, classID *uuid.UUID) ([]submissionModel.Submission, error)
	FindClassSubmissions(ctx context.Context, classID uuid.UUID, contestID uuid.UUID) ([]submissionModel.Submission, error)
}
