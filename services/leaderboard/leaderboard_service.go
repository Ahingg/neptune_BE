package leaderboardServ

import (
	"context"
	"fmt"
	leaderboardModel "neptune/backend/models/leaderboard"
	submissionModel "neptune/backend/models/submission"
	"neptune/backend/repositories/class"
	contestRepository "neptune/backend/repositories/contest"
	submissionRepo "neptune/backend/repositories/submission"
	userRepo "neptune/backend/repositories/user"
	"sort"
	"time"

	"github.com/google/uuid"
)

const penaltyPerWrongAttempt = 20

type Service interface {
	GetContestLeaderboard(ctx context.Context, classID, contestID uuid.UUID) ([]leaderboardModel.LeaderboardRow, error)
}

type serviceImpl struct {
	submissionRepo submissionRepo.SubmissionRepository
	contestRepo    contestRepository.ContestRepository
	classRepo      class.ClassRepository
	userRepo       userRepo.UserRepository
}

// NewService initializes the leaderboard service with required repositories.
func NewService(subRepo submissionRepo.SubmissionRepository,
	conRepo contestRepository.ContestRepository,
	clsRepo class.ClassRepository,
	userRepository userRepo.UserRepository) Service {
	return &serviceImpl{
		submissionRepo: subRepo,
		contestRepo:    conRepo,
		classRepo:      clsRepo,
		userRepo:       userRepository,
	}
}

func (s *serviceImpl) GetContestLeaderboard(ctx context.Context, classID, contestID uuid.UUID) ([]leaderboardModel.LeaderboardRow, error) {
	// Step 1: Fetch core contest data using the provided repositories
	classContest, err := s.contestRepo.FindClassContestByIDs(ctx, classID, contestID)
	if err != nil {
		return nil, fmt.Errorf("could not find contest assignment for this class: %w", err)
	}
	contestStartTime := classContest.StartTime

	contestCases, err := s.contestRepo.FindContestCases(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch problems for contest: %w", err)
	}
	if len(contestCases) == 0 {
		return []leaderboardModel.LeaderboardRow{}, nil // Return empty leaderboard if no problems
	}

	classInfo, err := s.classRepo.FindClassByTransactionID(ctx, classID.String())
	if err != nil {
		return nil, fmt.Errorf("could not fetch class details: %w", err)
	}

	// Step 2: Aggregate participants and problem IDs for the main query
	caseIDs := make([]uuid.UUID, len(contestCases))
	for i, cc := range contestCases {
		caseIDs[i] = cc.CaseID
	}

	userIDs := make([]uuid.UUID, len(classInfo.Students))
	// We also need a quick way to look up user names later
	userNameMap := make(map[uuid.UUID]struct {
		Name     string
		UserName string
	})
	for i, student := range classInfo.Students {
		userIDs[i] = student.UserID
		userNameMap[student.UserID] = struct {
			Name     string
			UserName string
		}{
			Name:     student.User.Name,
			UserName: student.User.Username,
		}
	}

	// Step 3: Fetch all relevant submissions in a single, efficient query
	allSubmissions, err := s.submissionRepo.FindAllForContest(ctx, caseIDs, userIDs, contestStartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}

	// Group submissions by user for easier processing
	submissionsByUser := make(map[uuid.UUID][]submissionModel.Submission)
	for _, sub := range allSubmissions {
		submissionsByUser[sub.UserID] = append(submissionsByUser[sub.UserID], sub)
	}

	// Step 4: Process submissions for each student to calculate results
	var leaderboardRows []leaderboardModel.LeaderboardRow
	for _, userID := range userIDs {
		userSubmissions := submissionsByUser[userID]

		row := leaderboardModel.LeaderboardRow{
			UserID:         userID,
			Name:           userNameMap[userID].Name,
			UserName:       userNameMap[userID].UserName,
			SolvedCount:    0,
			TotalPenalty:   0,
			ProblemResults: make(map[string]leaderboardModel.CaseResult),
		}

		submissionsByCase := make(map[uuid.UUID][]submissionModel.Submission)
		for _, sub := range userSubmissions {
			submissionsByCase[sub.CaseID] = append(submissionsByCase[sub.CaseID], sub)
		}

		for _, problem := range contestCases {
			problemSubmissions := submissionsByCase[problem.CaseID]
			problemResult := calculateProblemResult(problemSubmissions, contestStartTime)
			row.ProblemResults[problem.ProblemCode] = problemResult // Use ProblemCode like "A", "B"

			if problemResult.IsSolved {

				row.SolvedCount++
				row.TotalPenalty += problemResult.SolveTimeMinutes + (problemResult.WrongAttempts * penaltyPerWrongAttempt)
			}
		}
		leaderboardRows = append(leaderboardRows, row)
	}

	// Step 5: Sort the leaderboard based on ICPC rules
	sort.Slice(leaderboardRows, func(i, j int) bool {
		if leaderboardRows[i].SolvedCount != leaderboardRows[j].SolvedCount {
			return leaderboardRows[i].SolvedCount > leaderboardRows[j].SolvedCount
		}
		if leaderboardRows[i].TotalPenalty != leaderboardRows[j].TotalPenalty {
			return leaderboardRows[i].TotalPenalty < leaderboardRows[j].TotalPenalty
		}
		return leaderboardRows[i].UserName < leaderboardRows[j].UserName
	})

	// Step 6: Assign final ranks
	for i := range leaderboardRows {
		leaderboardRows[i].Rank = i + 1
	}

	return leaderboardRows, nil
}

// calculateProblemResult logic remains the same as it's pure computation.
func calculateProblemResult(submissions []submissionModel.Submission, contestStartTime time.Time) leaderboardModel.CaseResult {
	if len(submissions) == 0 {
		return leaderboardModel.CaseResult{Status: "Unsolved"}
	}

	wrongAttempts := 0
	for _, sub := range submissions {
		wrongAttempts++
		if sub.Status == submissionModel.SubmissionStatusAccepted {
			solveTime := int(sub.CreatedAt.Sub(contestStartTime).Minutes())
			return leaderboardModel.CaseResult{
				SubmissionID:     sub.ID,
				CaseID:           sub.CaseID,
				Status:           "AC",
				Score:            sub.Score,
				IsSolved:         true,
				SolveTimeMinutes: solveTime,
				WrongAttempts:    wrongAttempts,
			}
		}

	}

	wrongResult := leaderboardModel.CaseResult{
		Status:        "WA", // Or status of the last attempt
		Score:         0,
		IsSolved:      false,
		WrongAttempts: wrongAttempts,
	}
	latestUpdateTime := submissions[0].UpdatedAt

	for _, sub := range submissions {
		if sub.Status == submissionModel.SubmissionStatusWrongAnswer {
			afterOrEqual := sub.UpdatedAt.Equal(latestUpdateTime) || sub.UpdatedAt.After(latestUpdateTime)
			if sub.Score > wrongResult.Score && afterOrEqual {
				wrongResult.SubmissionID = sub.ID
				wrongResult.Score = sub.Score
				latestUpdateTime = sub.UpdatedAt
				wrongResult.CaseID = sub.CaseID
			}
		}
	}

	return wrongResult
}
