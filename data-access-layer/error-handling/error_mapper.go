// ============================================================================
// FILE: error_mapper.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - Error Handling
// PURPOSE: Map database errors to domain-specific error codes
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package errorhandling

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Domain error codes for geopolitical domain
const (
	// GEO-1xxx: Country errors
	ErrCountryNotFound     = "GEO-1001"
	ErrCountryDuplicate    = "GEO-1002" 
	ErrCountryInvalid      = "GEO-1003"
	ErrCountryVersionConflict = "GEO-1004"
	
	// GEO-2xxx: Region errors
	ErrRegionNotFound      = "GEO-2001"
	ErrRegionDuplicate     = "GEO-2002"
	
	// GEO-3xxx: Language errors
	ErrLanguageNotFound    = "GEO-3001"
	ErrLanguageDuplicate   = "GEO-3002"
	
	// GEO-9xxx: System errors
	ErrDatabaseConnection  = "GEO-9001"
	ErrDatabaseTimeout     = "GEO-9002"
)

// DomainError represents a business domain error
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Cause   error  `json:"-"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorMapper maps database errors to domain errors
type ErrorMapper struct{}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapError converts database errors to domain-specific errors
func (m *ErrorMapper) MapError(err error, entity string) error {
	if err == nil {
		return nil
	}

	// Handle standard SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		return m.mapNotFoundError(entity)
	}

	// Handle PostgreSQL constraint violations
	errMsg := strings.ToLower(err.Error())
	
	// Unique constraint violations
	if strings.Contains(errMsg, "duplicate key") {
		return m.mapDuplicateError(entity, errMsg)
	}
	
	// Check constraint violations
	if strings.Contains(errMsg, "check constraint") {
		return m.mapValidationError(entity, errMsg)
	}
	
	// Foreign key violations
	if strings.Contains(errMsg, "foreign key") {
		return m.mapForeignKeyError(entity, errMsg)
	}
	
	// Connection errors
	if strings.Contains(errMsg, "connection") {
		return &DomainError{
			Code:    ErrDatabaseConnection,
			Message: "Database connection failed",
			Cause:   err,
		}
	}
	
	// Default to system error
	return &DomainError{
		Code:    "GEO-9999",
		Message: "Unknown database error",
		Cause:   err,
	}
}

func (m *ErrorMapper) mapNotFoundError(entity string) error {
	switch strings.ToLower(entity) {
	case "country":
		return &DomainError{
			Code:    ErrCountryNotFound,
			Message: "Country not found",
		}
	case "region":
		return &DomainError{
			Code:    ErrRegionNotFound,
			Message: "Region not found",
		}
	case "language":
		return &DomainError{
			Code:    ErrLanguageNotFound,
			Message: "Language not found",
		}
	default:
		return &DomainError{
			Code:    "GEO-0001",
			Message: fmt.Sprintf("%s not found", entity),
		}
	}
}

func (m *ErrorMapper) mapDuplicateError(entity, errMsg string) error {
	switch strings.ToLower(entity) {
	case "country":
		if strings.Contains(errMsg, "country_code") {
			return &DomainError{
				Code:    ErrCountryDuplicate,
				Message: "Country code already exists",
				Field:   "country_code",
			}
		}
	case "region":
		return &DomainError{
			Code:    ErrRegionDuplicate,
			Message: "Region code already exists",
			Field:   "region_code",
		}
	}
	
	return &DomainError{
		Code:    "GEO-0002",
		Message: "Duplicate entry",
	}
}

func (m *ErrorMapper) mapValidationError(entity, errMsg string) error {
	return &DomainError{
		Code:    ErrCountryInvalid,
		Message: "Validation constraint violated",
	}
}

func (m *ErrorMapper) mapForeignKeyError(entity, errMsg string) error {
	return &DomainError{
		Code:    "GEO-0003",
		Message: "Referenced entity does not exist",
	}
}