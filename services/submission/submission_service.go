package submissionServ

import (
	"context"
	"github.com/google/uuid"
	submissionModel "neptune/backend/models/submission"
	"neptune/backend/pkg/requests"
)

type SubmissionService interface {
	SubmitCode(ctx context.Context, request *requests.SubmitCodeRequest, userID uuid.UUID) (*submissionModel.Submission, error)
	StartListeners() error
}
