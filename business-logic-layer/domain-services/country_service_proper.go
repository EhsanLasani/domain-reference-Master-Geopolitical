// BUSINESS LOGIC LAYER - Domain Service
package services

import (
	"reference-master-geopolitical/data-access-layer/orm-odm-abstractions"
	"reference-master-geopolitical/data-access-layer/repositories-daos"
)

type CountryService struct {
	repo *repositories.CountryRepository
}

func NewCountryService(repo *repositories.CountryRepository) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) GetActiveCountries() ([]models.Country, error) {
	// Business logic: validate, apply rules, etc.
	countries, err := s.repo.GetAllActiveCountries()
	if err != nil {
		return nil, err
	}
	
	// Business rule: Filter out countries without proper codes
	var validCountries []models.Country
	for _, country := range countries {
		if len(country.CountryCode) == 2 && country.CountryName != "" {
			validCountries = append(validCountries, country)
		}
	}
	
	return validCountries, nil
}