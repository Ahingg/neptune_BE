package internal_semester

import (
	"context"
	model "neptune/backend/models/semester"
)

type SemesterService interface {
	SyncSemester(ctx context.Context, requestMakerID string) error
	GetInternalSemesters(ctx context.Context) ([]model.Semester, error)
}
