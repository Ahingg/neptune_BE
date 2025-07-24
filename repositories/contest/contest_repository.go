package contestRepository

import (
	"context"
	"github.com/google/uuid"
	contestModel "neptune/backend/models/contest"
)

type ContestRepository interface {
	// Global Contest Management
	SaveGlobalContestDetail(ctx context.Context, detail *contestModel.GlobalContestDetail) error
	FindAllActiveGlobalContests(ctx context.Context) ([]contestModel.Contest, error)

	SaveContest(ctx context.Context, contest *contestModel.Contest) error
	FindContestByID(ctx context.Context, contestID uuid.UUID) (*contestModel.Contest, error)
	FindAllContests(ctx context.Context) ([]contestModel.Contest, error)
	DeleteContest(ctx context.Context, contestID uuid.UUID) error // Soft delete

	// ContestCase (Problems in a Contest) Management
	AddCasesToContest(ctx context.Context, contestID uuid.UUID, cases []contestModel.ContestCase) error
	ClearContestCases(ctx context.Context, contestID uuid.UUID) error // Clear all problems for a contest
	FindContestCases(ctx context.Context, contestID uuid.UUID) ([]contestModel.ContestCase, error)
	GetCaseCountInContest(ctx context.Context, contestID uuid.UUID) (int, error)
	GetContestCaseByCaseID(ctx context.Context, contestID, caseID uuid.UUID) (*contestModel.ContestCase, error)

	// ClassContest (Contest Assignment to Class) Management
	AssignContestToClass(ctx context.Context, classContest *contestModel.ClassContest) error
	FindContestsByClassTransactionID(ctx context.Context, classTransactionID uuid.UUID) ([]contestModel.ClassContest, error)
	FindClassContestByIDs(ctx context.Context, classTransactionID, contestID uuid.UUID) (*contestModel.ClassContest, error)
	RemoveContestFromClass(ctx context.Context, classTransactionID, contestID uuid.UUID) error
}
