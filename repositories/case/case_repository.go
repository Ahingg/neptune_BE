package caseRepository

import (
	"context"
	"github.com/google/uuid"
	caseModel "neptune/backend/models/contest"
)

type CaseRepository interface {
	SaveCase(ctx context.Context, problemCase *caseModel.Case) error
	FindCaseByID(ctx context.Context, caseID uuid.UUID) (*caseModel.Case, error)
	FindAllCases(ctx context.Context) ([]caseModel.Case, error)
	DeleteCase(ctx context.Context, caseID uuid.UUID) error // Soft delete
}
