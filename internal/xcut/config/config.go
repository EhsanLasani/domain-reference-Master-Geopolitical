// Package config implements externalized configuration management
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Auth     AuthConfig     `yaml:"auth"`
	Tracing  TracingConfig  `yaml:"tracing"`
	Server   ServerConfig   `yaml:"server"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host" env:"DB_HOST" default:"localhost"`
	Port            int           `yaml:"port" env:"DB_PORT" default:"5432"`
	Name            string        `yaml:"name" env:"DB_NAME" default:"geopolitical"`
	User            string        `yaml:"user" env:"DB_USER" default:"postgres"`
	Password        string        `yaml:"password" env:"DB_PASSWORD"`
	SSLMode         string        `yaml:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	MaxConnections  int           `yaml:"max_connections" env:"DB_MAX_CONNECTIONS" default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" default:"localhost"`
	Port     int    `yaml:"port" env:"REDIS_PORT" default:"6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" default:"0"`
}

type AuthConfig struct {
	JWTSecret     string        `yaml:"jwt_secret" env:"JWT_SECRET"`
	TokenExpiry   time.Duration `yaml:"token_expiry" env:"TOKEN_EXPIRY" default:"24h"`
	RefreshExpiry time.Duration `yaml:"refresh_expiry" env:"REFRESH_EXPIRY" default:"168h"`
}

type TracingConfig struct {
	Enabled     bool   `yaml:"enabled" env:"TRACING_ENABLED" default:"true"`
	ServiceName string `yaml:"service_name" env:"SERVICE_NAME" default:"geopolitical-service"`
	Endpoint    string `yaml:"endpoint" env:"TRACING_ENDPOINT"`
}

type ServerConfig struct {
	Host         string        `yaml:"host" env:"SERVER_HOST" default:"0.0.0.0"`
	Port         int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	
	// Load from environment variables
	if err := loadFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}
	
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return config, nil
}

func loadFromEnv(config *Config) error {
	// Database config
	config.Database.Host = getEnvString("DB_HOST", "localhost")
	config.Database.Port = getEnvInt("DB_PORT", 5432)
	config.Database.Name = getEnvString("DB_NAME", "geopolitical")
	config.Database.User = getEnvString("DB_USER", "postgres")
	config.Database.Password = getEnvString("DB_PASSWORD", "")
	config.Database.SSLMode = getEnvString("DB_SSL_MODE", "disable")
	config.Database.MaxConnections = getEnvInt("DB_MAX_CONNECTIONS", 25)
	config.Database.ConnMaxLifetime = getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	
	// Redis config
	config.Redis.Host = getEnvString("REDIS_HOST", "localhost")
	config.Redis.Port = getEnvInt("REDIS_PORT", 6379)
	config.Redis.Password = getEnvString("REDIS_PASSWORD", "")
	config.Redis.DB = getEnvInt("REDIS_DB", 0)
	
	// Auth config
	config.Auth.JWTSecret = getEnvString("JWT_SECRET", "")
	config.Auth.TokenExpiry = getEnvDuration("TOKEN_EXPIRY", 24*time.Hour)
	config.Auth.RefreshExpiry = getEnvDuration("REFRESH_EXPIRY", 168*time.Hour)
	
	// Tracing config
	config.Tracing.Enabled = getEnvBool("TRACING_ENABLED", true)
	config.Tracing.ServiceName = getEnvString("SERVICE_NAME", "geopolitical-service")
	config.Tracing.Endpoint = getEnvString("TRACING_ENDPOINT", "")
	
	// Server config
	config.Server.Host = getEnvString("SERVER_HOST", "0.0.0.0")
	config.Server.Port = getEnvInt("SERVER_PORT", 8080)
	config.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second)
	config.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)
	
	return nil
}

func validateConfig(config *Config) error {
	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

// Helper functions
func getEnvString(key, defaultValue string) string {
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

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}