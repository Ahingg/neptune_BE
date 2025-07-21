package container

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	caseHandler "neptune/backend/handlers/case"
	classHand "neptune/backend/handlers/class"
	contestHandler "neptune/backend/handlers/contest"
	"neptune/backend/handlers/language"
	"neptune/backend/handlers/semester"
	submissionHand "neptune/backend/handlers/submission"
	"neptune/backend/handlers/test_case"
	userHand "neptune/backend/handlers/user"
	websocketHand "neptune/backend/handlers/websocket"
	"neptune/backend/messier/auth/log_on"
	"neptune/backend/messier/auth/me"
	externalClass "neptune/backend/messier/class"
	externalSemester "neptune/backend/messier/semester"
	caseRepository "neptune/backend/repositories/case"
	internalClassRepo "neptune/backend/repositories/class"
	contestRepository "neptune/backend/repositories/contest"
	"neptune/backend/repositories/messier_token"
	internalSemesterRepo "neptune/backend/repositories/semester"
	submissionRepo "neptune/backend/repositories/submission"
	testCaseRepo "neptune/backend/repositories/test_case"
	userRepo "neptune/backend/repositories/user"
	caseService "neptune/backend/services/case"
	contestService "neptune/backend/services/contest"
	"neptune/backend/services/internal_class"
	"neptune/backend/services/internal_semester"
	judgeServ "neptune/backend/services/judge0"
	submissionServ "neptune/backend/services/submission"
	testCaseServ "neptune/backend/services/test_case"
	userService "neptune/backend/services/user"
	webSocketService "neptune/backend/services/web_socket_service"
	"os"

	"gorm.io/gorm"
)

type HandlerContainer struct {
	UserHandler             userHand.UserHandler
	InternalSemesterHandler semester.SemesterHandler
	ClassHandler            classHand.ClassHandler
	CaseHandler             caseHandler.CaseHandler
	ContestHandler          contestHandler.ContestHandler
	TestCaseHandler         testCaseHand.TestCaseHandler
	SubmissionHandler       submissionHand.SubmissionHandler
	WebSocketHandler        websocketHand.WebSocketHandler
	LanguageHandler         language.LanguageHandler
}

func NewHandlerContainer(db *gorm.DB) *HandlerContainer {
	// Messier
	logOnService := log_on.NewLogOnService()
	meService := me.NewMeService()
	messierSemesterService := externalSemester.NewExternalSemesterService()
	messierClassService := externalClass.NewMessierClassService()

	// Core
	judge0client := judgeServ.NewJudge0Client()
	webSocketServ := webSocketService.NewWebSocketService()
	webSocketHandler := websocketHand.NewWebSocketHandler(webSocketServ)

	rabbitConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}

	ch, err := rabbitConnection.Channel()
	if err != nil {
		panic("Failed to open a channel: " + err.Error())
	}
	// repo
	messierTokenRepository := messier_token.NewMessierTokenRepository(db)
	semesterRepository := internalSemesterRepo.NewSemesterRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	classRepo := internalClassRepo.NewClassRepository(db)
	caseRepo := caseRepository.NewCaseRepository(db)
	contestRepo := contestRepository.NewContestRepository(db)
	testCaseRepository := testCaseRepo.NewTestCaseRepository(db)
	submissionRepository := submissionRepo.NewSubmissionRepository(db)

	// semester
	semesterService := internal_semester.NewSemesterService(semesterRepository, messierSemesterService, messierTokenRepository)
	internalSemesterHandler := semester.NewSemesterHandler(semesterService)
	// class
	classService := internal_class.NewClassService(messierClassService, classRepo, userRepository, messierTokenRepository)
	classHandler := classHand.NewClassHandler(classService)

	// user
	userServ := userService.NewUserService(userRepository, logOnService, meService, messierTokenRepository, classRepo, semesterRepository)
	userHandler := userHand.NewUserHandler(userServ)

	// case
	caseServ := caseService.NewCaseService(caseRepo)
	caseHand := caseHandler.NewCaseHandler(caseServ)

	// contest
	contestServ := contestService.NewContestService(contestRepo, caseRepo)
	contestHand := contestHandler.NewContestHandler(contestServ)

	// test_case
	testCaseService := testCaseServ.NewTestCaseService(testCaseRepository, caseRepo)
	testCaseHandler := testCaseHand.NewTestCaseHandler(testCaseService, caseServ)

	// submission
	submissionService := submissionServ.NewSubmissionService(submissionRepository, testCaseRepository, ch, judge0client, webSocketServ)
	submissionHandler := submissionHand.NewSubmissionHandler(submissionService)

	go func() {
		if err := submissionService.StartListeners(); err != nil {
			log.Fatalf("Failed to start submission listeners: %v", err)
		}
	}()

	languageHandler := language.NewLanguageHandler()

	return &HandlerContainer{
		UserHandler:             *userHandler,
		InternalSemesterHandler: *internalSemesterHandler,
		ClassHandler:            *classHandler,
		CaseHandler:             *caseHand,
		ContestHandler:          *contestHand,
		TestCaseHandler:         *testCaseHandler,
		WebSocketHandler:        *webSocketHandler,
		SubmissionHandler:       *submissionHandler,
		LanguageHandler:         *languageHandler,
	}
}
