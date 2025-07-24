package me

import (
	"context"
	"fmt"
	"neptune/backend/messier"
	"neptune/backend/models/user"
	"net/http"
	"os"
	"strings"
)

type meService struct {
	baseURL string
}

func NewMeService() MeService {
	return &meService{
		baseURL: "https://your-lab-api.com", // Replace with actual Lab API URL
	}
}

func (m *meService) GetAssistantProfile(ctx context.Context, authToken string) (*MeResponse, error) {
	var resp MeResponse

	baseURL := strings.TrimRight(os.Getenv("MESSIER_API_URL"), "/")
	err := messier.SendRequest(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/Account/Me", baseURL),
		nil, &resp, authToken)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch me: %w", err)
	}

	// pasti assistant
	var role user.Role = user.RoleAssistant
	resp.Role = role
	return &resp, nil
}
