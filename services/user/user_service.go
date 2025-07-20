package user

import (
	"context"
	"github.com/google/uuid"
	model "neptune/backend/models/user"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	"time"
)

type UserService interface {
	LoginAssistant(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error)
	LoginStudent(ctx context.Context, req *requests.LoginRequest) (*responses.LoginResponse, string, time.Time, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)
	DeleteUserAccessToken(ctx context.Context, userID string) error
	GetDetailedUserProfile(ctx context.Context, userID uuid.UUID) (*responses.UserMeResponse, error)
}
