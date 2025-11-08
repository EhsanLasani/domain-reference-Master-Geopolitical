// ============================================================================
// FILE: bootstrap.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer
// PURPOSE: Database validation and setup - side effects isolated
// ============================================================================

package dal

import (
	"database/sql"
	"fmt"
)

// ValidateDatabase checks database connectivity and schema
// Side effect - database validation
func ValidateDatabase(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	return nil
}

// VerifySchema checks required tables exist
// Side effect - schema verification
func VerifySchema(db *sql.DB, schema string) error {
	query := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1`
	var count int
	err := db.QueryRow(query, schema).Scan(&count)
	if err != nil {
		return fmt.Errorf("schema verification failed: %w", err)
	}
	if count < 6 {
		return fmt.Errorf("incomplete schema: expected 6 tables, found %d", count)
	}
	return nil
}

// Bootstrap initializes data access layer
// Side effect - performs all validation and setup
func Bootstrap(db *sql.DB) error {
	if err := ValidateDatabase(db); err != nil {
		return err
	}
	return VerifySchema(db, "domain_reference_master_geopolitical")
}