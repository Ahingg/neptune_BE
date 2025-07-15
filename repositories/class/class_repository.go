package class

import (
	"context"
	models "neptune/backend/models/class"

	"github.com/google/uuid"
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
	FindClassBySemesterCourseAndStudent(ctx context.Context, semesterID, courseOutlineID, userID string) ([]models.Class, error)
}
