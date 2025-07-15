package log_on

import (
	"context"
	"fmt"
	"neptune/backend/messier"
	"net/http"
	"os"
	"strings"
)

type logOnService struct {
	baseURL string
}

func NewLogOnService() LogOnService {
	return &logOnService{
		baseURL: os.Getenv("MESSIER_API_URL"),
	}
}

func (l *logOnService) LogOnAssistant(ctx context.Context, username, password string) (*LogOnResponse, error) {
	messierURL := os.Getenv("MESSIER_API_URL")
	if messierURL == "" {
		return nil, fmt.Errorf("MESSIER_API_URL environment variable is not set")
	}

	req := LogOnRequest{
		Username: username,
		Password: password,
	}
	var resp LogOnResponse
	baseURL := strings.TrimRight(messierURL, "/")
	url := fmt.Sprintf("%s/Account/LogOn", baseURL)

	fmt.Printf("Attempting to log on assistant with URL: %s\n", url)
	err := messier.SendRequest(ctx, http.MethodPost, url, req, &resp, "")
	if err != nil {
		return nil, fmt.Errorf("failed to log on: %w", err)
	}
	return &resp, nil
}

func (l *logOnService) LogOnStudent(ctx context.Context, username, password string) (*LogOnStudentResponse, error) {
	messierURL := os.Getenv("MESSIER_API_URL")
	if messierURL == "" {
		return nil, fmt.Errorf("MESSIER_API_URL environment variable is not set")
	}

	req := LogOnRequest{
		Username: username,
		Password: password,
	}
	var resp LogOnStudentResponse
	baseURL := strings.TrimRight(messierURL, "/")
	url := fmt.Sprintf("%s/Account/LogOnBinusian", baseURL)

	fmt.Printf("Attempting to log on student with URL: %s\n", url)
	err := messier.SendRequest(ctx, http.MethodPost, url, req, &resp, "")
	if err != nil {
		return nil, fmt.Errorf("failed to log on student: %w", err)
	}
	return &resp, nil
}
