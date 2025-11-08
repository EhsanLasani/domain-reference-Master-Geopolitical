// ============================================================================
// FILE: country_validator.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - Validation
// PURPOSE: Complete validation before database operations
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package validation

import (
	"regexp"
	"strings"
)

// CountryValidator implements all database constraints in code
type CountryValidator struct {
	iso2Pattern  *regexp.Regexp
	iso3Pattern  *regexp.Regexp
	phonePattern *regexp.Regexp
}

func NewCountryValidator() *CountryValidator {
	return &CountryValidator{
		iso2Pattern:  regexp.MustCompile(`^[A-Z]{2}$`),
		iso3Pattern:  regexp.MustCompile(`^[A-Z]{3}$`),
		phonePattern: regexp.MustCompile(`^\+[1-9]\d{0,3}$`),
	}
}

// ValidateCountryCode validates ISO 3166-1 alpha-2 format
func (v *CountryValidator) ValidateCountryCode(code string) *ValidationResult {
	result := &ValidationResult{IsValid: true, Errors: make(map[string]string)}
	
	if code == "" {
		result.AddError("country_code", "Country code is required")
	} else if !v.iso2Pattern.MatchString(code) {
		result.AddError("country_code", "Must be 2 uppercase letters (ISO 3166-1)")
	}
	
	return result
}

// ValidateCountryName validates name constraints
func (v *CountryValidator) ValidateCountryName(name string) *ValidationResult {
	result := &ValidationResult{IsValid: true, Errors: make(map[string]string)}
	
	if name == "" {
		result.AddError("country_name", "Country name is required")
	} else if len(name) < 2 || len(name) > 100 {
		result.AddError("country_name", "Must be 2-100 characters")
	} else if strings.Contains(name, "  ") {
		result.AddError("country_name", "Cannot contain double spaces")
	}
	
	return result
}

type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  map[string]string `json:"errors"`
}

func (vr *ValidationResult) AddError(field, message string) {
	vr.IsValid = false
	vr.Errors[field] = message
}