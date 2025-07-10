package class

import (
	"context"
	"github.com/google/uuid"
	models "neptune/backend/models/class"
)

type ClassRepository interface {
	SaveClass(ctx context.Context, class *models.Class) error
	AddClassStudents(ctx context.Context, classTransactionID string, studentUserIDs []uuid.UUID) error
	AddClassAssistants(ctx context.Context, classTransactionID string, assistantUserIDs []uuid.UUID) error
	ClearClassStudents(ctx context.Context, classTransactionID string) error
	ClearClassAssistants(ctx context.Context, classTransactionID string) error

	FindAllClassesBySemesterAndCourse(ctx context.Context, semesterID, courseOutlineID string) ([]models.Class, error)
	FindFirstStudentByClassTransactionID(ctx context.Context, classTransactionID string) (*models.ClassStudent, error)
	FindClassByTransactionID(ctx context.Context, classTransactionID string) (*models.Class, error)
	FindClassBasicInfoBySemesterAndCourse(ctx context.Context, semesterID, courseOutlineID string) ([]models.Class, error)
}
