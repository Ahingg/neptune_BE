package caseService

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	caseModel "neptune/backend/models/contest"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	caseRepository "neptune/backend/repositories/case"
)

type caseServiceImpl struct {
	caseRepo caseRepository.CaseRepository
}

func NewCaseService(caseRepo caseRepository.CaseRepository) CaseService {
	return &caseServiceImpl{caseRepo: caseRepo}
}

// CreateCase creates a new problem case.
func (s *caseServiceImpl) CreateCase(ctx context.Context, req requests.CreateCaseRequest) (*responses.CaseResponse, error) {
	problemCase := &caseModel.Case{
		ID:            uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		PDFFileUrl:    req.PDFFileUrl,
		TimeLimitMs:   req.TimeLimitMs,
		MemoryLimitMb: req.MemoryLimitMb,
	}
	if err := s.caseRepo.SaveCase(ctx, problemCase); err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}
	return &responses.CaseResponse{
		ID:            problemCase.ID,
		Name:          problemCase.Name,
		Description:   problemCase.Description,
		PDFFileUrl:    problemCase.PDFFileUrl,
		TimeLimitMs:   problemCase.TimeLimitMs,
		MemoryLimitMb: problemCase.MemoryLimitMb,
		CreatedAt:     problemCase.CreatedAt,
	}, nil
}

// GetCaseByID retrieves a problem case.
func (s *caseServiceImpl) GetCaseByID(ctx context.Context, caseID uuid.UUID) (*responses.CaseResponse, error) {
	problemCase, err := s.caseRepo.FindCaseByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}
	if problemCase == nil {
		return nil, nil // Not found
	}
	return &responses.CaseResponse{
		ID:            problemCase.ID,
		Name:          problemCase.Name,
		Description:   problemCase.Description,
		PDFFileUrl:    problemCase.PDFFileUrl,
		TimeLimitMs:   problemCase.TimeLimitMs,
		MemoryLimitMb: problemCase.MemoryLimitMb,
		CreatedAt:     problemCase.CreatedAt,
		UpdatedAt:     problemCase.UpdatedAt,
	}, nil
}

// GetAllCases retrieves all problem cases.
func (s *caseServiceImpl) GetAllCases(ctx context.Context) ([]responses.CaseResponse, error) {
	cases, err := s.caseRepo.FindAllCases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all cases: %w", err)
	}

	resp := make([]responses.CaseResponse, len(cases))
	for i, c := range cases {
		resp[i] = responses.CaseResponse{
			ID:            c.ID,
			Name:          c.Name,
			Description:   c.Description,
			PDFFileUrl:    c.PDFFileUrl,
			TimeLimitMs:   c.TimeLimitMs,
			MemoryLimitMb: c.MemoryLimitMb,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
		}
	}
	return resp, nil
}

// UpdateCase updates an existing problem case.
func (s *caseServiceImpl) UpdateCase(ctx context.Context, caseID uuid.UUID, req requests.UpdateCaseRequest) (*responses.CaseResponse, error) {
	problemCase, err := s.caseRepo.FindCaseByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to find case for update: %w", err)
	}
	if problemCase == nil {
		return nil, fmt.Errorf("case with ID %s not found", caseID.String())
	}

	problemCase.Name = req.Name
	problemCase.Description = req.Description
	problemCase.PDFFileUrl = req.PDFFileUrl
	problemCase.TimeLimitMs = req.TimeLimitMs
	problemCase.MemoryLimitMb = req.MemoryLimitMb

	if err := s.caseRepo.SaveCase(ctx, problemCase); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}

	return &responses.CaseResponse{
		ID:            problemCase.ID,
		Name:          problemCase.Name,
		Description:   problemCase.Description,
		PDFFileUrl:    problemCase.PDFFileUrl,
		TimeLimitMs:   problemCase.TimeLimitMs,
		MemoryLimitMb: problemCase.MemoryLimitMb,
		CreatedAt:     problemCase.CreatedAt,
		UpdatedAt:     problemCase.UpdatedAt,
	}, nil
}

// DeleteCase soft deletes a problem case.
func (s *caseServiceImpl) DeleteCase(ctx context.Context, caseID uuid.UUID) error {
	// Implement checks if case is part of any active contests before deleting
	// For now, just soft delete
	return s.caseRepo.DeleteCase(ctx, caseID)
}
