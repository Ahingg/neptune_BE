package judgeServ

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"neptune/backend/pkg/requests"
	"net/http"
	"os"
	"time"
)

type judge0ClientImpl struct {
	apiURL     string
	httpClient *http.Client
}

func (c judge0ClientImpl) SubmitCode(sourceCode, stdin string, languageID int) (*Judge0Result, error) {
	if c.apiURL == "" {
		return nil, fmt.Errorf("JUDGE0_API_URL environment variable not set")
	}

	reqBody := requests.Judge0SubmissionRequest{
		SourceCode: sourceCode,
		LanguageID: languageID,
		Stdin:      stdin,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Judge0 request body: %w", err)
	}

	// We add ?wait=true to get the result synchronously
	reqURL := fmt.Sprintf("%s/submissions?wait=true&base64_encoded=false", c.apiURL)

	req, err := http.NewRequestWithContext(context.Background(), "POST", reqURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create Judge0 request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Judge0: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Judge0 returned a non-success status code %d: %s", resp.StatusCode, string(respBody))
	}

	var result Judge0Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Judge0 response: %w", err)
	}

	return &result, nil
}

func NewJudge0Client() Judge0Client {
	return &judge0ClientImpl{
		apiURL: os.Getenv("JUDGE0_API_URL"), // e.g., "http://localhost:2358"
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Set a reasonable timeout
		},
	}
}
