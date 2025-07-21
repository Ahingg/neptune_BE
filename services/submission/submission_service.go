package submissionServ

import (
	"context"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type SubmissionService interface {
	SubmitCode(ctx context.Context, request *requests.SubmitCodeRequest, userID uuid.UUID) (*responses.SubmitCodeResponse, error)
	StartListeners() error
}
