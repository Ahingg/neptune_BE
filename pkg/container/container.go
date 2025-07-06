package container

import (
	"gorm.io/gorm"
	"neptune/backend/handlers/semester"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	externalSemester "neptune/backend/messier/semester"
	"neptune/backend/repositories/messier_token"
	internalSemesterRepo "neptune/backend/repositories/semester"
	userRepo "neptune/backend/repositories/user"
	"neptune/backend/services/internal_semester"
	userService "neptune/backend/services/user"
)

type HandlerContainer struct {
	UserHandler             userHand.UserHandler
	InternalSemesterHandler semester.SemesterHandler
}

func NewHandlerContainer(db *gorm.DB) *HandlerContainer {
	// Messier
	logOnService := log_on.NewLogOnService()
	meService := me.NewMeService()
	messierSemesterService := externalSemester.NewExternalSemesterService()

	// user
	userRepository := userRepo.NewUserRepository(db)
	messierTokenRepository := messier_token.NewMessierTokenRepository(db)
	userServ := userService.NewUserService(userRepository, logOnService, meService, messierTokenRepository)
	userHandler := userHand.NewUserHandler(userServ)

	// semester
	semesterRepository := internalSemesterRepo.NewSemesterRepository(db)
	semesterService := internal_semester.NewSemesterService(semesterRepository, messierSemesterService, messierTokenRepository)
	internalSemesterHandler := semester.NewSemesterHandler(semesterService)

	return &HandlerContainer{
		UserHandler:             *userHandler,
		InternalSemesterHandler: *internalSemesterHandler,
	}
}
