package user

import (
	"context"
	"github.com/google/uuid"
	model "neptune/backend/models/user"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
