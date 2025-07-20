package responses

type LoginResponse struct {
	UserID      string                 `json:"id"`
	Username    string                 `json:"username"`
	Name        string                 `json:"name"`
	Role        string                 `json:"role"`
	Enrollments []UserEnrollmentDetail `json:"enrollments,omitempty"`
}

type UserEnrollmentDetail struct {
	ClassTransactionID string `json:"class_transaction_id"`
	ClassName          string `json:"class_name"`
	CourseOutlineID    string `json:"course_outline_id"` // Assuming UUID is string in JSON
	SemesterID         string `json:"semester_id"`       // Assuming UUID is string in JSON
}

type UserMeResponse struct {
	ID          string                 `json:"id"`
	Username    string                 `json:"username"`
	Name        string                 `json:"name"`
	Role        string                 `json:"role"`
	Enrollments []UserEnrollmentDetail `json:"enrollments,omitempty"` // NEW: Optional enrollment details
}
