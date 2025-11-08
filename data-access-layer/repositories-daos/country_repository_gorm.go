// Package repositories implements data access patterns with GORM ORM
package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/cache"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

// CountryRepositoryInterface defines the contract for country data access
type CountryRepositoryInterface interface {
	GetAllActiveCountries(ctx context.Context, tenantID string) ([]models.Country, error)
	GetByCode(ctx context.Context, tenantID, code string) (*models.Country, error)
	Create(ctx context.Context, tenantID string, country *models.Country) error
	Update(ctx context.Context, tenantID string, country *models.Country) error
	Delete(ctx context.Context, tenantID, code string) error
	BulkCreate(ctx context.Context, tenantID string, countries []models.Country) error
}

// CountryRepositoryGORM implements CountryRepositoryInterface using GORM
type CountryRepositoryGORM struct {
	db     *gorm.DB
	cache  cache.Cache
	logger logging.Logger
}

// NewCountryRepositoryGORM creates a new GORM-based country repository
func NewCountryRepositoryGORM(db *gorm.DB, cache cache.Cache, logger logging.Logger) CountryRepositoryInterface {
	return &CountryRepositoryGORM{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// GetAllActiveCountries retrieves all active countries
func (r *CountryRepositoryGORM) GetAllActiveCountries(ctx context.Context, tenantID string) ([]models.Country, error) {
	var countries []models.Country
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND is_deleted = ?", true, false).
		Order("country_name ASC").
		Find(&countries).Error

	if err != nil {
		fmt.Printf("DEBUG: Database query error: %v\n", err)
		return nil, fmt.Errorf("failed to query countries: %w", err)
	}

	return countries, nil
}

// GetByCode retrieves a country by its code
func (r *CountryRepositoryGORM) GetByCode(ctx context.Context, tenantID, code string) (*models.Country, error) {
	var country models.Country
	err := r.db.WithContext(ctx).
		Where("country_code = ? AND is_active = ? AND is_deleted = ?", code, true, false).
		First(&country).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query country by code: %w", err)
	}

	return &country, nil
}

// Create creates a new country
func (r *CountryRepositoryGORM) Create(ctx context.Context, tenantID string, country *models.Country) error {
	now := time.Now()
	country.CountryID = uuid.New()
	country.CreatedAt = &now
	country.UpdatedAt = &now
	country.IsActive = true
	country.IsDeleted = false
	country.Version = 1
	country.SourceSystem = "reference_master"

	if err := r.db.WithContext(ctx).Create(country).Error; err != nil {
		return fmt.Errorf("failed to create country: %w", err)
	}

	return nil
}

// Update updates an existing country
func (r *CountryRepositoryGORM) Update(ctx context.Context, tenantID string, country *models.Country) error {
	now := time.Now()
	country.UpdatedAt = &now

	result := r.db.WithContext(ctx).
		Where("country_code = ? AND version = ?", country.CountryCode, country.Version).
		Updates(map[string]interface{}{
			"country_name":   country.CountryName,
			"iso3_code":      country.ISO3Code,
			"official_name":  country.OfficialName,
			"updated_at":     now,
			"version":        gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found or version conflict")
	}

	return nil
}

// Delete soft deletes a country
func (r *CountryRepositoryGORM) Delete(ctx context.Context, tenantID, code string) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Where("country_code = ? AND is_deleted = ?", code, false).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at":  now,
			"updated_at":  now,
			"version":     gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to delete country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

// BulkCreate creates multiple countries
func (r *CountryRepositoryGORM) BulkCreate(ctx context.Context, tenantID string, countries []models.Country) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		
		for i := range countries {
			countries[i].CountryID = uuid.New()
			countries[i].CreatedAt = &now
			countries[i].UpdatedAt = &now
			countries[i].IsActive = true
			countries[i].IsDeleted = false
			countries[i].Version = 1
			countries[i].SourceSystem = "reference_master"
		}

		if err := tx.CreateInBatches(countries, 100).Error; err != nil {
			return fmt.Errorf("failed to bulk create countries: %w", err)
		}

		return nil
	})
}