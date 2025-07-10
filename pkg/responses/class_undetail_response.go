package responses

type GetClassWithoutDetailResponse struct {
	ClassCode          string `json:"class_code"`
	ClassTransactionID string `json:"class_transaction_id"`
	SemesterID         string `json:"semester_id"`
	CourseOutlineID    string `json:"course_outline_id"`
}
