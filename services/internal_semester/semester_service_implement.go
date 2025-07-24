package internal_semester

import (
	"context"
	"fmt"
	"log"
	messierSemester "neptune/backend/messier/semester"
	model "neptune/backend/models/semester"
	"neptune/backend/pkg/responses"
	messierTokenRepo "neptune/backend/repositories/messier_token"
	"neptune/backend/repositories/semester"
	"time"
)

type semesterService struct {
	semesterRepository      semester.SemesterRepository
	externalSemesterService messierSemester.MessierSemesterService
	messierTokenRepository  messierTokenRepo.MessierTokenRepository
}

func (s *semesterService) SyncSemester(ctx context.Context, requestMakerID string) error {
	log.Printf("Starting semester sync for user: %s", requestMakerID)

	adminToken, err := s.messierTokenRepository.GetMessierTokenByUserID(ctx, requestMakerID)
	if err != nil {
		log.Printf("Failed to get admin token for user %s: %v", requestMakerID, err)
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	log.Printf("Found admin token for user %s", requestMakerID)

	if adminToken == nil || adminToken.MessierAccessToken == "" {
		log.Printf("Admin token not found or empty for user ID: %s", requestMakerID)
		return fmt.Errorf("admin token not found or empty for user ID: %s", requestMakerID)
	}

	if adminToken.MessierTokenExpires.Before(time.Now()) {
		log.Printf("Admin token has expired for user ID: %s", requestMakerID)
		return fmt.Errorf("admin token has expired for user ID: %s", requestMakerID)
	}

	log.Printf("Admin token is valid, fetching semesters from external API")

	messierAccessToken := adminToken.MessierAccessToken

	externalSemesters, err := s.externalSemesterService.GetSemesters(ctx, messierAccessToken)
	if err != nil {
		log.Printf("Failed to get external semesters: %v", err)
		return fmt.Errorf("failed to get external semesters: %w", err)
	}

	log.Printf("Successfully fetched %d semesters from external API", len(externalSemesters))

	for _, ms := range externalSemesters {
		if ms.Start.IsZero() { // Check if the parsed Start time is its zero value
			return fmt.Errorf("received an invalid (zero) Start time for semester ID %s, which is unexpected", ms.SemesterID)
		}
		var endTimePtr *time.Time
		if !ms.End.IsZero() { // If the parsed End time is NOT its zero value
			// Create a new time.Time variable and take its address
			actualEndTime := ms.End.Time
			endTimePtr = &actualEndTime // Assign the address to the pointer
		}
		sem := &model.Semester{
			ID:          ms.SemesterID,
			Description: ms.Description,
			Start:       ms.Start.Time,
			End:         endTimePtr,
		}
		log.Printf("Saving semester: ID=%s, Description=%s", sem.ID, sem.Description)
		if err := s.semesterRepository.Save(ctx, sem); err != nil {
			log.Printf("Warning: Failed to save or update semester %s: %v", ms.SemesterID, err)
			// You might choose to return an error here or continue processing
		} else {
			log.Printf("Successfully saved semester: ID=%s, Description=%s", sem.ID, sem.Description)
		}
	}

	log.Printf("Successfully synced %d semesters from Messier API.", len(externalSemesters))
	return nil
}

func (s *semesterService) GetInternalSemesters(ctx context.Context) ([]responses.SemesterResponse, error) {
	log.Printf("Getting internal semesters from database")
	semesters, err := s.semesterRepository.FindAll(ctx)
	if err != nil {
		log.Printf("Failed to retrieve semesters from internal DB: %v", err)
		return nil, fmt.Errorf("failed to retrieve semesters from internal DB: %w", err)
	}

	log.Printf("Successfully retrieved %d semesters from internal DB", len(semesters))

	semResponses := make([]responses.SemesterResponse, len(semesters))
	for i, sem := range semesters {
		semResponses[i] = responses.SemesterResponse{
			SemesterID:  sem.ID,
			Description: sem.Description,
			Start:       sem.Start,
			End:         sem.End, // This will be nil if End is not set in the DB
		}
		log.Printf("Semester %d: ID=%s, Description=%s", i+1, sem.ID, sem.Description)
	}
	return semResponses, nil
}

func (s *semesterService) GetCurrentSemester(ctx context.Context) (*responses.SemesterResponse, error) {
	currSemester, err := s.semesterRepository.FindCurrentSemester(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current semester: %w", err)
	}

	semResponse := &responses.SemesterResponse{
		SemesterID:  currSemester.ID,
		Description: currSemester.Description,
		Start:       currSemester.Start,
		End:         currSemester.End,
	}
	return semResponse, nil
}

func NewSemesterService(semesterRepo semester.SemesterRepository,
	externalSemesterRepo messierSemester.MessierSemesterService,
	messierTokenRepository messierTokenRepo.MessierTokenRepository) SemesterService {
	return &semesterService{
		semesterRepository:      semesterRepo,
		externalSemesterService: externalSemesterRepo,
		messierTokenRepository:  messierTokenRepository,
	}
}
