package log_on

import (
	"context"
	"fmt"
	"neptune/backend/messier"
	"net/http"
	"os"
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
	req := LogOnRequest{
		Username: username,
		Password: password,
	}
	var resp LogOnResponse
	url := fmt.Sprintf("%s/Account/LogOn", os.Getenv("MESSIER_API_URL"))

	err := messier.SendRequest(ctx, http.MethodPost, url, req, &resp, "")
	if err != nil {
		return nil, fmt.Errorf("failed to log on: %w", err)
	}
	return &resp, nil
}

func (l *logOnService) LogOnStudent(ctx context.Context, username, password string) (*LogOnStudentResponse, error) {
	req := LogOnRequest{
		Username: username,
		Password: password,
	}
	var resp LogOnStudentResponse
	url := fmt.Sprintf("%s/Account/LogOnBinusian", os.Getenv("MESSIER_API_URL"))

	err := messier.SendRequest(ctx, http.MethodPost, url, req, &resp, "")
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
