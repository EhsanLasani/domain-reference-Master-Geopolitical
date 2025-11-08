# Cross-Cutting Concerns

## Overview
Cross-cutting concerns that span across all architectural layers in the Reference Master Geopolitical domain.

## Components

### Security
- **auth_middleware.go**: Authentication and authorization middleware
- Rate limiting and token validation
- Security headers and CORS handling

### Monitoring  
- **health_check.go**: Health check endpoints
- Database connectivity monitoring
- Service readiness checks

### Logging
- **logger.go**: Structured logging with LASANI audit compliance
- Performance metrics logging
- Error tracking and audit trails

## Usage

```go
// Security
authMiddleware := security.NewAuthMiddleware("secret-key")
router.Use(authMiddleware.RequireAuth())
router.Use(authMiddleware.RateLimit(100))

// Monitoring
healthChecker := monitoring.NewHealthChecker(db.Ping)
router.GET("/health", healthChecker.HealthCheck)
router.GET("/ready", healthChecker.ReadinessCheck)

// Logging
logger := logging.NewLogger()
logger.LogAudit(ctx, logging.AuditLog{
    Action: "CREATE_COUNTRY",
    EntityID: "country-123",
    UserID: "user-456",
})
```