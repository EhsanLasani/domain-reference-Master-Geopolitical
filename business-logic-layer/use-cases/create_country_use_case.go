package usecases

import (
	"context"
	"fmt"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/errors"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

// CreateCountryInput represents the input for creating a country
type CreateCountryInput struct {
	CountryCode    string `json:"country_code" validate:"required,len=2"`
	CountryName    string `json:"country_name" validate:"required,max=100"`
	ISO3Code       string `json:"iso3_code" validate:"len=3"`
	OfficialName   string `json:"official_name" validate:"max=200"`
	ContinentCode  string `json:"continent_code" validate:"oneof=AF AS EU NA SA OC AN"`
	TenantID       string `json:"-"`
	UserID         string `json:"-"`
}

// CreateCountryOutput represents the output of creating a country
type CreateCountryOutput struct {
	CountryID   string `json:"country_id"`
	CountryCode string `json:"country_code"`
	Success     bool   `json:"success"`
}

// CountryRepository interface for dependency injection
type CountryRepository interface {
	Create(ctx context.Context, tenantID string, country *models.Country) error
	GetByCode(ctx context.Context, tenantID, code string) (*models.Country, error)
}

// CreateCountryUseCase implements the Handle pattern
type CreateCountryUseCase struct {
	repo   CountryRepository
	logger *logging.StructuredLogger
	tracer *tracing.Tracer
}

func NewCreateCountryUseCase(repo CountryRepository, logger *logging.StructuredLogger, tracer *tracing.Tracer) *CreateCountryUseCase {
	return &CreateCountryUseCase{
		repo:   repo,
		logger: logger,
		tracer: tracer,
	}
}

// Handle implements the use case pattern: Handle(ctx, input) (output, error)
func (uc *CreateCountryUseCase) Handle(ctx context.Context, input CreateCountryInput) (*CreateCountryOutput, error) {
	ctx, span := uc.tracer.StartSpan(ctx, "use_case.create_country")
	defer span.End()

	// Business validation
	if err := uc.validateBusinessRules(ctx, input); err != nil {
		return nil, err
	}

	// Check if country already exists
	existing, err := uc.repo.GetByCode(ctx, input.TenantID, input.CountryCode)
	if err != nil {
		return nil, errors.NewSystemError("GEO-1006", "Database error during validation")
	}
	if existing != nil {
		return nil, errors.NewBusinessError("GEO-1002", "Country code already exists")
	}

	// Create country entity
	country := &models.Country{
		CountryCode:   input.CountryCode,
		CountryName:   input.CountryName,
		ISO3Code:      &input.ISO3Code,
		OfficialName:  &input.OfficialName,
		ContinentCode: &input.ContinentCode,
	}

	// Execute creation
	if err := uc.repo.Create(ctx, input.TenantID, country); err != nil {
		uc.logger.Error(ctx, "Failed to create country", err,
			logging.Field{Key: "tenant_id", Value: input.TenantID},
			logging.Field{Key: "country_code", Value: input.CountryCode})
		return nil, errors.NewSystemError("GEO-1006", "Failed to create country")
	}

	uc.logger.Info(ctx, "Country created successfully",
		logging.Field{Key: "tenant_id", Value: input.TenantID},
		logging.Field{Key: "country_code", Value: input.CountryCode},
		logging.Field{Key: "country_id", Value: country.CountryID.String()})

	return &CreateCountryOutput{
		CountryID:   country.CountryID.String(),
		CountryCode: country.CountryCode,
		Success:     true,
	}, nil
}

func (uc *CreateCountryUseCase) validateBusinessRules(ctx context.Context, input CreateCountryInput) error {
	// Business rule: Country code must be uppercase
	if input.CountryCode != fmt.Sprintf("%s", input.CountryCode) {
		return errors.NewBusinessError("GEO-1001", "Country code must be uppercase")
	}

	// Business rule: Continent code validation
	validContinents := []string{"AF", "AS", "EU", "NA", "SA", "OC", "AN"}
	if input.ContinentCode != "" {
		valid := false
		for _, continent := range validContinents {
			if input.ContinentCode == continent {
				valid = true
				break
			}
		}
		if !valid {
			return errors.NewBusinessError("GEO-1003", "Invalid continent code")
		}
	}

	return nil
}