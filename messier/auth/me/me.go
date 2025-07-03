package me

import (
	"context"
)

type MeService interface {
	GetAssistantProfile(ctx context.Context, token string) (*MeResponse, error)
}
