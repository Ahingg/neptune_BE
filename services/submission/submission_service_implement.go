package submissionServ

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	submissionModel "neptune/backend/models/submission"
	"neptune/backend/pkg/amqp_messages"
	"neptune/backend/pkg/requests"
	"neptune/backend/pkg/responses"
	submissionRepo "neptune/backend/repositories/submission"
	testCaseRepo "neptune/backend/repositories/test_case"
	userRepo "neptune/backend/repositories/user"
	contestService "neptune/backend/services/contest"
	judgeServ "neptune/backend/services/judge0"
	webSocketService "neptune/backend/services/web_socket_service"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type submissionService struct {
	submissionRepository submissionRepo.SubmissionRepository
	testCaseRepository   testCaseRepo.TestCaseRepository
	contestService       contestService.ContestService
	rabbitChannel        *amqp.Channel
	judgeClient          judgeServ.Judge0Client
	webSocketManager     webSocketService.WebSocketService
	userRepository       userRepo.UserRepository
}

func (s *submissionService) SubmitCode(ctx context.Context, req *requests.SubmitCodeRequest, userID uuid.UUID) (*responses.SubmitCodeResponse, error) {
	submission := &submissionModel.Submission{
		ID:                 uuid.New(),
		CaseID:             req.CaseID,
		UserID:             userID,
		LanguageID:         req.LanguageID,
		ClassTransactionID: req.ClassTransactionID,
		ContestID:          &req.ContestID,
		Status:             submissionModel.SubmissionStatusJudging, // Start as In Queue
		Score:              0,
	}

	submissionDir := filepath.Join("public/submissions", submission.ID.String())
	if err := os.MkdirAll(submissionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create submission directory: %w", err)
	}

	// Use the validated extension from the request struct
	fileName := "main" + req.FileExtension
	sourcePath := filepath.Join(submissionDir, fileName)
	submission.SourceCodePath = "/" + sourcePath // Store URL path

	// Use the byte slice from the request struct, which could have come from the string or file
	if err := os.WriteFile(sourcePath, req.SourceCodeBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write source code: %w", err)
	}

	// Save initial submission record to DB
	if err := s.submissionRepository.Save(ctx, submission); err != nil {
		return nil, fmt.Errorf("failed to save submission record: %w", err)
	}

	// --- Publish to RabbitMQ (logic remains the same) ---
	msgBody, err := json.Marshal(amqp_messages.JudgeQueueMessage{SubmissionID: submission.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal judge queue message: %w", err)
	}

	err = s.rabbitChannel.Publish(
		"",                           // exchange
		amqp_messages.JudgeQueueName, // routing key (queue name)
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish to judge queue: %w", err)
	}

	log.Printf("Successfully queued submission %s for judging", submission.ID)

	resp := &responses.SubmitCodeResponse{
		SubmissionID: submission.ID.String(),
		Status:       submission.Status.String(),
	}
	return resp, nil
}

func (s *submissionService) GetSubmissionByUserInContest(ctx context.Context, userID uuid.UUID, contestID uuid.UUID, classTransactionID *uuid.UUID) ([]responses.GetSubmissionByUserInContestResponse, error) {
	submissions, err := s.submissionRepository.FindByUserInContest(ctx, contestID, userID, classTransactionID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions for user %s in contest %s: %w", userID, contestID, err)
	}

	var resp []responses.GetSubmissionByUserInContestResponse
	for _, sub := range submissions {
		submissionCase, err := s.contestService.GetContentCaseByCaseID(ctx, contestID, sub.CaseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get contest case on submission ln 107 :%w", err)
		}

		resp = append(resp, responses.GetSubmissionByUserInContestResponse{
			SubmissionID: sub.ID.String(),
			ContestID:    sub.ContestID.String(),
			CaseID:       sub.CaseID.String(),
			CaseCode:     submissionCase.CaseCode,
			Status:       sub.Status.String(),
			Score:        sub.Score,
			SubmitTime:   sub.CreatedAt.Format(time.RFC3339),
			LanguageID:   sub.LanguageID,
		})
	}
	return resp, nil
}

func (s *submissionService) GetClassContestSubmissions(ctx context.Context, classTransactionID uuid.UUID, contestID uuid.UUID) ([]responses.GetSubmissionByClassContestResponse, error) {
	submissions, err := s.submissionRepository.FindClassSubmissions(ctx, classTransactionID, contestID)
	if err != nil {

		return nil, fmt.Errorf("failed to fetch class submissions for class %s in contest %s: %w", classTransactionID, contestID, err)
	}
	var resp []responses.GetSubmissionByClassContestResponse
	for _, sub := range submissions {
		submissionCase, err := s.contestService.GetContentCaseByCaseID(ctx, contestID, sub.CaseID)

		if err != nil {
			return nil, fmt.Errorf("failed to get contest case on submission ln 127 :%w", err)
		}

		user, err := s.userRepository.GetUserByID(ctx, sub.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by ID %s for submission %s: %w", sub.UserID, sub.ID, err)
		}

		resp = append(resp, responses.GetSubmissionByClassContestResponse{
			UserID:       sub.UserID.String(),
			Username:     user.Username,
			Name:         user.Name,
			SubmissionID: sub.ID.String(),
			ContestID:    sub.ContestID.String(),
			CaseID:       sub.CaseID.String(),
			CaseCode:     submissionCase.CaseCode,
			Status:       sub.Status.String(),
			Score:        sub.Score,
			SubmitTime:   sub.CreatedAt.Format(time.RFC3339),
			LanguageID:   sub.LanguageID,
		})
	}

	return resp, nil
}

func (s *submissionService) StartListeners() error {
	_, err := s.rabbitChannel.QueueDeclare(amqp_messages.JudgeQueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare judge queue: %w", err)
	}
	_, err = s.rabbitChannel.QueueDeclare(amqp_messages.ResultQueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare result queue: %w", err)
	}

	// Listener for judge_queue
	judgeMsgs, err := s.rabbitChannel.Consume(amqp_messages.JudgeQueueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume from judge queue: %w", err)
	}

	// Listener for result_queue
	resultMsgs, err := s.rabbitChannel.Consume(amqp_messages.ResultQueueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume from result queue: %w", err)
	}

	// Start processing in background goroutines
	go func() {
		for d := range judgeMsgs {
			log.Printf("Received judging job: %s", d.Body)
			s.processSubmissionJob(context.Background(), d)
		}
	}()

	go func() {
		for d := range resultMsgs {
			log.Printf("Received result job: %s", d.Body)
			s.processResultJob(context.Background(), d)
		}
	}()

	log.Println("Submission listeners started successfully")
	return nil
}

func (s *submissionService) processSubmissionJob(ctx context.Context, d amqp.Delivery) {
	defer d.Ack(false) // Acknowledge message when done

	var msg amqp_messages.JudgeQueueMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("Error unmarshalling judge job: %v", err)
		return
	}

	submission, err := s.submissionRepository.FindByID(ctx, msg.SubmissionID.String())
	if err != nil {
		log.Printf("Error finding submission %s: %v", msg.SubmissionID, err)
		return
	}

	// Helper function to publish an error status and exit
	publishError := func(status submissionModel.SubmissionStatus) {
		resultMsg := amqp_messages.ResultQueueMessage{
			SubmissionID: submission.ID,
			FinalStatus:  status,
			Results:      []submissionModel.SubmissionResult{},
			Score:        0,
		}
		resultBody, _ := json.Marshal(resultMsg)
		s.rabbitChannel.Publish("", amqp_messages.ResultQueueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        resultBody,
		})
	}

	// Update status to Judging and notify client
	submission.Status = submissionModel.SubmissionStatusJudging
	err = s.submissionRepository.Update(ctx, submission)
	if err != nil {
		log.Printf("Error updating submission %s to Judging: %v", submission.ID, err)
		publishError(submissionModel.SubmissionStatusInternalError)
		return
	}

	resp := responses.FinalResultResponse{
		SubmissionID: submission.ID.String(),
		Status:       submission.Status.String(),
		CaseID:       submission.CaseID.String(),
		TestCases:    []responses.TestCaseJudgeResponse{},
	}

	s.webSocketManager.SendUpdateToClient(submission.ID, resp)

	testcases, err := s.testCaseRepository.FindTestCaseByCaseID(ctx, submission.CaseID.String())
	if err != nil {
		log.Printf("Error fetching testcases for case %s: %v", submission.CaseID, err)
		publishError(submissionModel.SubmissionStatusInternalError)
		return
	}

	sourceCodeBytes, err := os.ReadFile(submission.SourceCodePath[1:]) // remove leading '/'
	if err != nil {
		log.Printf("Error reading source code for submission %s: %v", submission.ID, err)
		publishError(submissionModel.SubmissionStatusInternalError)
		return
	}

	var results []submissionModel.SubmissionResult
	overallStatus := submissionModel.SubmissionStatusAccepted

	// ---- Main Judging Loop ----
	for _, tc := range testcases {
		inputBytes, err := os.ReadFile(tc.InputUrl[1:])
		if err != nil {
			log.Printf("Error reading input file %s: %v", tc.InputUrl, err)
			overallStatus = submissionModel.SubmissionStatusInternalError
			break
		}

		expectedOutputBytes, err := os.ReadFile(tc.OutputUrl[1:])
		if err != nil {
			log.Printf("Error reading output file %s: %v", tc.OutputUrl, err)
			overallStatus = submissionModel.SubmissionStatusInternalError
			break
		}
		expectedOutput := string(expectedOutputBytes)

		// Call Judge0
		judgeResult, err := s.judgeClient.SubmitCode(string(sourceCodeBytes), string(inputBytes), submission.LanguageID)
		if err != nil {
			log.Printf("Error submitting to Judge0: %v", err)
			overallStatus = submissionModel.SubmissionStatusInternalError
			break
		}

		// 1. Convert Judge0 status to our internal status
		currentTestcaseStatus := mapJudge0Status(judgeResult.Status.ID, judgeResult.Stdout, expectedOutput)

		// 2. Build the detailed result object
		time, _ := strconv.ParseFloat(judgeResult.Time, 64)
		newResult := submissionModel.SubmissionResult{
			SubmissionID:   submission.ID,
			TestcaseNumber: tc.Number,
			Status:         currentTestcaseStatus,
			TimeSeconds:    time,
			MemoryKB:       judgeResult.Memory,
			Input:          string(inputBytes),
			ExpectedOutput: expectedOutput,
			ActualOutput:   judgeResult.Stdout,
		}

		// Append compile or runtime error details to the output for user feedback
		if judgeResult.Stderr != "" {
			newResult.ActualOutput += "\n--- STDERR ---\n" + judgeResult.Stderr
		}
		if judgeResult.CompileOutput != "" {
			newResult.ActualOutput += "\n--- COMPILE OUTPUT ---\n" + judgeResult.CompileOutput
		}

		results = append(results, newResult)

		// 3. Update overall status only on the first failure
		if overallStatus == submissionModel.SubmissionStatusAccepted && currentTestcaseStatus != submissionModel.SubmissionStatusAccepted {
			overallStatus = currentTestcaseStatus
		}

		// 4. If a test case failed, stop judging the rest
		if overallStatus != submissionModel.SubmissionStatusAccepted {
			break
		}
	}

	// ---- Post-Judging ----
	resultMsg := amqp_messages.ResultQueueMessage{
		SubmissionID: submission.ID,
		FinalStatus:  overallStatus,
		Results:      results,
		Score:        0,
	}

	if overallStatus == submissionModel.SubmissionStatusAccepted {
		resultMsg.Score = 100
	}

	resultBody, _ := json.Marshal(resultMsg)
	s.rabbitChannel.Publish("", amqp_messages.ResultQueueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        resultBody,
	})
}

