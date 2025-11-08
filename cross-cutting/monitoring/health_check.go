// ============================================================================
// FILE: health_check.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Cross-Cutting - Monitoring
// PURPOSE: Health check endpoints and monitoring
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package monitoring

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

type HealthChecker struct {
	dbPing func() error
}

func NewHealthChecker(dbPing func() error) *HealthChecker {
	return &HealthChecker{
		dbPing: dbPing,
	}
}

func (h *HealthChecker) HealthCheck(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	// Check database
	if err := h.dbPing(); err != nil {
		status.Status = "unhealthy"
		status.Services["database"] = "down"
	} else {
		status.Services["database"] = "up"
	}

	httpStatus := http.StatusOK
	if status.Status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, status)
}

func (h *HealthChecker) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"timestamp": time.Now(),
	})
}