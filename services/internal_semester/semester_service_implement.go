package internal_semester

import (
	"context"
	"fmt"
	"log"
	messierSemester "neptune/backend/messier/semester"
	model "neptune/backend/models/semester"
	messierTokenRepo "neptune/backend/repositories/messier_token"
	"neptune/backend/repositories/semester"
	"time"
)

type semesterService struct {
	semesterRepository      semester.SemesterRepository
	externalSemesterService messierSemester.MessierSemesterService
	messierTokenRepository  messierTokenRepo.MessierTokenRepository
}

func (s semesterService) SyncSemester(ctx context.Context, requestMakerID string) error {
	adminToken, err := s.messierTokenRepository.GetMessierTokenByUserID(ctx, requestMakerID)
	if err != nil {
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	if adminToken == nil || adminToken.MessierAccessToken == "" {
		return fmt.Errorf("admin token not found or empty for user ID: %s", requestMakerID)
	}

	if adminToken.MessierTokenExpires.Before(time.Now()) {
		return fmt.Errorf("admin token has expired for user ID: %s", requestMakerID)
	}

	messierAccessToken := adminToken.MessierAccessToken

	externalSemesters, err := s.externalSemesterService.GetSemesters(ctx, messierAccessToken)
	if err != nil {
		return fmt.Errorf("failed to get external semesters: %w", err)
	}

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
		if err := s.semesterRepository.Save(ctx, sem); err != nil {
			log.Printf("Warning: Failed to save or update semester %s: %v", ms.SemesterID, err)
			// You might choose to return an error here or continue processing
		}
	}

	log.Printf("Successfully synced %d semesters from Messier API.", len(externalSemesters))
	return nil
}

func (s semesterService) GetInternalSemesters(ctx context.Context) ([]model.Semester, error) {
	semesters, err := s.semesterRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve semesters from internal DB: %w", err)
	}
	return semesters, nil
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
