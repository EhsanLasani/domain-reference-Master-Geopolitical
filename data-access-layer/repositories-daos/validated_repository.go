package repositories

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/validation"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/errors"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

type ValidatedCountryRepository struct {
	db        *gorm.DB
	validator *validation.SchemaValidator
}

func NewValidatedCountryRepository(db *gorm.DB) CountryRepositoryInterface {
	return &ValidatedCountryRepository{
		db:        db,
		validator: validation.NewSchemaValidator(db),
	}
}

func (r *ValidatedCountryRepository) GetAllActiveCountries(ctx context.Context, tenantID string) ([]models.Country, error) {
	var countries []models.Country
	
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ? AND is_deleted = ?", tenantID, true, false).
		Order("country_name ASC").
		Find(&countries).Error
	
	if err != nil {
		return nil, errors.NewRepositoryError("QUERY_FAILED", "Failed to retrieve countries", err)
	}
	
	return countries, nil
}

func (r *ValidatedCountryRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Country, error) {
	// Input validation
	if code == "" {
		return nil, errors.NewValidationError("country_code", "cannot be empty")
	}
	if len(code) != 2 {
		return nil, errors.NewValidationError("country_code", "must be exactly 2 characters")
	}
	
	var country models.Country
	err := r.db.WithContext(ctx).
		Where("country_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).
		First(&country).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.NewRepositoryError("QUERY_FAILED", "Failed to retrieve country", err)
	}
	
	return &country, nil
}

func (r *ValidatedCountryRepository) Create(ctx context.Context, tenantID string, country *models.Country) error {
	// Schema validation
	if err := r.validator.ValidateStruct("countries", country); err != nil {
		return errors.NewSchemaError("Schema validation failed", err)
	}
	
	// Business validation
	if country.CountryCode == "" {
		return errors.NewValidationError("country_code", "is required")
	}
	if country.CountryName == "" {
		return errors.NewValidationError("country_name", "is required")
	}
	if len(country.CountryCode) != 2 {
		return errors.NewValidationError("country_code", "must be exactly 2 characters")
	}
	
	// Set audit fields
	now := time.Now()
	country.CountryID = uuid.New()
	country.TenantID = tenantID
	country.CreatedAt = &now
	country.UpdatedAt = &now
	country.IsActive = true
	country.IsDeleted = false
	country.Version = 1
	country.SourceSystem = "reference_master_geopolitical"
	
	// Database operation
	if err := r.db.WithContext(ctx).Create(country).Error; err != nil {
		// Map database errors
		if isDuplicateKeyError(err) {
			return errors.NewRepositoryError("DUPLICATE_KEY", "Country code already exists", err)
		}
		if isConstraintViolationError(err) {
			return errors.NewRepositoryError("CONSTRAINT_VIOLATION", "Data violates database constraints", err)
		}
		return errors.NewRepositoryError("CREATE_FAILED", "Failed to create country", err)
	}
	
	return nil
}

func (r *ValidatedCountryRepository) Update(ctx context.Context, tenantID string, country *models.Country) error {
	// Schema validation
	if err := r.validator.ValidateStruct("countries", country); err != nil {
		return errors.NewSchemaError("Schema validation failed", err)
	}
	
	// Business validation
	if country.CountryCode == "" {
		return errors.NewValidationError("country_code", "is required")
	}
	if country.CountryName == "" {
		return errors.NewValidationError("country_name", "is required")
	}
	
	// Set audit fields
	now := time.Now()
	country.UpdatedAt = &now
	
	result := r.db.WithContext(ctx).
		Where("country_code = ? AND tenant_id = ? AND version = ?", country.CountryCode, tenantID, country.Version).
		Updates(map[string]interface{}{
			"country_name":   country.CountryName,
			"iso3_code":      country.ISO3Code,
			"official_name":  country.OfficialName,
			"capital_city":   country.CapitalCity,
			"continent_code": country.ContinentCode,
			"phone_prefix":   country.PhonePrefix,
			"is_active":      country.IsActive,
			"updated_at":     now,
			"version":        gorm.Expr("version + 1"),
		})
	
	if result.Error != nil {
		if isConstraintViolationError(result.Error) {
			return errors.NewRepositoryError("CONSTRAINT_VIOLATION", "Data violates database constraints", result.Error)
		}
		return errors.NewRepositoryError("UPDATE_FAILED", "Failed to update country", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.NewRepositoryError("NOT_FOUND", "Country not found or version conflict", nil)
	}
	
	return nil
}

func (r *ValidatedCountryRepository) Delete(ctx context.Context, tenantID, code string) error {
	// Input validation
	if code == "" {
		return errors.NewValidationError("country_code", "cannot be empty")
	}
	
	now := time.Now()
	result := r.db.WithContext(ctx).
		Where("country_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at":  now,
			"updated_at":  now,
			"version":     gorm.Expr("version + 1"),
		})
	
	if result.Error != nil {
		return errors.NewRepositoryError("DELETE_FAILED", "Failed to delete country", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.NewRepositoryError("NOT_FOUND", "Country not found", nil)
	}
	
	return nil
}

func (r *ValidatedCountryRepository) BulkCreate(ctx context.Context, tenantID string, countries []models.Country) error {
	// Validate each country
	for i, country := range countries {
		if err := r.validator.ValidateStruct("countries", &country); err != nil {
			return errors.NewSchemaError(fmt.Sprintf("Validation failed for country %d", i), err)
		}
	}
	
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		
		for i := range countries {
			countries[i].CountryID = uuid.New()
			countries[i].TenantID = tenantID
			countries[i].CreatedAt = &now
			countries[i].UpdatedAt = &now
			countries[i].IsActive = true
			countries[i].IsDeleted = false
			countries[i].Version = 1
			countries[i].SourceSystem = "reference_master_geopolitical"
		}
		
		if err := tx.CreateInBatches(countries, 100).Error; err != nil {
			return errors.NewRepositoryError("BULK_CREATE_FAILED", "Failed to bulk create countries", err)
		}
		
		return nil
	})
}

// Helper functions for error detection
func isDuplicateKeyError(err error) bool {
	return err != nil && (
		containsString(err.Error(), "duplicate key") ||
		containsString(err.Error(), "UNIQUE constraint"))
}

func isConstraintViolationError(err error) bool {
	return err != nil && (
		containsString(err.Error(), "violates") ||
		containsString(err.Error(), "constraint"))
}

func containsString(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     (str[:len(substr)] == substr || 
		      str[len(str)-len(substr):] == substr ||
		      findSubstring(str, substr))))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}