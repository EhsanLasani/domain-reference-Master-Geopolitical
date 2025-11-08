// Package logging implements structured logging with correlation IDs and audit trails
package logging

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

type StructuredLogger struct {
	logger *slog.Logger
}

type LogEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         string    `json:"level"`
	Message       string    `json:"message"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	TenantID      string    `json:"tenant_id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	SQLKey        string    `json:"sql_key,omitempty"`
	UseCase       string    `json:"use_case,omitempty"`
	Component     string    `json:"component"`
	Error         string    `json:"error,omitempty"`
}

func NewStructuredLogger(component string) Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	
	return &StructuredLogger{
		logger: slog.New(handler).With("component", component),
	}
}

func (l *StructuredLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelInfo, msg, nil, fields...)
}

func (l *StructuredLogger) Error(ctx context.Context, msg string, err error, fields ...Field) {
	l.log(ctx, slog.LevelError, msg, err, fields...)
}

func (l *StructuredLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelWarn, msg, nil, fields...)
}

func (l *StructuredLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelDebug, msg, nil, fields...)
}

func (l *StructuredLogger) log(ctx context.Context, level slog.Level, msg string, err error, fields ...Field) {
	attrs := []slog.Attr{
		slog.Time("timestamp", time.Now()),
		slog.String("message", msg),
	}
	
	// Extract correlation context
	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		attrs = append(attrs, slog.String("correlation_id", correlationID))
	}
	if tenantID := GetTenantID(ctx); tenantID != "" {
		attrs = append(attrs, slog.String("tenant_id", tenantID))
	}
	if userID := GetUserID(ctx); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	
	// Add error if present
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	
	// Add custom fields
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	
	l.logger.LogAttrs(ctx, level, msg, attrs...)
}

// Context helpers
type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	tenantIDKey      contextKey = "tenant_id"
	userIDKey        contextKey = "user_id"
)

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tenantIDKey, id)
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

func GetTenantID(ctx context.Context) string {
	if id, ok := ctx.Value(tenantIDKey).(string); ok {
		return id
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}