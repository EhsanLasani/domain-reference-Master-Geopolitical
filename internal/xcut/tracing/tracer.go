package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Tracer interface defines tracing operations
type Tracer interface {
	StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, oteltrace.Span)
	StartSQLSpan(ctx context.Context, sqlKey, operation string, tenantID string) (context.Context, oteltrace.Span)
}

// OtelTracer implements Tracer interface
type OtelTracer struct {
	tracer oteltrace.Tracer
}

func NewTracer(serviceName string) (Tracer, error) {
	exporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	tracer := otel.Tracer(serviceName)

	return &OtelTracer{tracer: tracer}, nil
}

func (t *OtelTracer) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, oteltrace.Span) {
	return t.tracer.Start(ctx, name, oteltrace.WithAttributes(attrs...))
}

func (t *OtelTracer) StartSQLSpan(ctx context.Context, sqlKey, operation string, tenantID string) (context.Context, oteltrace.Span) {
	return t.tracer.Start(ctx, "sql."+operation, oteltrace.WithAttributes(
		attribute.String("sql_key", sqlKey),
		attribute.String("db.operation", operation),
		attribute.String("tenant_id", tenantID),
	))
}