// DATA ACCESS LAYER - Models
package models

type Country struct {
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	ISO3Code     string `json:"iso3_code"`
	OfficialName string `json:"official_name"`
	IsActive     bool   `json:"is_active"`
}