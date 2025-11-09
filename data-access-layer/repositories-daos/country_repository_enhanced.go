package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/cache"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

// EnhancedCountryRepository implements enterprise guidelines
type EnhancedCountryRepository struct {
	db     *gorm.DB
	cache  cache.Cache
	logger logging.Logger
	tracer tracing.Tracer
}

func NewEnhancedCountryRepository(db *gorm.DB, cache cache.Cache, logger logging.Logger, tracer tracing.Tracer) *EnhancedCountryRepository {
	return &EnhancedCountryRepository{
		db:     db,
		cache:  cache,
		logger: logger,
		tracer: tracer,
	}
}

// GetAllActiveCountries with tenant isolation and tracing
func (r *EnhancedCountryRepository) GetAllActiveCountries(ctx context.Context, tenantID string) ([]models.Country, error) {
	ctx, span := r.tracer.StartSQLSpan(ctx, "countries.list_active", "SELECT", tenantID)
	defer span.End()

	// Check cache first
	cacheKey := "countries:active"
	var countries []models.Country
	if err := r.cache.Get(ctx, tenantID, cacheKey, &countries); err == nil {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		return countries, nil
	}

	// Set tenant context for RLS
	if err := r.db.Exec("SET app.tenant_id = ?", tenantID).Error; err != nil {
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}

	start := time.Now()
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND is_deleted = ?", true, false).
		Order("country_name ASC").
		Find(&countries).Error

	duration := time.Since(start)
	span.SetAttributes(
		attribute.Int64("query_duration_ms", duration.Milliseconds()),
		attribute.Int("result_count", len(countries)),
		attribute.Bool("cache_hit", false),
	)

	if err != nil {
		r.logger.Error(ctx, "Failed to query countries", err,
			logging.Field{Key: "tenant_id", Value: tenantID},
			logging.Field{Key: "sql_key", Value: "countries.list_active"},
			logging.Field{Key: "duration_ms", Value: duration.Milliseconds()})
		return nil, fmt.Errorf("GEO-1006: %w", err)
	}

	// Cache results for 5 minutes
	r.cache.Set(ctx, tenantID, cacheKey, countries, 5*time.Minute)

	return countries, nil
}

// Create with LASANI audit fields
func (r *EnhancedCountryRepository) Create(ctx context.Context, tenantID string, country *models.Country) error {
	ctx, span := r.tracer.StartSQLSpan(ctx, "countries.create", "INSERT", tenantID)
	defer span.End()

	// Set LASANI audit fields
	now := time.Now()
	country.CountryID = uuid.New()
	country.CreatedAt = &now
	country.UpdatedAt = &now
	country.IsActive = true
	country.IsDeleted = false
	country.Version = 1
	country.SourceSystem = "reference_master_geopolitical"

	// Set tenant context
	if err := r.db.Exec("SET app.tenant_id = ?", tenantID).Error; err != nil {
		return fmt.Errorf("failed to set tenant context: %w", err)
	}

	start := time.Now()
	err := r.db.WithContext(ctx).Create(country).Error
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("query_duration_ms", duration.Milliseconds()),
		attribute.String("country_code", country.CountryCode),
	)

	if err != nil {
		r.logger.Error(ctx, "Failed to create country", err,
			logging.Field{Key: "tenant_id", Value: tenantID},
			logging.Field{Key: "country_code", Value: country.CountryCode},
			logging.Field{Key: "sql_key", Value: "countries.create"})
		return fmt.Errorf("GEO-1002: %w", err)
	}

	// Invalidate cache
	r.cache.Delete(ctx, tenantID, "countries:active")

	return nil
}