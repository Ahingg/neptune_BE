package internal_class

import (
	"context"
	"neptune/backend/pkg/responses"
)

type ClassService interface {
	SyncClasses(ctx context.Context, semesterID string, requestMakerID string) error
	SyncClassStudents(ctx context.Context, semesterID string, courseOutlineID string, requestMakerID string) error
	SyncClassAssistants(ctx context.Context, semesterID string, courseOutlineID string, requestMakerID string) error

	GetClassesBySemesterAndCourse(ctx context.Context, semesterID string, courseID string) ([]responses.GetClassWithoutDetailResponse, error)
	GetClassDetailBySemesterAndCourse(ctx context.Context, semesterID string, courseID string) ([]responses.GetDetailClassResponse, error)
	GetClassDetailByTransactionID(ctx context.Context, classTransactionID string) (*responses.GetDetailClassResponse, error) // For future CRUD}
}
