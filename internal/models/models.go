// In internal/models/models.go
package models

// QueuedData is the generic wrapper for any data pushed to the queue.
type QueuedData struct {
	RetryCount int         `json:"retry_count"`
	Type       string      `json:"type"` // e.g., "general", "gps"
	Data       interface{} `json:"data"`
}

// APIError is the single, standard error response for the entire API.
type APIError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationErrorDetail is the specific struct for a single validation failure.
// This is what the 'details' field will contain for a validation error.
type ValidationErrorDetail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}