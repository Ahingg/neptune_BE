package container

import (
	caseHandler "neptune/backend/handlers/case"
	classHand "neptune/backend/handlers/class"
	contestHandler "neptune/backend/handlers/contest"
	"neptune/backend/handlers/semester"
	"neptune/backend/handlers/test_case"
	userHand "neptune/backend/handlers/user"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	externalClass "neptune/backend/messier/class"
	externalSemester "neptune/backend/messier/semester"
	caseRepository "neptune/backend/repositories/case"
	internalClassRepo "neptune/backend/repositories/class"
	contestRepository "neptune/backend/repositories/contest"
	"neptune/backend/repositories/messier_token"
	internalSemesterRepo "neptune/backend/repositories/semester"
	testCaseRepo "neptune/backend/repositories/test_case"
	userRepo "neptune/backend/repositories/user"
	caseService "neptune/backend/services/case"
	contestService "neptune/backend/services/contest"
	"neptune/backend/services/internal_class"
	"neptune/backend/services/internal_semester"
	testCaseServ "neptune/backend/services/test_case"
	userService "neptune/backend/services/user"

	"gorm.io/gorm"
)

type HandlerContainer struct {
	UserHandler             userHand.UserHandler
	InternalSemesterHandler semester.SemesterHandler
	ClassHandler            classHand.ClassHandler
	CaseHandler             caseHandler.CaseHandler
	ContestHandler          contestHandler.ContestHandler
	TestCaseHandler         testCaseHand.TestCaseHandler
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

	// case
	caseRepo := caseRepository.NewCaseRepository(db)
	caseServ := caseService.NewCaseService(caseRepo)
	caseHand := caseHandler.NewCaseHandler(caseServ)

	// contest
	contestRepo := contestRepository.NewContestRepository(db)
	contestServ := contestService.NewContestService(contestRepo, caseRepo)
	contestHand := contestHandler.NewContestHandler(contestServ)

	// test_case
	testCaseRepository := testCaseRepo.NewTestCaseRepository(db)
	testCaseService := testCaseServ.NewTestCaseService(testCaseRepository, caseRepo)
	testCaseHandler := testCaseHand.NewTestCaseHandler(testCaseService, caseServ)

	return &HandlerContainer{
		UserHandler:             *userHandler,
		InternalSemesterHandler: *internalSemesterHandler,
		ClassHandler:            *classHandler,
		CaseHandler:             *caseHand,
		ContestHandler:          *contestHand,
		TestCaseHandler:         *testCaseHandler,
	}
}
