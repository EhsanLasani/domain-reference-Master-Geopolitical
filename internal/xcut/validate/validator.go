// Package validate implements schema validation and business rules
package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(data interface{}) error
	ValidateStruct(s interface{}) error
}

type StructValidator struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	v := validator.New()
	
	// Register custom validations
	v.RegisterValidation("country_code", validateCountryCode)
	v.RegisterValidation("tenant_id", validateTenantID)
	
	return &StructValidator{validator: v}
}

func (sv *StructValidator) Validate(data interface{}) error {
	return sv.validator.Struct(data)
}

func (sv *StructValidator) ValidateStruct(s interface{}) error {
	err := sv.validator.Struct(s)
	if err != nil {
		return formatValidationError(err)
	}
	return nil
}

func validateCountryCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	return len(code) == 2 && strings.ToUpper(code) == code
}

func validateTenantID(fl validator.FieldLevel) bool {
	tenantID := fl.Field().String()
	return len(tenantID) >= 8 && len(tenantID) <= 36
}

func formatValidationError(err error) error {
	var errors []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
			case "country_code":
				errors = append(errors, fmt.Sprintf("%s must be a valid 2-letter country code", e.Field()))
			case "tenant_id":
				errors = append(errors, fmt.Sprintf("%s must be a valid tenant ID", e.Field()))
			default:
				errors = append(errors, fmt.Sprintf("%s validation failed: %s", e.Field(), e.Tag()))
			}
		}
	}
	
	return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
}

// Business rule validators
type BusinessRuleValidator struct {
	validator Validator
}

func NewBusinessRuleValidator() *BusinessRuleValidator {
	return &BusinessRuleValidator{
		validator: NewValidator(),
	}
}

func (brv *BusinessRuleValidator) ValidateCountryCreation(country interface{}) error {
	// Basic struct validation
	if err := brv.validator.ValidateStruct(country); err != nil {
		return err
	}
	
	// Business rule validation
	v := reflect.ValueOf(country)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	// Check if country code is not empty
	codeField := v.FieldByName("CountryCode")
	if codeField.IsValid() && codeField.String() == "" {
		return fmt.Errorf("country code cannot be empty")
	}
	
	return nil
}