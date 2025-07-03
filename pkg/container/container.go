package container

import (
	"gorm.io/gorm"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	userRepo "neptune/backend/repositories/user"
	userService "neptune/backend/services/user"
)

type HandlerContainer struct {
	UserHandler userHand.UserHandler
}

func NewHandlerContainer(db *gorm.DB) *HandlerContainer {
	// Messier
	logOnService := log_on.NewLogOnService()
	meService := me.NewMeService()

	// user
	userRepository := userRepo.NewUserRepository(db)
	userServ := userService.NewUserService(userRepository, logOnService, meService)
	userHandler := userHand.NewUserHandler(userServ)

	return &HandlerContainer{
		UserHandler: *userHandler,
	}
}
