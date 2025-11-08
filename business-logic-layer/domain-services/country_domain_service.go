// ============================================================================
// FILE: country_domain_service.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Business Logic Layer
// PURPOSE: Pure country business logic - no I/O
// ============================================================================

package domain

import "strings"

// Country domain entity - immutable
type Country struct {
	ID           string
	Code         string
	Name         string
	OfficialName string
	IsActive     bool
	IsDeleted    bool
}

// CountryDomainService contains pure business logic only
type CountryDomainService struct{}

// ValidateCountryCode checks ISO 3166-1 alpha-2 format
// Pure function - no side effects
func (s CountryDomainService) ValidateCountryCode(code string) error {
	normalized := strings.TrimSpace(code)
	if len(normalized) != 2 {
		return fmt.Errorf("country code must be 2 characters")
	}
	if normalized != strings.ToUpper(normalized) {
		return fmt.Errorf("country code must be uppercase")
	}
	return nil
}

// IsEligibleForActivation checks business rules for activation
// Pure function - deterministic business logic
func (s CountryDomainService) IsEligibleForActivation(country Country) bool {
	return !country.IsDeleted && 
		   len(strings.TrimSpace(country.Code)) == 2 &&
		   len(strings.TrimSpace(country.Name)) > 0
}

// CalculateDisplayName determines display name based on business rules
// Pure function - returns new string
func (s CountryDomainService) CalculateDisplayName(country Country) string {
	if country.OfficialName != "" {
		return country.OfficialName
	}
	return country.Name
}

// CompareCountries determines sort order
// Pure function - comparison logic
func (s CountryDomainService) CompareCountries(a, b Country) int {
	return strings.Compare(a.Name, b.Name)
}