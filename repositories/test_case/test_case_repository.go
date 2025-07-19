package testCaseRepo

import (
	"context"
	testCaseModel "neptune/backend/models/test_case"
)

type TestCaseRepository interface {
	SaveTestCase(ctx context.Context, testCase *testCaseModel.TestCase) error
	SaveTestCaseBatch(ctx context.Context, testCases []testCaseModel.TestCase) error
	FindTestCaseByCaseID(ctx context.Context, caseID string) ([]testCaseModel.TestCase, error)
	DeleteTestCaseByCaseID(ctx context.Context, caseID string) error // Hard Delete
}
