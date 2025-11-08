// ============================================================================
// FILE: transaction_manager.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Business Logic Layer - Transaction Management
// PURPOSE: Transaction scope management and consistency
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package transactionmanagement

import (
	"context"
	"time"
	"gorm.io/gorm"
)

// TransactionManager handles transaction boundaries
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
	WithReadOnlyTransaction(ctx context.Context, fn func(context.Context) error) error
}

type gormTransactionManager struct {
	db      *gorm.DB
	timeout time.Duration
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{
		db:      db,
		timeout: 30 * time.Second,
	}
}

// WithTransaction executes function within transaction scope
func (tm *gormTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	// Apply transaction timeout
	ctx, cancel := context.WithTimeout(ctx, tm.timeout)
	defer cancel()
	
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pass transaction context to function
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// WithReadOnlyTransaction for read operations with consistency
func (tm *gormTransactionManager) WithReadOnlyTransaction(ctx context.Context, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, tm.timeout)
	defer cancel()
	
	// Read-only transaction for consistent reads
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set transaction to read-only
		tx.Exec("SET TRANSACTION READ ONLY")
		
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}