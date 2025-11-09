// Package bootstrap provides dependency injection and application container
package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/cache"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/database"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/models"
	applicationservices "github.com/EhsanLasani/domain-reference-Master-Geopolitical/business-logic-layer/application-services"
	repositories "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/repositories-daos"
)

// Container holds all application dependencies
type Container struct {
	Config            *config.Config
	Logger            logging.Logger
	Tracer            tracing.Tracer
	Cache             cache.Cache
	DB                *sql.DB
	DBManager         *database.Database
	CountryRepository repositories.CountryRepositoryInterface
	CountryAppService *applicationservices.CountryAppService
	RegionRepository  *repositories.RegionRepository
	LanguageRepository *repositories.LanguageRepository
}

// NewContainer creates and initializes the application container
func NewContainer(cfg *config.Config, logger logging.Logger) (*Container, error) {
	// Initialize tracing
	tracer, err := tracing.NewTracer("geopolitical-service")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize cache
	cacheClient := cache.NewCache(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB)

	// Initialize database
	dbManager, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Validate existing table structure
	if err := dbManager.CheckAllTables(); err != nil {
		fmt.Printf("Table validation completed with warnings: %v\n", err)
	}
	
	// Auto-migrate all models
	if err := dbManager.AutoMigrate(
		&models.Country{},
		&models.Region{},
		&models.Language{},
		&models.Timezone{},
		&models.CountrySubdivision{},
		&models.Locales{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// Initialize repositories with validation
	countryRepo := repositories.NewValidatedCountryRepository(dbManager.DB)
	regionRepo := repositories.NewRegionRepository(dbManager.DB)
	languageRepo := repositories.NewLanguageRepository(dbManager.DB)

	// Initialize application services
	countryAppService := applicationservices.NewCountryAppService(countryRepo, logger, tracer)

	return &Container{
		Config:            cfg,
		Logger:            logger,
		Tracer:            tracer,
		Cache:             cacheClient,
		DB:                nil, // Will be set from GORM DB if needed
		DBManager:         dbManager,
		CountryRepository: countryRepo,
		CountryAppService: countryAppService,
		RegionRepository:  regionRepo,
		LanguageRepository: languageRepo,
	}, nil
}

// Close closes all resources
func (c *Container) Close() error {
	if c.DBManager != nil {
		return c.DBManager.Close()
	}
	return nil
}