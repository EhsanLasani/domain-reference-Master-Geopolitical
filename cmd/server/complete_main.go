package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/bootstrap"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/presentation-layer/middleware"
	v1 "github.com/EhsanLasani/domain-reference-Master-Geopolitical/presentation-layer/rest-api/v1"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/repositories-daos"
)

func main() {
	ctx := context.Background()
	
	// Load configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "@Salman2021"),
			Name:     getEnv("DB_NAME", "referencemaster"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: config.ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8081),
		},
	}
	
	// Initialize logger
	logger := logging.NewStructuredLogger("debug")
	
	// Initialize container
	container, err := bootstrap.NewContainer(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Initialize all repositories
	regionRepo := repositories.NewRegionRepository(container.DBManager.DB)
	languageRepo := repositories.NewLanguageRepository(container.DBManager.DB)
	timezoneRepo := repositories.NewTimezoneRepository(container.DBManager.DB)
	subdivisionRepo := repositories.NewSubdivisionRepository(container.DBManager.DB)
	localeRepo := repositories.NewLocaleRepository(container.DBManager.DB)

	// Initialize all handlers
	countriesHandler := v1.NewCountriesHandler(container.CountryAppService, container.Logger)
	regionsHandler := v1.NewRegionsHandler(regionRepo)
	languagesHandler := v1.NewLanguagesHandler(languageRepo)
	timezonesHandler := v1.NewTimezonesHandler(timezoneRepo)
	subdivisionsHandler := v1.NewSubdivisionsHandler(subdivisionRepo)
	localesHandler := v1.NewLocalesHandler(localeRepo)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(1000, 100)
	router.Use(rateLimiter.Middleware())
	
	// CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Tenant-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	
	// Tenant context
	router.Use(func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default-tenant"
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	})

	// Health endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "reference-master-geopolitical",
			"version": "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	// Static files
	router.Static("/web-ui", "./presentation-layer/web-ui")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web-ui/business-central.html")
	})

	// API v1 routes - ALL TABLES CRUD
	v1Group := router.Group("/api/v1")
	{
		// Countries CRUD
		countries := v1Group.Group("/countries")
		{
			countries.GET("", countriesHandler.GetAllCountries)
			countries.POST("", countriesHandler.CreateCountry)
			countries.GET("/:code", countriesHandler.GetCountryByCode)
			countries.PUT("/:code", countriesHandler.UpdateCountry)
			countries.DELETE("/:code", countriesHandler.DeleteCountry)
		}

		// Regions CRUD
		regions := v1Group.Group("/regions")
		{
			regions.GET("", regionsHandler.GetAll)
			regions.POST("", regionsHandler.Create)
			regions.GET("/:code", regionsHandler.GetByCode)
			regions.PUT("/:code", regionsHandler.Update)
			regions.DELETE("/:code", regionsHandler.Delete)
		}

		// Languages CRUD
		languages := v1Group.Group("/languages")
		{
			languages.GET("", languagesHandler.GetAll)
			languages.POST("", languagesHandler.Create)
			languages.GET("/:code", languagesHandler.GetByCode)
			languages.PUT("/:code", languagesHandler.Update)
			languages.DELETE("/:code", languagesHandler.Delete)
		}

		// Timezones CRUD
		timezones := v1Group.Group("/timezones")
		{
			timezones.GET("", timezonesHandler.GetAll)
			timezones.POST("", timezonesHandler.Create)
			timezones.GET("/:code", timezonesHandler.GetByCode)
			timezones.PUT("/:code", timezonesHandler.Update)
			timezones.DELETE("/:code", timezonesHandler.Delete)
		}

		// Subdivisions CRUD
		subdivisions := v1Group.Group("/subdivisions")
		{
			subdivisions.GET("", subdivisionsHandler.GetAll)
			subdivisions.POST("", subdivisionsHandler.Create)
			subdivisions.GET("/country/:countryId", subdivisionsHandler.GetByCountry)
			subdivisions.PUT("/:id", subdivisionsHandler.Update)
			subdivisions.DELETE("/:id", subdivisionsHandler.Delete)
		}

		// Locales CRUD
		locales := v1Group.Group("/locales")
		{
			locales.GET("", localesHandler.GetAll)
			locales.POST("", localesHandler.Create)
			locales.GET("/:code", localesHandler.GetByCode)
			locales.PUT("/:code", localesHandler.Update)
			locales.DELETE("/:code", localesHandler.Delete)
		}
	}

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info(ctx, "🚀 Complete CRUD API Server started",
		logging.Field{Key: "port", Value: cfg.Server.Port},
		logging.Field{Key: "endpoints", Value: "countries, regions, languages, timezones, subdivisions, locales"})

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info(ctx, "Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}