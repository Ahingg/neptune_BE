package contestService

import (
	"context"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type ContestService interface {
	// Global Contest Management
	FindAllActiveGlobalContests(ctx context.Context) ([]responses.GlobalContestResponse, error)
	FindAllActiveGlobalContestsDetail(ctx context.Context) ([]responses.GlobalContestDetailResponse, error)

	// Contest Management
	CreateContest(ctx context.Context, req requests.CreateContestRequest) (*responses.ContestResponse, error)
	GetContestByID(ctx context.Context, contestID uuid.UUID) (*responses.ContestDetailResponse, error)
	GetAllContests(ctx context.Context) ([]responses.ContestResponse, error)
	UpdateContest(ctx context.Context, contestID uuid.UUID, req requests.UpdateContestRequest) (*responses.ContestResponse, error)
	DeleteContest(ctx context.Context, contestID uuid.UUID) error

	// Contest-Case (Problem) Management
	AddCasesToContest(ctx context.Context, contestID uuid.UUID, req requests.AddCasesToContestRequest) error
	GetContestCases(ctx context.Context, contestID uuid.UUID) ([]responses.ContestCaseResponse, error)
	GetContentCaseByCaseID(ctx context.Context, contestID, caseID uuid.UUID) (*responses.ContestCaseResponse, error)

	// Class-Contest Assignment
	AssignContestToClass(ctx context.Context, classTransactionID uuid.UUID, req requests.AssignContestToClassRequest) (*responses.ClassContestAssignmentResponse, error)
	GetContestsForClass(ctx context.Context, classTransactionID uuid.UUID) ([]responses.ClassContestAssignmentResponse, error)
	RemoveContestFromClass(ctx context.Context, classTransactionID, contestID uuid.UUID) error
}
