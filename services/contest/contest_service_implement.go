package contestService

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	contestModel "neptune/backend/models/contest"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	"neptune/backend/pkg/utils"
	caseRepository "neptune/backend/repositories/case"
	contestRepository "neptune/backend/repositories/contest"
)

type contestServiceImpl struct {
	contestRepo contestRepository.ContestRepository
	caseRepo    caseRepository.CaseRepository // Need to lookup cases by ID
}

func NewContestService(contestRepo contestRepository.ContestRepository, caseRepo caseRepository.CaseRepository) ContestService {
	return &contestServiceImpl{
		contestRepo: contestRepo,
		caseRepo:    caseRepo,
	}
}

// CreateContest creates a new contest.
func (s *contestServiceImpl) CreateContest(ctx context.Context, req requests.CreateContestRequest) (*responses.ContestResponse, error) {
	contest := &contestModel.Contest{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.contestRepo.SaveContest(ctx, contest); err != nil {
		return nil, fmt.Errorf("failed to create contest: %w", err)
	}
	return &responses.ContestResponse{
		ID:          contest.ID,
		Name:        contest.Name,
		Scope:       contest.Scope,
		Description: contest.Description,
		CreatedAt:   contest.CreatedAt,
	}, nil
}

// GetContestByID retrieves a contest with its associated cases.
func (s *contestServiceImpl) GetContestByID(ctx context.Context, contestID uuid.UUID) (*responses.ContestDetailResponse, error) {
	contest, err := s.contestRepo.FindContestByID(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest: %w", err)
	}
	if contest == nil {
		return nil, nil // Not found
	}

	resp := &responses.ContestDetailResponse{
		ID:          contest.ID,
		Name:        contest.Name,
		Scope:       contest.Scope,
		Description: contest.Description,
		CreatedAt:   contest.CreatedAt,
	}

	for _, cc := range contest.ContestCases {
		if cc.Case.ID != uuid.Nil { // Ensure case was loaded
			resp.Cases = append(resp.Cases, responses.ContestCaseProblemResponse{
				CaseID:        cc.Case.ID,
				ProblemCode:   cc.ProblemCode,
				Name:          cc.Case.Name,
				TimeLimitMs:   cc.Case.TimeLimitMs,
				MemoryLimitMb: cc.Case.MemoryLimitMb,
				PDFFileUrl:    cc.Case.PDFFileUrl,
			})
		}
	}
	return resp, nil
}

// GetAllContests retrieves all contests (basic info).
func (s *contestServiceImpl) GetAllContests(ctx context.Context) ([]responses.ContestResponse, error) {
	contests, err := s.contestRepo.FindAllContests(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all contests: %w", err)
	}

	resp := make([]responses.ContestResponse, len(contests))
	for i, c := range contests {
		resp[i] = responses.ContestResponse{
			ID:          c.ID,
			Name:        c.Name,
			Scope:       c.Scope,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
	}
	return resp, nil
}

// UpdateContest updates an existing contest.
func (s *contestServiceImpl) UpdateContest(ctx context.Context, contestID uuid.UUID, req requests.UpdateContestRequest) (*responses.ContestResponse, error) {
	contest, err := s.contestRepo.FindContestByID(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to find contest for update: %w", err)
	}
	if contest == nil {
		return nil, fmt.Errorf("contest with ID %s not found", contestID.String())
	}

	contest.Name = req.Name
	contest.Description = req.Description

	if err := s.contestRepo.SaveContest(ctx, contest); err != nil {
		return nil, fmt.Errorf("failed to update contest: %w", err)
	}

	return &responses.ContestResponse{
		ID:          contest.ID,
		Name:        contest.Name,
		Scope:       contest.Scope,
		Description: contest.Description,
		CreatedAt:   contest.CreatedAt,
		UpdatedAt:   contest.UpdatedAt,
	}, nil
}

// DeleteContest soft deletes a contest.
func (s *contestServiceImpl) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	// Implement checks if contest is assigned to any active classes before deleting
	// For now, just soft delete
	return s.contestRepo.DeleteContest(ctx, contestID)
}

// AddCasesToContest adds/updates problems for a contest.
func (s *contestServiceImpl) AddCasesToContest(ctx context.Context, contestID uuid.UUID, req requests.AddCasesToContestRequest) error {
	var contestCases []contestModel.ContestCase
	for _, problem := range req.Problems {
		// Cek apakah case nya beneran ada
		_, err := s.caseRepo.FindCaseByID(ctx, problem.CaseID)
		if err != nil {
			return fmt.Errorf("failed to find case by ID %s: %w", problem.CaseID.String(), err)
		}

		// If a problem code is duplicated, the unique index on (ContestID, CaseID) will prevent it
		//caseCode := getProblemCodeByContestCaseCount
		caseCode, err := s.getProblemCodeByContestCaseCount(contestID)
		if err != nil {
			return fmt.Errorf("failed to generate problem code: %w", err)
		}

		contestCases = append(contestCases, contestModel.ContestCase{
			ContestID:   contestID,
			CaseID:      problem.CaseID,
			ProblemCode: caseCode,
		})
	}

	if err := s.contestRepo.AddCasesToContest(ctx, contestID, contestCases); err != nil {
		return fmt.Errorf("failed to add cases to contest: %w", err)
	}
	return nil
}

func (s *contestServiceImpl) GetContestCases(ctx context.Context, contestID uuid.UUID) ([]responses.ContestCaseResponse, error) {
	cases, err := s.contestRepo.FindContestCases(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contest cases: %w", err)
	}

	resp := make([]responses.ContestCaseResponse, len(cases))
	for i, cc := range cases {
		resp[i] = responses.ContestCaseResponse{
			CaseID:   cc.Case.ID,
			CaseCode: cc.ProblemCode,
			CaseName: cc.Case.Name,
		}
	}
	return resp, nil
}

// AssignContestToClass assigns a contest to a specific class with a duration.
func (s *contestServiceImpl) AssignContestToClass(ctx context.Context, classTransactionID uuid.UUID, req requests.AssignContestToClassRequest) (*responses.ClassContestAssignmentResponse, error) {
	// Optional: Verify ContestID and ClassTransactionID exist before assigning
	// contest, err := s.contestRepo.FindContestByID(ctx, req.ContestID)
	// if err != nil || contest == nil { return nil, fmt.Errorf("contest not found") }
	// class, err := s.classRepo.FindClassByTransactionID(ctx, classTransactionID) // Assuming ClassRepo is available
	// if err != nil || class == nil { return nil, fmt.Errorf("class not found") }

	classContest := &contestModel.ClassContest{
		ClassTransactionID: classTransactionID,
		ContestID:          req.ContestID,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
	}

	if err := s.contestRepo.AssignContestToClass(ctx, classContest); err != nil {
		return nil, fmt.Errorf("failed to assign contest to class: %w", err)
	}

	return &responses.ClassContestAssignmentResponse{
		ClassTransactionID: classContest.ClassTransactionID,
		ContestID:          classContest.ContestID,
		StartTime:          classContest.StartTime,
		EndTime:            classContest.EndTime,
		CreatedAt:          classContest.CreatedAt,
		UpdatedAt:          classContest.UpdatedAt,
	}, nil
}

// GetContestsForClass retrieves all contests assigned to a specific class.
func (s *contestServiceImpl) GetContestsForClass(ctx context.Context, classTransactionID uuid.UUID) ([]responses.ClassContestAssignmentResponse, error) {
	classContests, err := s.contestRepo.FindContestsByClassTransactionID(ctx, classTransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contests for class: %w", err)
	}

	resp := make([]responses.ClassContestAssignmentResponse, len(classContests))
	for i, cc := range classContests {

		resp[i] = responses.ClassContestAssignmentResponse{
			ClassTransactionID: cc.ClassTransactionID,
			ContestID:          cc.ContestID,
			StartTime:          cc.StartTime,
			EndTime:            cc.EndTime,
			CreatedAt:          cc.CreatedAt,
			UpdatedAt:          cc.UpdatedAt,
			Contest: responses.ContestResponse{
				ID:          cc.Contest.ID,
				Name:        cc.Contest.Name,
				Scope:       cc.Contest.Scope,
				Description: cc.Contest.Description,
				CreatedAt:   cc.Contest.CreatedAt,
			},
		}
	}
	return resp, nil
}

func (s *contestServiceImpl) getProblemCodeByContestCaseCount(contestID uuid.UUID) (string, error) {
	// Get all contest cases for the given contest ID
	caseCount, err := s.contestRepo.GetCaseCountInContest(context.Background(), contestID)
	if err != nil {
		// Handle error (you might want to panic or return a default)
		return "", err
	}

	return utils.GenerateAlphabetCode(caseCount), nil
}

// RemoveContestFromClass removes a contest assignment from a class.
func (s *contestServiceImpl) RemoveContestFromClass(ctx context.Context, classTransactionID, contestID uuid.UUID) error {
	return s.contestRepo.RemoveContestFromClass(ctx, classTransactionID, contestID)
}
