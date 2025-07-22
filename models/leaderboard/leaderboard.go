package leaderboardModel

import "github.com/google/uuid"

type CaseResult struct {
	CaseID           uuid.UUID `json:"case_id"`            // Unique identifier for the case
	SubmissionID     uuid.UUID `json:"submission_id"`      // Unique identifier for the submission
	Status           string    `json:"status"`             // Status of the case result (Accepted, Wrong Answer, No Submission)
	Score            int       `json:"score"`              // Score for the case, if applicable
	IsSolved         bool      `json:"is_solved"`          // Whether the case is solved
	SolveTimeMinutes int       `json:"solve_time_minutes"` // Time taken to solve the case in minutes
	WrongAttempts    int       `json:"wrong_attempts"`     // Number of wrong attempts before solving
}

type LeaderboardRow struct {
	Rank           int                   `json:"rank"`
	UserID         uuid.UUID             `json:"user_id"`
	UserName       string                `json:"username"`
	Name           string                `json:"name"`
	SolvedCount    int                   `json:"solved_count"`
	TotalPenalty   int                   `json:"total_penalty"`
	ProblemResults map[string]CaseResult `json:"case_results"`
}
