package responses

import "time"

type SemesterResponse struct {
	SemesterID  string     `json:"semester_id"`
	Description string     `json:"description"`
	Start       time.Time  `json:"start"`         // Assuming Start is a string representation of time
	End         *time.Time `json:"end,omitempty"` // Pointer to string to allow null values
}
