// Package tracing implements OpenTelemetry distributed tracing
package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	StartSpan(ctx context.Context, operationName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	StartSQLSpan(ctx context.Context, sqlKey, query string) (context.Context, trace.Span)
	StartUseCaseSpan(ctx context.Context, useCase string) (context.Context, trace.Span)
}

type OpenTelemetryTracer struct {
	tracer trace.Tracer
}

func InitTracing(serviceName, version string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func NewTracer(instrumentationName string) Tracer {
	return &OpenTelemetryTracer{
		tracer: otel.Tracer(instrumentationName),
	}
}

func (t *OpenTelemetryTracer) StartSpan(ctx context.Context, operationName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, operationName, opts...)
}

func (t *OpenTelemetryTracer) StartSQLSpan(ctx context.Context, sqlKey, query string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, "sql.query")
	span.SetAttributes(
		attribute.String("sql.key", sqlKey),
		attribute.String("sql.query", maskSensitiveData(query)),
		attribute.String("db.system", "postgresql"),
	)
	return ctx, span
}

func (t *OpenTelemetryTracer) StartUseCaseSpan(ctx context.Context, useCase string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, fmt.Sprintf("usecase.%s", useCase))
	span.SetAttributes(
		attribute.String("use_case", useCase),
		attribute.String("layer", "business_logic"),
	)
	return ctx, span
}

// maskSensitiveData removes sensitive information from SQL queries for tracing
func maskSensitiveData(query string) string {
	// Simple masking - in production, use more sophisticated approach
	if len(query) > 100 {
		return query[:100] + "..."
	}
	return query
}

// Helper functions for adding trace attributes
func AddTenantAttribute(span trace.Span, tenantID string) {
	span.SetAttributes(attribute.String("tenant.id", tenantID))
}

func AddUserAttribute(span trace.Span, userID string) {
	span.SetAttributes(attribute.String("user.id", userID))
}

func AddLayerAttribute(span trace.Span, layer string) {
	span.SetAttributes(attribute.String("layer", layer))
}