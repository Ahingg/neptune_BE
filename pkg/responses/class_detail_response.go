package responses

type GetDetailClassResponse struct {
	ClassTransactionID string              `json:"class_transaction_id"`
	SemesterID         string              `json:"semester_id"`
	CourseOutlineID    string              `json:"course_outline_id"`
	ClassCode          string              `json:"class_code"`
	Students           []ClassUserResponse `json:"students"`
	Assistants         []ClassUserResponse `json:"assistants"`
}

type ClassUserResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
}
