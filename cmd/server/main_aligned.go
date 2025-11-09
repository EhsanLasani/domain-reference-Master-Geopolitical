package main

import (
	"context"
	"database/sql"
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
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

func main() {
	ctx := context.Background()
	
	// Load configuration from environment
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "@Salman2021"),
			Name:     getEnv("DB_NAME", "geopolitical"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: config.RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Server: config.ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8082),
		},
	}
	
	// Initialize logger
	logger := logging.NewStructuredLogger("debug")
	
	logger.Info(ctx, "Starting Aligned Reference Master Geopolitical service",
		logging.Field{Key: "version", Value: "2.0.0-aligned"},
		logging.Field{Key: "port", Value: cfg.Server.Port},
		logging.Field{Key: "db_name", Value: cfg.Database.Name},
		logging.Field{Key: "schema", Value: "domain_reference_master_geopolitical"})

	// Initialize aligned container with all dependencies
	container, err := bootstrap.NewAlignedContainer(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize aligned container: %v", err)
	}
	defer container.Close()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(1000, 100)
	router.Use(rateLimiter.Middleware())
	
	// CORS middleware
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
	
	// Tenant context middleware
	router.Use(func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default-tenant"
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	})

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "reference-master-geopolitical-aligned",
			"version":   "2.0.0-aligned",
			"timestamp": time.Now().UTC(),
			"schema":    "domain_reference_master_geopolitical",
		})
	})
	
	router.GET("/health/db", func(c *gin.Context) {
		sqlDB, err := container.DBManager.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
			"schema": "domain_reference_master_geopolitical",
		})
	})

	// Static file serving for web UI
	router.Static("/web-ui", "./presentation-layer/web-ui")
	
	// Redirect root to web UI
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web-ui/enterprise-demo.html")
	})

	// API v2 routes (aligned)
	v2 := router.Group("/api/v2")
	{
		// Countries endpoints using aligned service
		countries := v2.Group("/countries")
		{
			countries.GET("", func(c *gin.Context) {
				tenantID := c.GetString("tenant_id")
				
				countries, err := container.CountryAppService.GetAllCountries(c.Request.Context(), tenantID)
				if err != nil {
					logger.Error(c.Request.Context(), "Failed to get countries", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve countries"})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
					"countries": countries,
					"count":     len(countries),
					"schema":    "domain_reference_master_geopolitical",
					"version":   "2.0.0-aligned",
				})
			})
			
			countries.POST("", func(c *gin.Context) {
				tenantID := c.GetString("tenant_id")
				
				var country models.Country
				if err := c.ShouldBindJSON(&country); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
					return
				}
				
				if err := container.CountryAppService.CreateCountry(c.Request.Context(), tenantID, &country); err != nil {
					logger.Error(c.Request.Context(), "Failed to create country", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusCreated, gin.H{
					"message": "Country created successfully",
					"country": country,
					"schema":  "domain_reference_master_geopolitical",
				})
			})
			
			countries.GET("/:code", func(c *gin.Context) {
				tenantID := c.GetString("tenant_id")
				code := c.Param("code")
				
				country, err := container.CountryAppService.GetCountryByCode(c.Request.Context(), tenantID, code)
				if err != nil {
					logger.Error(c.Request.Context(), "Failed to get country by code", err)
					c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
					"country": country,
					"schema":  "domain_reference_master_geopolitical",
				})
			})
			
			countries.PUT("/:code", func(c *gin.Context) {
				tenantID := c.GetString("tenant_id")
				code := c.Param("code")
				
				// Get existing country
				country, err := container.CountryAppService.GetCountryByCode(c.Request.Context(), tenantID, code)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
					return
				}
				
				// Bind update data
				var updateData models.Country
				if err := c.ShouldBindJSON(&updateData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
					return
				}
				
				// Update fields
				country.CountryName = updateData.CountryName
				country.ISO3Code = updateData.ISO3Code
				country.NumericCode = updateData.NumericCode
				country.OfficialName = updateData.OfficialName
				country.CapitalCity = updateData.CapitalCity
				country.ContinentCode = updateData.ContinentCode
				country.PhonePrefix = updateData.PhonePrefix
				country.IsActive = updateData.IsActive
				
				if err := container.CountryAppService.UpdateCountry(c.Request.Context(), tenantID, country); err != nil {
					logger.Error(c.Request.Context(), "Failed to update country", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
					"message": "Country updated successfully",
					"country": country,
					"schema":  "domain_reference_master_geopolitical",
				})
			})
			
			countries.DELETE("/:code", func(c *gin.Context) {
				tenantID := c.GetString("tenant_id")
				code := c.Param("code")
				
				if err := container.CountryAppService.DeleteCountry(c.Request.Context(), tenantID, code); err != nil {
					logger.Error(c.Request.Context(), "Failed to delete country", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
					"message": "Country deleted successfully",
					"code":    code,
					"schema":  "domain_reference_master_geopolitical",
				})
			})
		}
		
		// Schema info endpoint
		v2.GET("/schema", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"schema": "domain_reference_master_geopolitical",
				"version": "2.0.0-aligned",
				"entities": []string{
					"countries", "regions", "languages", 
					"timezones", "country_subdivisions", "locales",
				},
				"alignment_status": "✅ Fully Aligned",
				"lasani_compliance": "27+ fields per entity",
			})
		})
	}

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info(ctx, "Aligned server started successfully",
		logging.Field{Key: "port", Value: cfg.Server.Port},
		logging.Field{Key: "schema", Value: "domain_reference_master_geopolitical"},
		logging.Field{Key: "alignment", Value: "✅ Complete"})

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down aligned server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info(ctx, "Aligned server exited")
}

// Helper functions
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

func nullStringToString(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func nullInt32ToInt(ni sql.NullInt32) interface{} {
	if ni.Valid {
		return ni.Int32
	}
	return nil
}