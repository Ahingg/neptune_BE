package class

type GetClassBySemesterAndCourseResponse []struct {
	ClassName          string `json:"ClassName"`
	ClassTransactionID string `json:"ClassTransactionId"`
	CourseOutlineID    string `json:"CourseOutlineId"`
}
