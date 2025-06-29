package user

import (
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
)

type UserService interface {
	Login(req *requests.LoginRequest) (responses.LoginResponse, error)
}
