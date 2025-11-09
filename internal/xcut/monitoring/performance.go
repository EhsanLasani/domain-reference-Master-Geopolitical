package monitoring

import (
	"context"
	"time"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/metrics"
)

// PerformanceMonitor implements guideline 14 - Monitoring & Alerting
type PerformanceMonitor struct {
	logger  *logging.StructuredLogger
	metrics *metrics.Metrics
}

func NewPerformanceMonitor(logger *logging.StructuredLogger, metrics *metrics.Metrics) *PerformanceMonitor {
	return &PerformanceMonitor{
		logger:  logger,
		metrics: metrics,
	}
}

// MonitorQuery tracks database query performance
func (pm *PerformanceMonitor) MonitorQuery(ctx context.Context, sqlKey string, tenantID string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	// Log slow queries (>100ms)
	if duration > 100*time.Millisecond {
		pm.logger.Warn(ctx, "Slow query detected",
			logging.Field{Key: "sql_key", Value: sqlKey},
			logging.Field{Key: "tenant_id", Value: tenantID},
			logging.Field{Key: "duration_ms", Value: duration.Milliseconds()})
	}

	// Record metrics
	operation := "SELECT"
	if err != nil {
		operation = "ERROR"
	}
	pm.metrics.RecordDBQuery(ctx, sqlKey, operation, duration, tenantID)

	return err
}

// MonitorHTTP tracks HTTP request performance
func (pm *PerformanceMonitor) MonitorHTTP(ctx context.Context, method, path string, tenantID string, fn func() (int, error)) (int, error) {
	start := time.Now()
	status, err := fn()
	duration := time.Since(start)

	statusStr := "200"
	if err != nil {
		statusStr = "500"
	}

	pm.metrics.RecordHTTPRequest(ctx, method, path, statusStr, duration, tenantID)

	// Log slow requests (>1s)
	if duration > time.Second {
		pm.logger.Warn(ctx, "Slow HTTP request",
			logging.Field{Key: "method", Value: method},
			logging.Field{Key: "path", Value: path},
			logging.Field{Key: "tenant_id", Value: tenantID},
			logging.Field{Key: "duration_ms", Value: duration.Milliseconds()})
	}

	return status, err
}