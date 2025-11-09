package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/bootstrap"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	v1 "github.com/EhsanLasani/domain-reference-Master-Geopolitical/presentation-layer/rest-api/v1"
)

func main() {
	ctx := context.Background()
	
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
	
	logger := logging.NewStructuredLogger("debug")
	
	container, err := bootstrap.NewContainer(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
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
	
	router.Use(func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default-tenant"
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "reference-master-geopolitical",
			"version": "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	router.Static("/web-ui", "./presentation-layer/web-ui")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web-ui/business-central.html")
	})

	v1Group := router.Group("/api/v1")
	{
		countriesHandler := v1.NewCountriesHandler(container.CountryAppService, container.Logger)
		
		countries := v1Group.Group("/countries")
		{
			countries.GET("", countriesHandler.GetAllCountries)
			countries.POST("", countriesHandler.CreateCountry)
			countries.GET("/:code", countriesHandler.GetCountryByCode)
		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info(ctx, "🚀 Simple API Server started", logging.Field{Key: "port", Value: cfg.Server.Port})

	select {}
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