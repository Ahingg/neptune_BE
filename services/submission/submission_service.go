package submissionServ

import (
	"context"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type SubmissionService interface {
	SubmitCode(ctx context.Context, request *requests.SubmitCodeRequest, userID uuid.UUID) (*responses.SubmitCodeResponse, error)
	GetSubmissionByUserInContest(ctx context.Context, userID uuid.UUID, contestID uuid.UUID, classTransactionID *uuid.UUID) ([]responses.GetSubmissionByUserInContestResponse, error)
	GetClassContestSubmissions(ctx context.Context, classTransactionID uuid.UUID, contestID uuid.UUID) ([]responses.GetSubmissionByClassContestResponse, error)
	StartListeners() error
}
