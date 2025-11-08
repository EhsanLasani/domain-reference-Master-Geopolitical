// ============================================================================
// FILE: query_monitor.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - Performance Monitoring
// PURPOSE: Query performance monitoring and N+1 detection
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package performance

import (
	"context"
	"log"
	"time"
)

// Performance budget constants
const (
	MaxSingleEntityLatency = 10 * time.Millisecond
	MaxCollectionLatency   = 50 * time.Millisecond
	MaxComplexQueryLatency = 100 * time.Millisecond
	MaxBulkOperationLatency = 500 * time.Millisecond
)

// QueryMonitor tracks query performance
type QueryMonitor struct {
	queryCount int
	startTime  time.Time
}

func NewQueryMonitor() *QueryMonitor {
	return &QueryMonitor{
		startTime: time.Now(),
	}
}

// TrackQuery monitors individual query performance
func (qm *QueryMonitor) TrackQuery(operation string, budget time.Duration, fn func() error) error {
	start := time.Now()
	qm.queryCount++
	
	err := fn()
	
	duration := time.Since(start)
	if duration > budget {
		log.Printf("PERFORMANCE WARNING: %s exceeded budget: %v > %v", 
			operation, duration, budget)
	}
	
	return err
}

// DetectN1Queries warns if too many queries executed
func (qm *QueryMonitor) DetectN1Queries(expectedQueries int) {
	if qm.queryCount > expectedQueries {
		log.Printf("N+1 QUERY DETECTED: Expected %d queries, executed %d", 
			expectedQueries, qm.queryCount)
	}
}

// WithPerformanceMonitoring wraps repository methods with monitoring
func WithPerformanceMonitoring(ctx context.Context, operation string, budget time.Duration, fn func() error) error {
	monitor := NewQueryMonitor()
	return monitor.TrackQuery(operation, budget, fn)
}