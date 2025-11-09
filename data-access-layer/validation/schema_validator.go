package validation

import (
	"fmt"
	"strings"
	"reflect"
	"gorm.io/gorm"
)

type ColumnInfo struct {
	Name         string
	DataType     string
	IsNullable   bool
	MaxLength    int
	IsPrimaryKey bool
}

type SchemaValidator struct {
	db     *gorm.DB
	cache  map[string]map[string]ColumnInfo
}

func NewSchemaValidator(db *gorm.DB) *SchemaValidator {
	return &SchemaValidator{
		db:    db,
		cache: make(map[string]map[string]ColumnInfo),
	}
}

func (v *SchemaValidator) GetTableSchema(tableName string) (map[string]ColumnInfo, error) {
	if schema, exists := v.cache[tableName]; exists {
		return schema, nil
	}

	var columns []struct {
		ColumnName    string `gorm:"column:column_name"`
		DataType      string `gorm:"column:data_type"`
		IsNullable    string `gorm:"column:is_nullable"`
		MaxLength     *int   `gorm:"column:character_maximum_length"`
		IsPrimaryKey  bool   `gorm:"column:is_primary_key"`
	}

	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable,
			character_maximum_length,
			EXISTS(
				SELECT 1 FROM information_schema.key_column_usage k
				WHERE k.table_schema = c.table_schema 
				AND k.table_name = c.table_name 
				AND k.column_name = c.column_name
				AND k.constraint_name LIKE '%_pkey'
			) as is_primary_key
		FROM information_schema.columns c
		WHERE table_schema = 'domain_reference_master_geopolitical' 
		AND table_name = ?`

	err := v.db.Raw(query, tableName).Scan(&columns).Error
	if err != nil {
		return nil, fmt.Errorf("schema introspection failed: %w", err)
	}

	schema := make(map[string]ColumnInfo)
	for _, col := range columns {
		maxLen := 0
		if col.MaxLength != nil {
			maxLen = *col.MaxLength
		}
		
		schema[col.ColumnName] = ColumnInfo{
			Name:         col.ColumnName,
			DataType:     col.DataType,
			IsNullable:   col.IsNullable == "YES",
			MaxLength:    maxLen,
			IsPrimaryKey: col.IsPrimaryKey,
		}
	}

	v.cache[tableName] = schema
	return schema, nil
}

func (v *SchemaValidator) ValidateStruct(tableName string, data interface{}) error {
	schema, err := v.GetTableSchema(tableName)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		// Get database column name from gorm tag
		gormTag := fieldType.Tag.Get("gorm")
		columnName := extractColumnName(gormTag, fieldType.Name)
		
		if columnName == "" {
			continue
		}

		colInfo, exists := schema[columnName]
		if !exists {
			return fmt.Errorf("field %s (column: %s) does not exist in database", fieldType.Name, columnName)
		}

		// Validate field value
		if err := v.validateFieldValue(colInfo, field.Interface()); err != nil {
			return fmt.Errorf("field %s validation failed: %w", fieldType.Name, err)
		}
	}

	return nil
}

func (v *SchemaValidator) validateFieldValue(col ColumnInfo, value interface{}) error {
	// Handle nil values
	if value == nil {
		if !col.IsNullable {
			return fmt.Errorf("field %s cannot be null", col.Name)
		}
		return nil
	}

	// Handle pointer types
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			if !col.IsNullable {
				return fmt.Errorf("field %s cannot be null", col.Name)
			}
			return nil
		}
		val = val.Elem()
		value = val.Interface()
	}

	// Type validation
	switch strings.ToLower(col.DataType) {
	case "uuid":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field %s must be UUID string", col.Name)
		}
	case "varchar", "text", "char":
		if str, ok := value.(string); ok {
			if col.MaxLength > 0 && len(str) > col.MaxLength {
				return fmt.Errorf("field %s exceeds max length %d (got %d)", col.Name, col.MaxLength, len(str))
			}
		} else {
			return fmt.Errorf("field %s must be string", col.Name)
		}
	case "integer", "smallint":
		switch value.(type) {
		case int, int16, int32, int64:
			// Valid integer types
		default:
			return fmt.Errorf("field %s must be integer", col.Name)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field %s must be boolean", col.Name)
		}
	}

	return nil
}

func extractColumnName(gormTag, fieldName string) string {
	if gormTag == "" {
		return strings.ToLower(fieldName)
	}
	
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	
	return strings.ToLower(fieldName)
}