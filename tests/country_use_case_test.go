package tests

import (
	"context"
	"testing"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/business-logic-layer/use-cases"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

type MockCountryRepository struct {
	countries map[string]*models.Country
}

func (m *MockCountryRepository) Create(ctx context.Context, tenantID string, country *models.Country) error {
	m.countries[country.CountryCode] = country
	return nil
}

func (m *MockCountryRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Country, error) {
	if country, exists := m.countries[code]; exists {
		return country, nil
	}
	return nil, nil
}

func TestCreateCountryUseCase_Handle_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockCountryRepository{
		countries: make(map[string]*models.Country),
	}
	logger := logging.NewStructuredLogger("debug")
	tracer, _ := tracing.NewTracer("test")
	useCase := usecases.NewCreateCountryUseCase(mockRepo, logger, tracer)

	input := usecases.CreateCountryInput{
		CountryCode:   "US",
		CountryName:   "United States",
		ISO3Code:      "USA",
		OfficialName:  "United States of America",
		ContinentCode: "NA",
		TenantID:      "test-tenant",
		UserID:        "test-user",
	}

	// Act
	output, err := useCase.Handle(context.Background(), input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output, got nil")
	}
	if !output.Success {
		t.Error("Expected success to be true")
	}
	if output.CountryCode != "US" {
		t.Errorf("Expected country code US, got %s", output.CountryCode)
	}
}

func TestCreateCountryUseCase_Handle_DuplicateCode(t *testing.T) {
	// Arrange
	mockRepo := &MockCountryRepository{
		countries: make(map[string]*models.Country),
	}
	// Pre-populate with existing country
	mockRepo.countries["US"] = &models.Country{CountryCode: "US"}

	logger := logging.NewStructuredLogger("debug")
	tracer, _ := tracing.NewTracer("test")
	useCase := usecases.NewCreateCountryUseCase(mockRepo, logger, tracer)

	input := usecases.CreateCountryInput{
		CountryCode: "US",
		CountryName: "United States",
		TenantID:    "test-tenant",
	}

	// Act
	output, err := useCase.Handle(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate country code")
	}
	if output != nil {
		t.Error("Expected no output on error")
	}
}

func TestCreateCountryUseCase_Handle_InvalidCountryCode(t *testing.T) {
	// Arrange
	mockRepo := &MockCountryRepository{
		countries: make(map[string]*models.Country),
	}
	logger := logging.NewStructuredLogger("debug")
	tracer, _ := tracing.NewTracer("test")
	useCase := usecases.NewCreateCountryUseCase(mockRepo, logger, tracer)

	input := usecases.CreateCountryInput{
		CountryCode: "us", // lowercase - should fail validation
		CountryName: "United States",
		TenantID:    "test-tenant",
	}

	// Act
	output, err := useCase.Handle(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid country code format")
	}
	if output != nil {
		t.Error("Expected no output on validation error")
	}
}