package class

import "github.com/google/uuid"

type GetClassBySemesterAndCourseResponse struct {
	ClassCode          string    `json:"ClassName"`
	ClassTransactionID uuid.UUID `json:"ClassTransactionId"`
	CourseOutlineID    uuid.UUID `json:"CourseOutlineId"`
	SemesterID         uuid.UUID `json:"SemesterId"`
}

// GetClassStudentsResponse untuk dapetin semua mahasiswa berdasarkan Course dan Semester dan ClassCode
type GetClassStudentsResponse struct {
	BinusianID uuid.UUID `json:"BinusianId"`
	NIM        string    `json:"Number"`
	Name       string    `json:"Name"`
}

// GetStudentClassTransactionWithAssistantResponse untuk dapetin AssistantDetail dari NIM mahasiswa
type GetStudentClassTransactionWithAssistantResponse struct {
	Assistants         []string  `json:"Assistants"`
	ClassName          string    `json:"ClassName"`
	ClassTransactionID uuid.UUID `json:"ClassTransactionId"`
	CourseOutlineID    uuid.UUID `json:"CourseOutlineId"`
}

// GetAssistantDetailResponse untuk dapetin Assistant Detail setelah kita dapet Inisial dari StudentClassWithAsssistant
type GetAssistantDetailResponse struct {
	Name     string    `json:"Name"`
	UserID   uuid.UUID `json:"UserId"`
	Username string    `json:"Username"`
}
