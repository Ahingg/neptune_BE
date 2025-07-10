package responses

import models "neptune/backend/models/class"

type GetDetailClassResponse struct {
	ClassTransactionID string                  `json:"class_transaction_id"`
	SemesterID         string                  `json:"semester_id"`
	CourseOutlineID    string                  `json:"course_outline_id"`
	ClassCode          string                  `json:"class_code"`
	Students           []models.ClassStudent   `json:"students"`
	Assistants         []models.ClassAssistant `json:"assistants"`
}
