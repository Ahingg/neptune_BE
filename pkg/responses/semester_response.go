package responses

import "time"

type SemesterResponse struct {
	SemesterID  string     `json:"SemesterID"`
	Description string     `json:"Description"`
	Start       time.Time  `json:"Start"`         // Assuming Start is a string representation of time
	End         *time.Time `json:"End,omitempty"` // Pointer to string to allow null values
}
