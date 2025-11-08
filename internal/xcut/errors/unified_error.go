// Package errors implements unified error model for consistent error handling
package errors

import (
	"fmt"
	"strings"
)

type ErrorCode string

const (
	// Business Logic Errors
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	ErrCodeBusinessRule   ErrorCode = "BUSINESS_RULE_VIOLATION"
	
	// Technical Errors
	ErrCodeDatabase       ErrorCode = "DATABASE_ERROR"
	ErrCodeNetwork        ErrorCode = "NETWORK_ERROR"
	ErrCodeTimeout        ErrorCode = "TIMEOUT_ERROR"
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	
	// Security Errors
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeRateLimit      ErrorCode = "RATE_LIMIT_EXCEEDED"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type StandardError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Retryable  bool                   `json:"retryable"`
	Severity   Severity               `json:"severity"`
	Details    map[string]interface{} `json:"details,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	TenantID   string                 `json:"tenant_id,omitempty"`
	Cause      error                  `json:"-"`
}

func (e *StandardError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *StandardError) Unwrap() error {
	return e.Cause
}

func NewValidationError(message string) *StandardError {
	return &StandardError{
		Code:      ErrCodeValidation,
		Message:   message,
		Retryable: false,
		Severity:  SeverityMedium,
	}
}

func NewNotFoundError(resource string) *StandardError {
	return &StandardError{
		Code:      ErrCodeNotFound,
		Message:   fmt.Sprintf("%s not found", resource),
		Retryable: false,
		Severity:  SeverityLow,
	}
}

func NewDatabaseError(cause error) *StandardError {
	return &StandardError{
		Code:      ErrCodeDatabase,
		Message:   "Database operation failed",
		Retryable: true,
		Severity:  SeverityHigh,
		Cause:     cause,
	}
}

func MapDatabaseError(err error) *StandardError {
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key") {
		return &StandardError{
			Code:      ErrCodeAlreadyExists,
			Message:   "Resource already exists",
			Retryable: false,
			Severity:  SeverityMedium,
			Cause:     err,
		}
	}
	
	return NewDatabaseError(err)
}