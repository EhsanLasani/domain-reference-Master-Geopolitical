// ============================================================================
// FILE: country_mapper.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer
// PURPOSE: Country DTO mapping with pure functions
// ============================================================================

package mappers

import "strings"

// CountryDTO for external APIs
type CountryDTO struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	OfficialName string `json:"official_name,omitempty"`
	ISO3Code     string `json:"iso3_code,omitempty"`
}

// Country domain entity
type Country struct {
	ID           string
	Code         string
	Name         string
	OfficialName string
	ISO3Code     string
	IsActive     bool
}

// NormalizeCountryCode normalizes country code format
// Pure function - same input = same output
func NormalizeCountryCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

// NormalizeCountryName normalizes country name format
// Pure function - same input = same output
func NormalizeCountryName(name string) string {
	return strings.TrimSpace(name)
}

// CountryToDTO converts domain entity to DTO
// Pure function - no side effects
func CountryToDTO(country Country) CountryDTO {
	return CountryDTO{
		Code:         country.Code,
		Name:         country.Name,
		OfficialName: country.OfficialName,
		ISO3Code:     country.ISO3Code,
	}
}

// DTOToCountry converts DTO to domain entity
// Pure function - returns new immutable object
func DTOToCountry(dto CountryDTO) Country {
	return Country{
		Code:         NormalizeCountryCode(dto.Code),
		Name:         NormalizeCountryName(dto.Name),
		OfficialName: NormalizeCountryName(dto.OfficialName),
		ISO3Code:     NormalizeCountryCode(dto.ISO3Code),
		IsActive:     true,
	}
}

// MapCountriesToDTOs converts slice using pure function composition
// Pure function - no side effects, immutable result
func MapCountriesToDTOs(countries []Country) []CountryDTO {
	return Map(countries, CountryToDTO)
}

// Map applies function to each element
// Pure higher-order function
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}