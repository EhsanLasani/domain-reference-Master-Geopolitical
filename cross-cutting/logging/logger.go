// ============================================================================
// FILE: logger.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Cross-Cutting - Logging
// PURPOSE: Structured logging with LASANI audit compliance
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package logging

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type AuditLog struct {
	Action    string    `json:"action"`
	EntityID  string    `json:"entity_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Changes   string    `json:"changes"`
}

func NewLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	
	return &Logger{Logger: log}
}

func (l *Logger) LogAudit(ctx context.Context, audit AuditLog) {
	l.WithFields(logrus.Fields{
		"type":      "audit",
		"action":    audit.Action,
		"entity_id": audit.EntityID,
		"user_id":   audit.UserID,
		"changes":   audit.Changes,
	}).Info("Audit log entry")
}

func (l *Logger) LogError(ctx context.Context, err error, message string) {
	l.WithFields(logrus.Fields{
		"error": err.Error(),
		"type":  "error",
	}).Error(message)
}

func (l *Logger) LogPerformance(ctx context.Context, operation string, duration time.Duration) {
	l.WithFields(logrus.Fields{
		"type":      "performance",
		"operation": operation,
		"duration":  duration.Milliseconds(),
	}).Info("Performance metric")
}