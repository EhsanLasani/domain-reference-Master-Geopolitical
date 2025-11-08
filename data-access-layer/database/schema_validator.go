package database

import (
	"fmt"
)

type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func (d *Database) ValidateTableStructure(tableName string, expectedColumns map[string]string) error {
	var columns []ColumnInfo
	
	query := `
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_schema = 'domain_reference_master_geopolitical' 
		AND table_name = ?
		ORDER BY ordinal_position
	`
	
	if err := d.DB.Raw(query, tableName).Scan(&columns).Error; err != nil {
		return fmt.Errorf("failed to query table structure: %w", err)
	}
	
	fmt.Printf("✅ Table '%s' has %d columns\n", tableName, len(columns))
	
	// Check for missing expected columns
	existingCols := make(map[string]bool)
	for _, col := range columns {
		existingCols[col.ColumnName] = true
	}
	
	for expectedCol := range expectedColumns {
		if !existingCols[expectedCol] {
			fmt.Printf("⚠️  Missing column: %s\n", expectedCol)
		}
	}
	
	return nil
}

func (d *Database) CheckAllTables() error {
	tables := map[string]map[string]string{
		"countries": {
			"country_id":          "uuid",
			"country_code":        "character",
			"country_name":        "character varying",
			"tenant_id":           "character varying",
			"is_active":           "boolean",
			"is_deleted":          "boolean",
			"created_at":          "timestamp with time zone",
			"updated_at":          "timestamp with time zone",
		},
		"regions": {
			"region_id":    "uuid",
			"region_code":  "character varying",
			"region_name":  "character varying",
			"tenant_id":    "character varying",
		},
		"languages": {
			"language_id":   "uuid",
			"language_code": "character",
			"language_name": "character varying",
			"tenant_id":     "character varying",
		},
	}
	
	for tableName, expectedCols := range tables {
		if err := d.ValidateTableStructure(tableName, expectedCols); err != nil {
			fmt.Printf("❌ Table validation failed for %s: %v\n", tableName, err)
		}
	}
	
	return nil
}