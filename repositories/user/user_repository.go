package user

import (
	"github.com/google/uuid"
	model "neptune/backend/models/user"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (model.User, error)
	FindById(id uuid.UUID) (model.User, error)
}
