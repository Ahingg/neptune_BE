package class

import "context"

type MessierClassService interface {
	GetClassesBySemesterAndCourseOutline(ctx context.Context, semesterId string, courseId string,
		authToken string) ([]GetClassBySemesterAndCourseResponse, error)
	GetAssistantInitialFromStudentTransaction(ctx context.Context, studentNIM string,
		semesterId string, authToken string) ([]GetStudentClassTransactionWithAssistantResponse, error)
	GetAssistantDetailFromAssistantInitial(ctx context.Context, initial string,
		generation string, authToken string) (*GetAssistantDetailResponse, error)
	GetStudentFromClassTransaction(ctx context.Context, semesterId string, courseOutlineId string, className string,
		authToken string) ([]GetClassStudentsResponse, error)
}
