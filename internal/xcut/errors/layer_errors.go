package errors

import (
	"fmt"
	"net/http"
)

type LayerError struct {
	Layer     string `json:"layer"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Cause     error  `json:"-"`
}

func (e *LayerError) Error() string {
	return fmt.Sprintf("[%s:%s] %s", e.Layer, e.Code, e.Message)
}

// Database Layer Errors
func NewDatabaseError(code, message string, cause error) *LayerError {
	return &LayerError{
		Layer:   "database",
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: getErrorDetails(cause),
	}
}

// Data Access Layer Errors
func NewRepositoryError(code, message string, cause error) *LayerError {
	return &LayerError{
		Layer:   "repository",
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: getErrorDetails(cause),
	}
}

// Business Logic Layer Errors
func NewBusinessError(code, message string, cause error) *LayerError {
	return &LayerError{
		Layer:   "business",
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: getErrorDetails(cause),
	}
}

// Presentation Layer Errors
func NewPresentationError(code, message string, cause error) *LayerError {
	return &LayerError{
		Layer:   "presentation",
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: getErrorDetails(cause),
	}
}

// Validation Errors
func NewValidationError(field, message string) *LayerError {
	return &LayerError{
		Layer:   "validation",
		Code:    "VALIDATION_FAILED",
		Message: fmt.Sprintf("Field '%s': %s", field, message),
	}
}

// Schema Errors
func NewSchemaError(message string, cause error) *LayerError {
	return &LayerError{
		Layer:   "schema",
		Code:    "SCHEMA_MISMATCH",
		Message: message,
		Cause:   cause,
		Details: getErrorDetails(cause),
	}
}

func getErrorDetails(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// HTTP Status Code Mapping
func (e *LayerError) HTTPStatus() int {
	switch e.Code {
	case "VALIDATION_FAILED", "SCHEMA_MISMATCH":
		return http.StatusBadRequest
	case "NOT_FOUND":
		return http.StatusNotFound
	case "DUPLICATE_KEY", "CONSTRAINT_VIOLATION":
		return http.StatusConflict
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Error Response for API
type ErrorResponse struct {
	Error LayerError `json:"error"`
}

func NewErrorResponse(err *LayerError) ErrorResponse {
	return ErrorResponse{Error: *err}
}