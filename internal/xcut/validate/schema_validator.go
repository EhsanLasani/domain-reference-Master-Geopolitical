package validate

import (
	"fmt"
	"regexp"
)

type SchemaValidator struct {
	schemas map[string]Schema
}

type Schema struct {
	Type       string            `json:"type"`
	Properties map[string]Field  `json:"properties"`
	Required   []string          `json:"required"`
}

type Field struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern,omitempty"`
	MinLen  int    `json:"minLength,omitempty"`
	MaxLen  int    `json:"maxLength,omitempty"`
}

func NewSchemaValidator() *SchemaValidator {
	sv := &SchemaValidator{
		schemas: make(map[string]Schema),
	}
	sv.loadDefaultSchemas()
	return sv
}

func (sv *SchemaValidator) Validate(data interface{}, schemaName string) error {
	schema, exists := sv.schemas[schemaName]
	if !exists {
		return fmt.Errorf("schema not found: %s", schemaName)
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("data must be object")
	}

	// Check required fields
	for _, field := range schema.Required {
		if _, exists := dataMap[field]; !exists {
			return fmt.Errorf("required field missing: %s", field)
		}
	}

	// Validate fields
	for fieldName, fieldSchema := range schema.Properties {
		if value, exists := dataMap[fieldName]; exists {
			if err := sv.validateField(value, fieldSchema); err != nil {
				return fmt.Errorf("field %s: %w", fieldName, err)
			}
		}
	}

	return nil
}

func (sv *SchemaValidator) validateField(value interface{}, field Field) error {
	str, ok := value.(string)
	if !ok && field.Type == "string" {
		return fmt.Errorf("expected string")
	}

	if field.Pattern != "" {
		matched, _ := regexp.MatchString(field.Pattern, str)
		if !matched {
			return fmt.Errorf("pattern mismatch")
		}
	}

	if field.MinLen > 0 && len(str) < field.MinLen {
		return fmt.Errorf("too short")
	}

	if field.MaxLen > 0 && len(str) > field.MaxLen {
		return fmt.Errorf("too long")
	}

	return nil
}

func (sv *SchemaValidator) loadDefaultSchemas() {
	sv.schemas["country"] = Schema{
		Type: "object",
		Properties: map[string]Field{
			"country_code": {Type: "string", Pattern: "^[A-Z]{2}$", MinLen: 2, MaxLen: 2},
			"country_name": {Type: "string", MinLen: 1, MaxLen: 100},
			"iso3_code":    {Type: "string", Pattern: "^[A-Z]{3}$", MinLen: 3, MaxLen: 3},
		},
		Required: []string{"country_code", "country_name"},
	}
}