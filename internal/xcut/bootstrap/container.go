// Package bootstrap implements dependency injection and application initialization
package bootstrap

import (
	"context"
	"fmt"

	"github.com/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/security"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/tracing"
)

type Container struct {
	Config      *config.Config
	Logger      logging.Logger
	Tracer      tracing.Tracer
	AuthService security.AuthService
}

func InitializeContainer(ctx context.Context) (*Container, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger := logging.NewStructuredLogger("geopolitical-service")
	logger.Info(ctx, "Configuration loaded successfully")

	// Initialize tracing
	tracer := tracing.NewTracer("geopolitical-service")

	// Initialize authentication service
	authService := security.NewJWTAuthService(cfg.Auth.JWTSecret)

	logger.Info(ctx, "All services initialized successfully")

	return &Container{
		Config:      cfg,
		Logger:      logger,
		Tracer:      tracer,
		AuthService: authService,
	}, nil
}

func (c *Container) Shutdown(ctx context.Context) error {
	c.Logger.Info(ctx, "Shutting down application")
	return nil
}