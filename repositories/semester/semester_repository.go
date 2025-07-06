package semester

import (
	"context"
	model "neptune/backend/models/semester"
)

type SemesterRepository interface {
	Save(ctx context.Context, semester *model.Semester) error
	FindAll(ctx context.Context) ([]model.Semester, error)
	GetSemesterByID(ctx context.Context, semesterID string) (model.Semester, error)
}