func mapJudge0Status(judgeStatusID int, stdout, expectedOutput string) submissionModel.SubmissionStatus {
	switch judgeStatusID {
	case 3: // Accepted
		// CRITICAL: Judge0's "Accepted" only means the code ran. We must verify the output.
		if stdout == expectedOutput {
			return submissionModel.SubmissionStatusAccepted
		}
		return submissionModel.SubmissionStatusWrongAnswer
	case 4:
		return submissionModel.SubmissionStatusWrongAnswer
	case 5:
		return submissionModel.SubmissionStatusTimeLimitExceeded
	case 6:
		return submissionModel.SubmissionStatusCompileError
	case 7, 8, 9, 10, 11, 12:
		return submissionModel.SubmissionStatusRuntimeError
	default:
		return submissionModel.SubmissionStatusInternalError
	}
}

func (s *submissionService) processResultJob(ctx context.Context, d amqp.Delivery) {
	defer d.Ack(false)

	var msg amqp_messages.ResultQueueMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("Error unmarshalling result job: %v", err)
		return
	}

	// Find the original submission
	submission, err := s.submissionRepository.FindByID(ctx, msg.SubmissionID.String())
	if err != nil {
		log.Printf("Error finding submission %s for final update: %v", msg.SubmissionID, err)
		return
	}

	// Update the final status and score
	submission.Status = msg.FinalStatus
	submission.Score = msg.Score
	submission.UpdatedAt = time.Now()

	// Create Response to backend

	// Use a transaction to update submission and save results
	// tx := s.db.Begin() ... (For simplicity, not showing full transaction code)
	finalScore := getFinalScore(msg.Results)
	submission.Score = finalScore
	if err := s.submissionRepository.Update(ctx, submission); err != nil {
		log.Printf("Error performing final update on submission %s: %v", submission.ID, err)
		return
	}

	// Save the detailed per-testcase results
	if err := s.submissionRepository.SaveResultsBatch(ctx, msg.Results); err != nil {
		log.Printf("Error saving batch results for submission %s: %v", submission.ID, err)
		return
	}

	// Push final result to client via WebSocket
	testCases := make([]responses.TestCaseJudgeResponse, len(msg.Results))
	for i, result := range msg.Results {
		testCases[i] = responses.TestCaseJudgeResponse{
			Number:         result.TestcaseNumber,
			Verdict:        result.Status.String(),
			Input:          result.Input,
			ExpectedOutput: result.ExpectedOutput,
			ActualOutput:   result.ActualOutput,
			TimeMs:         int(result.TimeSeconds * 1000), // Convert seconds to milliseconds
			MemoryKB:       result.MemoryKB,
		}
	}

	finalResultResponse := &responses.FinalResultResponse{
		SubmissionID: submission.ID.String(),
		Status:       submission.Status.String(),
		CaseID:       submission.CaseID.String(),
		Score:        finalScore,
		TestCases:    testCases,
	}
	log.Printf("Pushing final update to WebSocket for submission %s", submission.ID)
	s.webSocketManager.SendUpdateToClient(submission.ID, finalResultResponse)
}

func getFinalScore(results []submissionModel.SubmissionResult) int {
	if len(results) == 0 {
		return 0
	}

	resultCount := len(results)
	testCaseCorrectCount := 0
	for _, result := range results {
		if result.Status == submissionModel.SubmissionStatusAccepted {
			testCaseCorrectCount++
		}
	}

	return int(float64(testCaseCorrectCount) / float64(resultCount) * 100)
}

func NewSubmissionService(repo submissionRepo.SubmissionRepository,
	testCaseRepo testCaseRepo.TestCaseRepository,
	rabbitChannel *amqp.Channel,
	judgeClient judgeServ.Judge0Client,
	webSocketManager webSocketService.WebSocketService,
	contestServ contestService.ContestService,
	userRepo userRepo.UserRepository) SubmissionService {
	return &submissionService{
		submissionRepository: repo,
		testCaseRepository:   testCaseRepo,
		rabbitChannel:        rabbitChannel,
		judgeClient:          judgeClient,
		webSocketManager:     webSocketManager,
		contestService:       contestServ,
		userRepository:       userRepo, // Assuming you have a user repository
	}
}
