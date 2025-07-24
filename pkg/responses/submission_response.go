package responses

type GetSubmissionPerContestResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	SubmissionID string `json:"submission_id"`
	ContestID    string `json:"contest_id"`
	CaseID       string `json:"case_id"`
	CaseCode     string `json:"case_code"`
	Status       string `json:"status"`
	Score        int    `json:"score"`
	SubmitTime   string `json:"submit_time"` // Time when the submission was made
	LanguageID   int    `json:"language_id"` // Name of the programming language used

}

type GetUserSubmissionsResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	SubmissionID string `json:"submission_id"`
	ContestID    string `json:"contest_id"`
	CaseID       string `json:"case_id"`
	CaseCode     string `json:"case_code"`
	Status       string `json:"status"`
	Score        int    `json:"score"`
	SubmitTime   string `json:"submit_time"` // Time when the submission was made
	LanguageID   int    `json:"language_id"` // Name of the programming language used
}
