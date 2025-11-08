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
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/bootstrap"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	countriesHandler "github.com/EhsanLasani/domain-reference-Master-Geopolitical/presentation-layer/rest-api/v1"
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
			Name:     getEnv("DB_NAME", "referencemaster"),
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
			Port: getEnvInt("SERVER_PORT", 8081),
		},
	}
	
	// Initialize logger
	logger := logging.NewStructuredLogger("debug")
	logger.Info(ctx, "Starting Reference Master Geopolitical service",
		logging.Field{Key: "version", Value: "1.0.0"},
		logging.Field{Key: "port", Value: cfg.Server.Port},
		logging.Field{Key: "db_name", Value: cfg.Database.Name},
		logging.Field{Key: "db_user", Value: cfg.Database.User})

	// Initialize container with all dependencies
	container, err := bootstrap.NewContainer(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
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
			"service":   "reference-master-geopolitical",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
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
		})
	})

	// Static file serving for web UI
	router.Static("/web-ui", "./presentation-layer/web-ui")
	
	// Redirect root to web UI
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web-ui/index.html")
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Countries endpoints
		countriesHandler := countriesHandler.NewCountriesHandler(
			container.CountryAppService,
			container.Logger,
		)
		
		countries := v1.Group("/countries")
		{
			countries.GET("", func(c *gin.Context) {
				// Get database connection
				db, err := container.DBManager.DB.DB()
				if err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection failed"})
					return
				}
				
				// Query all non-deleted countries (both active and inactive)
				rows, err := db.Query(`
					SELECT country_code, country_name, iso3_code, numeric_code, 
						   official_name, capital_city, continent_code, phone_prefix, 
						   is_active, created_at, updated_at
					FROM domain_reference_master_geopolitical.countries 
					WHERE tenant_id = 'default-tenant' AND is_deleted = false
					ORDER BY is_active DESC, country_code
				`)
				
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
					return
				}
				defer rows.Close()
				
				var countries []gin.H
				for rows.Next() {
					var countryCode, countryName string
					var iso3Code, officialName, capitalCity, continentCode, phonePrefix sql.NullString
					var numericCode sql.NullInt32
					var isActive bool
					var createdAt, updatedAt interface{}
					
					err := rows.Scan(
						&countryCode,
						&countryName,
						&iso3Code,
						&numericCode,
						&officialName,
						&capitalCity,
						&continentCode,
						&phonePrefix,
						&isActive,
						&createdAt,
						&updatedAt,
					)
					
					if err != nil {
						logger.Error(c.Request.Context(), "Failed to scan country row", err)
						continue
					}
					
					// Handle NULL values properly
					country := gin.H{
						"country_code": countryCode,
						"country_name": countryName,
						"iso3_code": nullStringToString(iso3Code),
						"numeric_code": nullInt32ToInt(numericCode),
						"official_name": nullStringToString(officialName),
						"capital_city": nullStringToString(capitalCity),
						"continent_code": nullStringToString(continentCode),
						"phone_prefix": nullStringToString(phonePrefix),
						"is_active": isActive,
					}
					
					countries = append(countries, country)
				}
				
				// Debug logging
				logger.Info(c.Request.Context(), "Countries retrieved",
					logging.Field{Key: "total_count", Value: len(countries)})
				
				for _, country := range countries {
					logger.Info(c.Request.Context(), "Country data",
						logging.Field{Key: "country_code", Value: country["country_code"]},
						logging.Field{Key: "country_name", Value: country["country_name"]},
						logging.Field{Key: "is_active", Value: country["is_active"]})
				}
				
				c.JSON(http.StatusOK, gin.H{
					"countries": countries,
					"count":     len(countries),
				})
			})
			countries.POST("", countriesHandler.CreateCountry)
			countries.GET("/:code", countriesHandler.GetCountryByCode)
			countries.PUT("/:code", func(c *gin.Context) {
				// PRESENTATION LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "PUT /countries/:code started", 
					logging.Field{Key: "country_code", Value: c.Param("code")},
					logging.Field{Key: "layer", Value: "presentation"})
				
				code := c.Param("code")
				if code == "" {
					logger.Error(c.Request.Context(), "Missing country code parameter", nil,
						logging.Field{Key: "layer", Value: "presentation"},
						logging.Field{Key: "error_code", Value: "PRES-001"})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code": "PRES-001",
							"message": "Country code parameter is required",
							"layer": "presentation",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				var updateData map[string]interface{}
				if err := c.ShouldBindJSON(&updateData); err != nil {
					logger.Error(c.Request.Context(), "JSON binding failed", err,
						logging.Field{Key: "layer", Value: "presentation"},
						logging.Field{Key: "error_code", Value: "PRES-002"})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code": "PRES-002",
							"message": "Invalid JSON format: " + err.Error(),
							"layer": "presentation",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				// DATA ACCESS LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "Accessing database for country update",
					logging.Field{Key: "layer", Value: "data_access"})
				
				db, err := container.DBManager.DB.DB()
				if err != nil {
					logger.Error(c.Request.Context(), "Database connection failed", err,
						logging.Field{Key: "layer", Value: "data_access"},
						logging.Field{Key: "error_code", Value: "DAL-001"})
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"error": gin.H{
							"code": "DAL-001",
							"message": "Database connection unavailable",
							"layer": "data_access",
							"timestamp": time.Now().UTC(),
							"details": err.Error(),
						},
					})
					return
				}
				
				// PRESENTATION LAYER VALIDATION
				logger.Info(c.Request.Context(), "Validating input data",
					logging.Field{Key: "layer", Value: "presentation"})
				
				// Validate continent enum
				if continentCode, ok := updateData["continent_code"].(string); ok && continentCode != "" {
					validContinents := []string{"AF", "AS", "EU", "NA", "SA", "OC", "AN"}
					validContinent := false
					for _, valid := range validContinents {
						if continentCode == valid {
							validContinent = true
							break
						}
					}
					if !validContinent {
						logger.Error(c.Request.Context(), "Invalid continent code", nil,
							logging.Field{Key: "layer", Value: "presentation"},
							logging.Field{Key: "error_code", Value: "PRES-003"},
							logging.Field{Key: "continent_code", Value: continentCode})
						c.JSON(http.StatusBadRequest, gin.H{
							"error": gin.H{
								"code": "PRES-003",
								"message": "Invalid continent code. Must be one of: AF, AS, EU, NA, SA, OC, AN",
								"layer": "presentation",
								"field": "continent_code",
								"provided_value": continentCode,
								"valid_values": []string{"AF", "AS", "EU", "NA", "SA", "OC", "AN"},
								"timestamp": time.Now().UTC(),
							},
						})
						return
					}
				}
				
				// BUSINESS LOGIC LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "Validating business rules",
					logging.Field{Key: "layer", Value: "business_logic"})
				
				// Validate required fields
				if updateData["country_name"] == nil || updateData["country_name"] == "" {
					logger.Error(c.Request.Context(), "Business validation failed: country_name required", nil,
						logging.Field{Key: "layer", Value: "business_logic"},
						logging.Field{Key: "error_code", Value: "BIZ-001"})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code": "BIZ-001",
							"message": "Country name is required",
							"layer": "business_logic",
							"field": "country_name",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				// DATABASE LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "Executing database update",
					logging.Field{Key: "layer", Value: "database"},
					logging.Field{Key: "country_code", Value: code},
					logging.Field{Key: "update_data", Value: updateData})
				
				// Log each parameter being sent to database
				logger.Info(c.Request.Context(), "Database update parameters",
					logging.Field{Key: "country_name", Value: updateData["country_name"]},
					logging.Field{Key: "iso3_code", Value: updateData["iso3_code"]},
					logging.Field{Key: "numeric_code", Value: updateData["numeric_code"]},
					logging.Field{Key: "official_name", Value: updateData["official_name"]},
					logging.Field{Key: "capital_city", Value: updateData["capital_city"]},
					logging.Field{Key: "continent_code", Value: updateData["continent_code"]},
					logging.Field{Key: "phone_prefix", Value: updateData["phone_prefix"]},
					logging.Field{Key: "is_active", Value: updateData["is_active"]},
					logging.Field{Key: "where_code", Value: code})
				
				result, err := db.Exec(`
					UPDATE domain_reference_master_geopolitical.countries 
					SET country_name = $1, 
						iso3_code = $2,
						numeric_code = $3,
						official_name = $4,
						capital_city = $5,
						continent_code = $6,
						phone_prefix = $7,
						is_active = $8,
						updated_at = NOW(),
						updated_by = '00000000-0000-0000-0000-000000000001'::uuid
					WHERE country_code = $9 AND tenant_id = 'default-tenant'
				`, 
					updateData["country_name"],
					updateData["iso3_code"],
					updateData["numeric_code"],
					updateData["official_name"],
					updateData["capital_city"],
					updateData["continent_code"],
					updateData["phone_prefix"],
					updateData["is_active"],
					code)
				
				if err != nil {
					logger.Error(c.Request.Context(), "Database update execution failed", err,
						logging.Field{Key: "layer", Value: "database"},
						logging.Field{Key: "error_code", Value: "DB-001"},
						logging.Field{Key: "country_code", Value: code},
						logging.Field{Key: "sql_error", Value: err.Error()})
					
					// Map specific database errors
					errorCode := "DB-001"
					errorMessage := "Database update failed"
					
					if strings.Contains(err.Error(), "duplicate key") {
						errorCode = "DB-002"
						errorMessage = "Duplicate country code"
					} else if strings.Contains(err.Error(), "foreign key") {
						errorCode = "DB-003"
						errorMessage = "Invalid reference data"
					} else if strings.Contains(err.Error(), "check constraint") {
						errorCode = "DB-004"
						errorMessage = "Data validation constraint failed"
					}
					
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code": errorCode,
							"message": errorMessage,
							"layer": "database",
							"timestamp": time.Now().UTC(),
							"details": err.Error(),
							"country_code": code,
						},
					})
					return
				}
				
				// Check if any rows were affected
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					logger.Error(c.Request.Context(), "Failed to get rows affected", err,
						logging.Field{Key: "layer", Value: "database"},
						logging.Field{Key: "error_code", Value: "DB-005"})
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code": "DB-005",
							"message": "Unable to verify update result",
							"layer": "database",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				if rowsAffected == 0 {
					logger.Error(c.Request.Context(), "Country not found for update", nil,
						logging.Field{Key: "layer", Value: "business_logic"},
						logging.Field{Key: "error_code", Value: "BIZ-002"},
						logging.Field{Key: "country_code", Value: code})
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"code": "BIZ-002",
							"message": "Country not found",
							"layer": "business_logic",
							"timestamp": time.Now().UTC(),
							"country_code": code,
						},
					})
					return
				}
				
				// SUCCESS RESPONSE
				logger.Info(c.Request.Context(), "Country updated successfully",
					logging.Field{Key: "layer", Value: "presentation"},
					logging.Field{Key: "country_code", Value: code},
					logging.Field{Key: "rows_affected", Value: rowsAffected})
				
				c.JSON(http.StatusOK, gin.H{
					"message": "Country updated successfully",
					"updated": true,
					"country_code": code,
					"rows_affected": rowsAffected,
					"timestamp": time.Now().UTC(),
				})
			})
			countries.DELETE("/:code", func(c *gin.Context) {
				// PRESENTATION LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "DELETE /countries/:code started", 
					logging.Field{Key: "country_code", Value: c.Param("code")},
					logging.Field{Key: "layer", Value: "presentation"})
				
				code := c.Param("code")
				if code == "" {
					logger.Error(c.Request.Context(), "Missing country code parameter", nil,
						logging.Field{Key: "layer", Value: "presentation"},
						logging.Field{Key: "error_code", Value: "PRES-001"})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code": "PRES-001",
							"message": "Country code parameter is required",
							"layer": "presentation",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				// DATA ACCESS LAYER ERROR HANDLING
				logger.Info(c.Request.Context(), "Accessing database for country delete",
					logging.Field{Key: "layer", Value: "data_access"})
				
				db, err := container.DBManager.DB.DB()
				if err != nil {
					logger.Error(c.Request.Context(), "Database connection failed", err,
						logging.Field{Key: "layer", Value: "data_access"},
						logging.Field{Key: "error_code", Value: "DAL-001"})
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"error": gin.H{
							"code": "DAL-001",
							"message": "Database connection unavailable",
							"layer": "data_access",
							"timestamp": time.Now().UTC(),
							"details": err.Error(),
						},
					})
					return
				}
				
				// BUSINESS LOGIC LAYER - Check if country is inactive before allowing delete
				logger.Info(c.Request.Context(), "Validating delete business rules",
					logging.Field{Key: "layer", Value: "business_logic"})
				
				var isActive bool
				err = db.QueryRow(`
					SELECT is_active 
					FROM domain_reference_master_geopolitical.countries 
					WHERE country_code = $1 AND tenant_id = 'default-tenant' AND is_deleted = false
				`, code).Scan(&isActive)
				
				if err != nil {
					logger.Error(c.Request.Context(), "Country not found for delete", err,
						logging.Field{Key: "layer", Value: "business_logic"},
						logging.Field{Key: "error_code", Value: "BIZ-003"},
						logging.Field{Key: "country_code", Value: code})
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"code": "BIZ-003",
							"message": "Country not found",
							"layer": "business_logic",
							"timestamp": time.Now().UTC(),
							"country_code": code,
						},
					})
					return
				}
				
				// Business rule: Cannot delete active countries
				if isActive {
					logger.Error(c.Request.Context(), "Cannot delete active country - must be inactive first", nil,
						logging.Field{Key: "layer", Value: "business_logic"},
						logging.Field{Key: "error_code", Value: "BIZ-004"},
						logging.Field{Key: "country_code", Value: code})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code": "BIZ-004",
							"message": "Cannot delete active country. Please deactivate it first by setting is_active to false.",
							"layer": "business_logic",
							"timestamp": time.Now().UTC(),
							"country_code": code,
							"required_action": "Set is_active = false before deletion",
						},
					})
					return
				}
				
				// DATABASE LAYER - SOFT DELETE (only inactive records)
				logger.Info(c.Request.Context(), "Executing database soft delete",
					logging.Field{Key: "layer", Value: "database"},
					logging.Field{Key: "country_code", Value: code})
				
				result, err := db.Exec(`
					UPDATE domain_reference_master_geopolitical.countries 
					SET is_deleted = true,
						deleted_at = NOW(),
						deleted_by = '00000000-0000-0000-0000-000000000001'::uuid,
						updated_at = NOW(),
						updated_by = '00000000-0000-0000-0000-000000000001'::uuid
					WHERE country_code = $1 AND tenant_id = 'default-tenant' AND is_deleted = false AND is_active = false
				`, code)
				
				if err != nil {
					logger.Error(c.Request.Context(), "Database delete execution failed", err,
						logging.Field{Key: "layer", Value: "database"},
						logging.Field{Key: "error_code", Value: "DB-001"},
						logging.Field{Key: "country_code", Value: code})
					
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code": "DB-001",
							"message": "Database delete failed",
							"layer": "database",
							"timestamp": time.Now().UTC(),
							"details": err.Error(),
							"country_code": code,
						},
					})
					return
				}
				
				// Check if any rows were affected
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					logger.Error(c.Request.Context(), "Failed to get rows affected", err,
						logging.Field{Key: "layer", Value: "database"},
						logging.Field{Key: "error_code", Value: "DB-005"})
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code": "DB-005",
							"message": "Unable to verify delete result",
							"layer": "database",
							"timestamp": time.Now().UTC(),
						},
					})
					return
				}
				
				if rowsAffected == 0 {
					logger.Error(c.Request.Context(), "Country not found for delete", nil,
						logging.Field{Key: "layer", Value: "business_logic"},
						logging.Field{Key: "error_code", Value: "BIZ-002"},
						logging.Field{Key: "country_code", Value: code})
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"code": "BIZ-002",
							"message": "Country not found or already deleted",
							"layer": "business_logic",
							"timestamp": time.Now().UTC(),
							"country_code": code,
						},
					})
					return
				}
				
				// SUCCESS RESPONSE
				logger.Info(c.Request.Context(), "Country deleted successfully",
					logging.Field{Key: "layer", Value: "presentation"},
					logging.Field{Key: "country_code", Value: code},
					logging.Field{Key: "rows_affected", Value: rowsAffected})
				
				c.JSON(http.StatusOK, gin.H{
					"message": "Country deleted successfully",
					"deleted": true,
					"country_code": code,
					"rows_affected": rowsAffected,
					"timestamp": time.Now().UTC(),
				})
			})
		}
		
		// Simple endpoints for regions and languages
		v1.GET("/regions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"regions": []gin.H{
					{"region_code": "NA", "region_name": "North America", "region_type": "CONTINENT", "is_active": true},
					{"region_code": "EU", "region_name": "Europe", "region_type": "CONTINENT", "is_active": true},
					{"region_code": "AS", "region_name": "Asia", "region_type": "CONTINENT", "is_active": true},
					{"region_code": "SA", "region_name": "South America", "region_type": "CONTINENT", "is_active": true},
					{"region_code": "AF", "region_name": "Africa", "region_type": "CONTINENT", "is_active": true},
					{"region_code": "OC", "region_name": "Oceania", "region_type": "CONTINENT", "is_active": true},
				},
			})
		})
		
		v1.GET("/languages", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"languages": []gin.H{
					{"language_code": "en", "language_name": "English", "iso3_code": "eng", "native_name": "English", "direction": "LTR", "is_active": true},
					{"language_code": "es", "language_name": "Spanish", "iso3_code": "spa", "native_name": "Español", "direction": "LTR", "is_active": true},
					{"language_code": "fr", "language_name": "French", "iso3_code": "fra", "native_name": "Français", "direction": "LTR", "is_active": true},
					{"language_code": "de", "language_name": "German", "iso3_code": "deu", "native_name": "Deutsch", "direction": "LTR", "is_active": true},
					{"language_code": "zh", "language_name": "Chinese", "iso3_code": "zho", "native_name": "中文", "direction": "LTR", "is_active": true},
					{"language_code": "ja", "language_name": "Japanese", "iso3_code": "jpn", "native_name": "日本語", "direction": "LTR", "is_active": true},
					{"language_code": "ar", "language_name": "Arabic", "iso3_code": "ara", "native_name": "العربية", "direction": "RTL", "is_active": true},
					{"language_code": "hi", "language_name": "Hindi", "iso3_code": "hin", "native_name": "हिन्दी", "direction": "LTR", "is_active": true},
					{"language_code": "pt", "language_name": "Portuguese", "iso3_code": "por", "native_name": "Português", "direction": "LTR", "is_active": true},
					{"language_code": "ru", "language_name": "Russian", "iso3_code": "rus", "native_name": "Русский", "direction": "LTR", "is_active": true},
				},
			})
		})
		
		v1.GET("/timezones", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"timezones": []gin.H{
					{"timezone_code": "UTC", "timezone_name": "Coordinated Universal Time", "utc_offset_hours": 0, "utc_offset_minutes": 0, "supports_dst": false, "is_active": true},
					{"timezone_code": "EST", "timezone_name": "Eastern Standard Time", "utc_offset_hours": -5, "utc_offset_minutes": 0, "supports_dst": true, "dst_offset_hours": 1, "is_active": true},
					{"timezone_code": "PST", "timezone_name": "Pacific Standard Time", "utc_offset_hours": -8, "utc_offset_minutes": 0, "supports_dst": true, "dst_offset_hours": 1, "is_active": true},
					{"timezone_code": "GMT", "timezone_name": "Greenwich Mean Time", "utc_offset_hours": 0, "utc_offset_minutes": 0, "supports_dst": true, "dst_offset_hours": 1, "is_active": true},
					{"timezone_code": "CET", "timezone_name": "Central European Time", "utc_offset_hours": 1, "utc_offset_minutes": 0, "supports_dst": true, "dst_offset_hours": 1, "is_active": true},
					{"timezone_code": "JST", "timezone_name": "Japan Standard Time", "utc_offset_hours": 9, "utc_offset_minutes": 0, "supports_dst": false, "is_active": true},
					{"timezone_code": "CST", "timezone_name": "China Standard Time", "utc_offset_hours": 8, "utc_offset_minutes": 0, "supports_dst": false, "is_active": true},
					{"timezone_code": "IST", "timezone_name": "India Standard Time", "utc_offset_hours": 5, "utc_offset_minutes": 30, "supports_dst": false, "is_active": true},
					{"timezone_code": "AEST", "timezone_name": "Australian Eastern Standard Time", "utc_offset_hours": 10, "utc_offset_minutes": 0, "supports_dst": true, "dst_offset_hours": 1, "is_active": true},
					{"timezone_code": "BRT", "timezone_name": "Brasília Time", "utc_offset_hours": -3, "utc_offset_minutes": 0, "supports_dst": false, "is_active": true},
				},
			})
		})
		
		v1.GET("/subdivisions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"subdivisions": []gin.H{}})
		})
		
		v1.GET("/locales", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"locales": []gin.H{}})
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

	logger.Info(ctx, "Server started successfully",
		logging.Field{Key: "port", Value: cfg.Server.Port})

	// Wait for interrupt signal to gracefully shutdown
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