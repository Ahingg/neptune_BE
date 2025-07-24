package testCaseServ

import (
	"context"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type TestCaseService interface {
	UploadTestCases(ctx context.Context, req requests.AddTestCaseRequest) error
	GetTestCasesByCaseID(ctx context.Context, caseID string) ([]responses.TestCaseResponse, error)
}
