package errors

import (
	"fmt"
	"time"
)

// UnifiedError implements guideline 03 - Unified Error Model
type UnifiedError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Retryable bool                   `json:"retryable"`
	Severity  string                 `json:"severity"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	TraceID   string                 `json:"trace_id,omitempty"`
	TenantID  string                 `json:"tenant_id,omitempty"`
}

func (e *UnifiedError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewError(code, message string, retryable bool) *UnifiedError {
	return &UnifiedError{
		Code:      code,
		Message:   message,
		Retryable: retryable,
		Severity:  "ERROR",
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]interface{}),
	}
}

func NewSystemError(code, message string) *UnifiedError {
	return NewError(code, message, true)
}

func NewNotFoundError(resource string) *UnifiedError {
	return NewError("GEO-1002", fmt.Sprintf("%s not found", resource), false)
}

func MapDatabaseError(err error) *UnifiedError {
	if err == nil {
		return nil
	}
	return NewError("GEO-1003", fmt.Sprintf("Database error: %v", err), true)
}