// Package applicationservices implements application layer orchestration
package applicationservices

import (
	"context"
	"fmt"

	"github.com/domain-reference-Master-Geopolitical/internal/xcut/errors"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/tracing"
)

type CountryApplicationService struct {
	countryDomainService CountryDomainService
	countryRepository    CountryRepository
	logger              logging.Logger
	tracer              tracing.Tracer
}

type CountryDTO struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	TenantID     string `json:"tenant_id"`
}

type CreateCountryRequest struct {
	Code     string `json:"code" validate:"required,len=2"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Region   string `json:"region" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
}

type UpdateCountryRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Region   string `json:"region" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
}

// Interfaces for dependencies
type CountryDomainService interface {
	ValidateCountryCreation(ctx context.Context, country *Country) error
	ValidateCountryUpdate(ctx context.Context, country *Country) error
}

type CountryRepository interface {
	Create(ctx context.Context, country *Country) error
	Update(ctx context.Context, country *Country) error
	FindByID(ctx context.Context, id string, tenantID string) (*Country, error)
	FindByCode(ctx context.Context, code string, tenantID string) (*Country, error)
	FindByTenant(ctx context.Context, tenantID string) ([]*Country, error)
	Delete(ctx context.Context, id string, tenantID string) error
}

type Country struct {
	ID       string
	Code     string
	Name     string
	Region   string
	TenantID string
}

func NewCountryApplicationService(
	domainService CountryDomainService,
	repository CountryRepository,
	logger logging.Logger,
	tracer tracing.Tracer,
) *CountryApplicationService {
	return &CountryApplicationService{
		countryDomainService: domainService,
		countryRepository:    repository,
		logger:              logger,
		tracer:              tracer,
	}
}

func (s *CountryApplicationService) CreateCountry(ctx context.Context, req *CreateCountryRequest) (*CountryDTO, error) {
	ctx, span := s.tracer.StartUseCaseSpan(ctx, "create_country")
	defer span.End()

	s.logger.Info(ctx, "Creating country", 
		logging.Field{Key: "country_code", Value: req.Code},
		logging.Field{Key: "tenant_id", Value: req.TenantID},
	)

	// Check if country already exists
	existing, err := s.countryRepository.FindByCode(ctx, req.Code, req.TenantID)
	if err == nil && existing != nil {
		return nil, errors.NewValidationError("Country with this code already exists")
	}

	// Create domain entity
	country := &Country{
		Code:     req.Code,
		Name:     req.Name,
		Region:   req.Region,
		TenantID: req.TenantID,
	}

	// Domain validation
	if err := s.countryDomainService.ValidateCountryCreation(ctx, country); err != nil {
		s.logger.Error(ctx, "Country validation failed", err,
			logging.Field{Key: "country_code", Value: req.Code},
		)
		return nil, err
	}

	// Persist
	if err := s.countryRepository.Create(ctx, country); err != nil {
		s.logger.Error(ctx, "Failed to create country", err,
			logging.Field{Key: "country_code", Value: req.Code},
		)
		return nil, errors.MapDatabaseError(err)
	}

	s.logger.Info(ctx, "Country created successfully",
		logging.Field{Key: "country_id", Value: country.ID},
		logging.Field{Key: "country_code", Value: country.Code},
	)

	return s.mapToDTO(country), nil
}

func (s *CountryApplicationService) UpdateCountry(ctx context.Context, req *UpdateCountryRequest) (*CountryDTO, error) {
	ctx, span := s.tracer.StartUseCaseSpan(ctx, "update_country")
	defer span.End()

	s.logger.Info(ctx, "Updating country",
		logging.Field{Key: "country_id", Value: req.ID},
		logging.Field{Key: "tenant_id", Value: req.TenantID},
	)

	// Find existing country
	existing, err := s.countryRepository.FindByID(ctx, req.ID, req.TenantID)
	if err != nil {
		return nil, errors.MapDatabaseError(err)
	}
	if existing == nil {
		return nil, errors.NewNotFoundError("Country")
	}

	// Update fields
	existing.Name = req.Name
	existing.Region = req.Region

	// Domain validation
	if err := s.countryDomainService.ValidateCountryUpdate(ctx, existing); err != nil {
		return nil, err
	}

	// Persist
	if err := s.countryRepository.Update(ctx, existing); err != nil {
		s.logger.Error(ctx, "Failed to update country", err,
			logging.Field{Key: "country_id", Value: req.ID},
		)
		return nil, errors.MapDatabaseError(err)
	}

	s.logger.Info(ctx, "Country updated successfully",
		logging.Field{Key: "country_id", Value: existing.ID},
	)

	return s.mapToDTO(existing), nil
}

func (s *CountryApplicationService) GetCountry(ctx context.Context, id string, tenantID string) (*CountryDTO, error) {
	ctx, span := s.tracer.StartUseCaseSpan(ctx, "get_country")
	defer span.End()

	country, err := s.countryRepository.FindByID(ctx, id, tenantID)
	if err != nil {
		return nil, errors.MapDatabaseError(err)
	}
	if country == nil {
		return nil, errors.NewNotFoundError("Country")
	}

	return s.mapToDTO(country), nil
}

func (s *CountryApplicationService) ListCountries(ctx context.Context, tenantID string) ([]*CountryDTO, error) {
	ctx, span := s.tracer.StartUseCaseSpan(ctx, "list_countries")
	defer span.End()

	countries, err := s.countryRepository.FindByTenant(ctx, tenantID)
	if err != nil {
		return nil, errors.MapDatabaseError(err)
	}

	dtos := make([]*CountryDTO, len(countries))
	for i, country := range countries {
		dtos[i] = s.mapToDTO(country)
	}

	return dtos, nil
}

func (s *CountryApplicationService) DeleteCountry(ctx context.Context, id string, tenantID string) error {
	ctx, span := s.tracer.StartUseCaseSpan(ctx, "delete_country")
	defer span.End()

	s.logger.Info(ctx, "Deleting country",
		logging.Field{Key: "country_id", Value: id},
		logging.Field{Key: "tenant_id", Value: tenantID},
	)

	if err := s.countryRepository.Delete(ctx, id, tenantID); err != nil {
		s.logger.Error(ctx, "Failed to delete country", err,
			logging.Field{Key: "country_id", Value: id},
		)
		return errors.MapDatabaseError(err)
	}

	s.logger.Info(ctx, "Country deleted successfully",
		logging.Field{Key: "country_id", Value: id},
	)

	return nil
}

func (s *CountryApplicationService) mapToDTO(country *Country) *CountryDTO {
	return &CountryDTO{
		ID:       country.ID,
		Code:     country.Code,
		Name:     country.Name,
		Region:   country.Region,
		TenantID: country.TenantID,
	}
}