package container

import (
	"gorm.io/gorm"
	classHand "neptune/backend/handlers/class"
	"neptune/backend/handlers/semester"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	externalClass "neptune/backend/messier/class"
	externalSemester "neptune/backend/messier/semester"
	internalClassRepo "neptune/backend/repositories/class"
	"neptune/backend/repositories/messier_token"
	internalSemesterRepo "neptune/backend/repositories/semester"
	userRepo "neptune/backend/repositories/user"
	"neptune/backend/services/internal_class"
	"neptune/backend/services/internal_semester"
	userService "neptune/backend/services/user"
)

type HandlerContainer struct {
	UserHandler             userHand.UserHandler
	InternalSemesterHandler semester.SemesterHandler
	ClassHandler            classHand.ClassHandler
}

func NewHandlerContainer(db *gorm.DB) *HandlerContainer {
	// Messier
	logOnService := log_on.NewLogOnService()
	meService := me.NewMeService()
	messierSemesterService := externalSemester.NewExternalSemesterService()
	messierClassService := externalClass.NewMessierClassService()

	// user
	userRepository := userRepo.NewUserRepository(db)
	messierTokenRepository := messier_token.NewMessierTokenRepository(db)
	userServ := userService.NewUserService(userRepository, logOnService, meService, messierTokenRepository)
	userHandler := userHand.NewUserHandler(userServ)

	// semester
	semesterRepository := internalSemesterRepo.NewSemesterRepository(db)
	semesterService := internal_semester.NewSemesterService(semesterRepository, messierSemesterService, messierTokenRepository)
	internalSemesterHandler := semester.NewSemesterHandler(semesterService)

	// class
	classRepo := internalClassRepo.NewClassRepository(db)
	classService := internal_class.NewClassService(messierClassService, classRepo, userRepository, messierTokenRepository)
	classHandler := classHand.NewClassHandler(classService)

	return &HandlerContainer{
		UserHandler:             *userHandler,
		InternalSemesterHandler: *internalSemesterHandler,
		ClassHandler:            *classHandler,
	}
}
