// ============================================================================
// FILE: idempotency_service.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Business Logic Layer - Resilience
// PURPOSE: Idempotency and retry patterns for production resilience
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package resilience

import (
	"context"
	"fmt"
	"math"
	"time"
)

// CommandResult represents operation outcome - immutable value object
type CommandResult struct {
	ID      string      `json:"id"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// IdempotencyPolicy defines idempotency rules - pure functions only
type IdempotencyPolicy struct{}

// ShouldExecute checks if operation should run - pure function
func (p IdempotencyPolicy) ShouldExecute(hasResult, isProcessing bool) bool {
	return !hasResult && !isProcessing
}

// IdempotencyStore handles persistence - side effects isolated
type IdempotencyStore interface {
	GetResult(ctx context.Context, key string) (*CommandResult, bool)
	IsProcessing(ctx context.Context, key string) bool
	MarkProcessing(ctx context.Context, key string) error
	StoreResult(ctx context.Context, key string, result *CommandResult, err error) error
}

// CommandExecutor wraps business operations
type CommandExecutor struct {
	fn func() (*CommandResult, error)
}

// Execute runs the wrapped operation - side effect
func (e CommandExecutor) Execute() (*CommandResult, error) {
	return e.fn()
}

// IdempotencyService orchestrates idempotent execution
type IdempotencyService struct {
	policy IdempotencyPolicy
	store  IdempotencyStore
}

func NewIdempotencyService(store IdempotencyStore) *IdempotencyService {
	return &IdempotencyService{
		policy: IdempotencyPolicy{},
		store:  store,
	}
}

// ExecuteIdempotent orchestrates execution - single responsibility
func (s *IdempotencyService) ExecuteIdempotent(ctx context.Context, key string, fn func() (*CommandResult, error)) (*CommandResult, error) {
	result, hasResult := s.store.GetResult(ctx, key)
	isProcessing := s.store.IsProcessing(ctx, key)
	
	if !s.policy.ShouldExecute(hasResult, isProcessing) {
		if hasResult {
			return result, nil
		}
		return nil, fmt.Errorf("request already processing")
	}
	
	if err := s.store.MarkProcessing(ctx, key); err != nil {
		return nil, err
	}
	
	executor := CommandExecutor{fn: fn}
	result, err := executor.Execute()
	
	s.store.StoreResult(ctx, key, result, err)
	return result, err
}

// RetryPolicy configuration for exponential backoff
type RetryPolicy struct {
	MaxAttempts   int           `default:"3"`
	InitialDelay  time.Duration `default:"100ms"`
	MaxDelay      time.Duration `default:"5s"`
	BackoffFactor float64       `default:"2.0"`
}

// RetryService implements retry with exponential backoff
type RetryService struct {
	policy RetryPolicy
}

func NewRetryService(policy RetryPolicy) *RetryService {
	return &RetryService{policy: policy}
}

// ExecuteWithRetry retries operation with exponential backoff
func (rs *RetryService) ExecuteWithRetry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error
	
	for attempt := 1; attempt <= rs.policy.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Don't retry on last attempt
		if attempt < rs.policy.MaxAttempts {
			delay := rs.calculateDelay(attempt)
			
			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	
	return fmt.Errorf("operation failed after %d attempts: %w", rs.policy.MaxAttempts, lastErr)
}

func (rs *RetryService) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(rs.policy.InitialDelay) * math.Pow(rs.policy.BackoffFactor, float64(attempt-1)))
	if delay > rs.policy.MaxDelay {
		delay = rs.policy.MaxDelay
	}
	return delay
}