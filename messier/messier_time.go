package messier

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessierTime This package provides a custom time type for handling Messier API date-time strings.
type MessierTime struct {
	time.Time
}

const messierTimeLayout = "2006-01-02T15:04:05.999"

// UnmarshalJSON implements the json.Unmarshaler interface for MessierTime.
// It parses the incoming string using the custom layout.
func (mt *MessierTime) UnmarshalJSON(b []byte) error {
	// First, check for the JSON 'null' literal
	if string(b) == "null" {
		mt.Time = time.Time{} // Assign zero value for time.Time if it's null
		return nil
	}

	// Try to unmarshal the bytes into a string
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		// If it's not a string or "null", return an error
		return fmt.Errorf("MessierTime: expecting a string or null, got %s: %w", string(b), err)
	}

	// If the string content itself is "null" (e.g., {"End": "null"}), handle it
	if s == "null" {
		mt.Time = time.Time{} // Assign zero value for time.Time
		return nil
	}

	// Now, proceed to parse the actual time string
	parsedTime, err := time.Parse(messierTimeLayout, s)
	if err != nil {
		// Provide more context in the error message for debugging
		return fmt.Errorf("failed to parse MessierTime '%s' with layout '%s': %w", s, messierTimeLayout, err)
	}

	mt.Time = parsedTime
	return nil
}
