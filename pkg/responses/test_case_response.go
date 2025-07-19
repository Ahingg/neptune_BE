package responses

// TestCaseResponse represents the response structure for a test case. Only for admin
type TestCaseResponse struct {
	CaseID    string `json:"case_id"`
	Number    int    `json:"number"`
	InputUrl  string `json:"input_url"`
	OutputUrl string `json:"output_url"`
}
