package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics implements guideline 02 - Metrics collection
type Metrics struct {
	meter           metric.Meter
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	errorCounter    metric.Int64Counter
	dbQueryCounter  metric.Int64Counter
	dbQueryDuration metric.Float64Histogram
}

func NewMetrics(serviceName string) (*Metrics, error) {
	meter := otel.Meter(serviceName)

	requestCounter, err := meter.Int64Counter("http_requests_total")
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram("http_request_duration_seconds")
	if err != nil {
		return nil, err
	}

	errorCounter, err := meter.Int64Counter("errors_total")
	if err != nil {
		return nil, err
	}

	dbQueryCounter, err := meter.Int64Counter("db_queries_total")
	if err != nil {
		return nil, err
	}

	dbQueryDuration, err := meter.Float64Histogram("db_query_duration_seconds")
	if err != nil {
		return nil, err
	}

	return &Metrics{
		meter:           meter,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		errorCounter:    errorCounter,
		dbQueryCounter:  dbQueryCounter,
		dbQueryDuration: dbQueryDuration,
	}, nil
}

func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path, status string, duration time.Duration, tenantID string) {
	m.requestCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.String("status", status),
		attribute.String("tenant_id", tenantID),
	))

	m.requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.String("tenant_id", tenantID),
	))
}

func (m *Metrics) RecordError(ctx context.Context, errorCode, layer string, tenantID string) {
	m.errorCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("error_code", errorCode),
		attribute.String("layer", layer),
		attribute.String("tenant_id", tenantID),
	))
}

func (m *Metrics) RecordDBQuery(ctx context.Context, sqlKey, operation string, duration time.Duration, tenantID string) {
	m.dbQueryCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("sql_key", sqlKey),
		attribute.String("operation", operation),
		attribute.String("tenant_id", tenantID),
	))

	m.dbQueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		attribute.String("sql_key", sqlKey),
		attribute.String("operation", operation),
		attribute.String("tenant_id", tenantID),
	))
}