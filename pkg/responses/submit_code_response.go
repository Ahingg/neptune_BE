package responses

type SubmitCodeResponse struct {
	SubmissionID string `json:"submission_id"` // UUID of the submission
	Status       string `json:"status"`        // ID of the programming language used
}

type FinalResultResponse struct {
	SubmissionID string                  `json:"submission_id"` // UUID of the submission
	FinalStatus  string                  `json:"final_status"`  // Final status of the submission
	Score        int                     `json:"score"`         // Score of the submission
	TestCases    []TestCaseJudgeResponse `json:"testcases"`     // List of test case results
}

type TestCaseJudgeResponse struct {
	Number         int    `json:"number"`          // Test case number
	Verdict        string `json:"verdict"`         // Verdict of the test case
	Input          string `json:"input"`           // Input for the test case
	ExpectedOutput string `json:"expected_output"` // Expected output for the test case
	ActualOutput   string `json:"actual_output"`   // Actual output for the test case
	TimeMs         int    `json:"time_ms"`         // Time taken for the test case in milliseconds
	MemoryKB       int    `json:"memory_kb"`       // Memory used for the test case in kilobytes
}
