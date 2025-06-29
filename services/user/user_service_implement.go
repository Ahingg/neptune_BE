package user

import (
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	"neptune/backend/repositories/user"
)

type userServiceImplement struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userServiceImplement{userRepo: userRepo}
}

func (u *userServiceImplement) Login(req *requests.LoginRequest) (responses.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}
