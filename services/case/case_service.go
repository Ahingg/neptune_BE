package caseService

import (
	"context"
	"github.com/google/uuid"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type CaseService interface {
	CreateCase(ctx context.Context, req requests.CreateCaseRequest, url string) (*responses.CaseResponse, error)
	GetCaseByID(ctx context.Context, caseID uuid.UUID) (*responses.CaseResponse, error)
	GetAllCases(ctx context.Context) ([]responses.CaseResponse, error)
	UpdateCase(ctx context.Context, caseID uuid.UUID, req requests.UpdateCaseRequest) (*responses.CaseResponse, error)
	DeleteCase(ctx context.Context, caseID uuid.UUID) error // Soft delete
}
