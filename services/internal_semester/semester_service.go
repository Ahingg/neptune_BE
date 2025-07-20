package internal_semester

import (
	"context"
	"neptune/backend/pkg/responses"
)

type SemesterService interface {
	SyncSemester(ctx context.Context, requestMakerID string) error
	GetInternalSemesters(ctx context.Context) ([]responses.SemesterResponse, error)
	GetCurrentSemester(ctx context.Context) (*responses.SemesterResponse, error)
}
