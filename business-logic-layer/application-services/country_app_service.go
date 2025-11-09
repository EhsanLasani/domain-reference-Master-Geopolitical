// Package applicationservices implements use case orchestration and business logic coordination
package applicationservices

import (
	"context"
	"fmt"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/errors"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
	repositories "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/repositories-daos"
)

// CountryAppService orchestrates country-related use cases
type CountryAppService struct {
	countryRepo repositories.CountryRepositoryInterface
	logger      logging.Logger
	tracer      tracing.Tracer
}

// NewCountryAppService creates a new country application service
func NewCountryAppService(
	countryRepo repositories.CountryRepositoryInterface,
	logger logging.Logger,
	tracer tracing.Tracer,
) *CountryAppService {
	return &CountryAppService{
		countryRepo: countryRepo,
		logger:      logger,
		tracer:      tracer,
	}
}

// GetAllCountries retrieves all active countries
func (s *CountryAppService) GetAllCountries(ctx context.Context, tenantID string) ([]models.Country, error) {
	ctx, span := s.tracer.StartSpan(ctx, "CountryAppService.GetAllCountries")
	defer span.End()

	s.logger.Info(ctx, "Retrieving all countries", 
		logging.Field{Key: "operation", Value: "get_all_countries"})

	// Reference data - available to all tenants
	countries, err := s.countryRepo.GetAllActiveCountries(ctx, tenantID)
	if err != nil {
		s.logger.Error(ctx, "Failed to retrieve countries", err,
			logging.Field{Key: "error", Value: err.Error()})
		return nil, errors.MapDatabaseError(err)
	}

	s.logger.Info(ctx, "Successfully retrieved countries",
		logging.Field{Key: "count", Value: len(countries)})

	return countries, nil
}

// CreateCountry creates a new country
func (s *CountryAppService) CreateCountry(ctx context.Context, tenantID string, country *models.Country) error {
	ctx, span := s.tracer.StartSpan(ctx, "CountryAppService.CreateCountry")
	defer span.End()

	s.logger.Info(ctx, "Creating new country",
		logging.Field{Key: "country_code", Value: country.CountryCode},
		logging.Field{Key: "operation", Value: "create_country"})

	// Business validation
	if country.CountryCode == "" || country.CountryName == "" {
		return errors.NewValidationError("country_code", "country code and name are required")
	}

	// Reference data - no tenant context needed
	err := s.countryRepo.Create(ctx, tenantID, country)
	if err != nil {
		s.logger.Error(ctx, "Failed to create country", err,
			logging.Field{Key: "country_code", Value: country.CountryCode})
		return errors.MapDatabaseError(err)
	}

	s.logger.Info(ctx, "Successfully created country",
		logging.Field{Key: "country_code", Value: country.CountryCode},
		logging.Field{Key: "country_id", Value: country.CountryID})

	return nil
}

// GetCountryByCode retrieves a country by its code
func (s *CountryAppService) GetCountryByCode(ctx context.Context, tenantID, countryCode string) (*models.Country, error) {
	ctx, span := s.tracer.StartSpan(ctx, "CountryAppService.GetCountryByCode")
	defer span.End()

	s.logger.Info(ctx, "Retrieving country by code",
		logging.Field{Key: "country_code", Value: countryCode},
		logging.Field{Key: "operation", Value: "get_country_by_code"})

	// Business validation
	if countryCode == "" {
		return nil, errors.NewValidationError("country_code", "country code is required")
	}

	// Reference data lookup
	country, err := s.countryRepo.GetByCode(ctx, tenantID, countryCode)
	if err != nil {
		s.logger.Error(ctx, "Failed to retrieve country by code", err,
			logging.Field{Key: "country_code", Value: countryCode})
		return nil, errors.MapDatabaseError(err)
	}

	if country == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("country with code %s", countryCode))
	}

	s.logger.Info(ctx, "Successfully retrieved country by code",
		logging.Field{Key: "country_code", Value: countryCode},
		logging.Field{Key: "country_id", Value: country.CountryID})

	return country, nil
}