package requests

type SyncClassStudentAndAssistantRequest struct {
	SemesterID string `json:"semester_id"`
	CourseID   string `json:"course_id"`
	ClassCode  string `json:"class_code"`
}
