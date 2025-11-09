package v1

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/business-logic-layer/application-services"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/errors"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

type ValidatedCountriesHandler struct {
	countryService *applicationservices.CountryAppService
	logger         logging.Logger
}

func NewValidatedCountriesHandler(
	countryService *applicationservices.CountryAppService,
	logger logging.Logger,
) *ValidatedCountriesHandler {
	return &ValidatedCountriesHandler{
		countryService: countryService,
		logger:         logger,
	}
}

func (h *ValidatedCountriesHandler) GetAllCountries(c *gin.Context) {
	tenantID := h.validateTenantID(c)
	if tenantID == "" {
		return
	}
	
	countries, err := h.countryService.GetAllCountries(c.Request.Context(), tenantID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"count":     len(countries),
		"tenant_id": tenantID,
	})
}

func (h *ValidatedCountriesHandler) CreateCountry(c *gin.Context) {
	tenantID := h.validateTenantID(c)
	if tenantID == "" {
		return
	}
	
	var country models.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		h.respondWithError(c, errors.NewPresentationError("INVALID_JSON", "Invalid JSON format", err))
		return
	}

	// Presentation layer validation
	if err := h.validateCountryInput(&country); err != nil {
		h.respondWithError(c, err)
		return
	}

	err := h.countryService.CreateCountry(c.Request.Context(), tenantID, &country)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Country created successfully",
		"country": country,
	})
}

func (h *ValidatedCountriesHandler) GetCountryByCode(c *gin.Context) {
	tenantID := h.validateTenantID(c)
	if tenantID == "" {
		return
	}
	
	code := c.Param("code")
	if err := h.validateCountryCode(code); err != nil {
		h.respondWithError(c, err)
		return
	}

	country, err := h.countryService.GetCountryByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if country == nil {
		h.respondWithError(c, errors.NewPresentationError("NOT_FOUND", "Country not found", nil))
		return
	}

	c.JSON(http.StatusOK, country)
}

func (h *ValidatedCountriesHandler) UpdateCountry(c *gin.Context) {
	tenantID := h.validateTenantID(c)
	if tenantID == "" {
		return
	}
	
	code := c.Param("code")
	if err := h.validateCountryCode(code); err != nil {
		h.respondWithError(c, err)
		return
	}

	var country models.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		h.respondWithError(c, errors.NewPresentationError("INVALID_JSON", "Invalid JSON format", err))
		return
	}

	// Ensure code consistency
	country.CountryCode = code

	// Presentation layer validation
	if err := h.validateCountryInput(&country); err != nil {
		h.respondWithError(c, err)
		return
	}

	// Note: Update method needs to be implemented in CountryAppService
	c.JSON(http.StatusOK, gin.H{
		"message": "Country update validation passed",
		"country": country,
	})
}

func (h *ValidatedCountriesHandler) DeleteCountry(c *gin.Context) {
	tenantID := h.validateTenantID(c)
	if tenantID == "" {
		return
	}
	
	code := c.Param("code")
	if err := h.validateCountryCode(code); err != nil {
		h.respondWithError(c, err)
		return
	}

	// Note: Delete method needs to be implemented in CountryAppService
	c.JSON(http.StatusOK, gin.H{
		"message": "Country delete validation passed",
		"code":    code,
	})
}

// Validation methods
func (h *ValidatedCountriesHandler) validateTenantID(c *gin.Context) string {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		h.respondWithError(c, errors.NewPresentationError("MISSING_TENANT", "Tenant ID is required", nil))
		return ""
	}
	return tenantID
}

func (h *ValidatedCountriesHandler) validateCountryCode(code string) *errors.LayerError {
	if code == "" {
		return errors.NewValidationError("country_code", "cannot be empty")
	}
	if len(code) != 2 {
		return errors.NewValidationError("country_code", "must be exactly 2 characters")
	}
	// Check if it's alphabetic
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return errors.NewValidationError("country_code", "must contain only uppercase letters")
		}
	}
	return nil
}

func (h *ValidatedCountriesHandler) validateCountryInput(country *models.Country) *errors.LayerError {
	if country.CountryCode == "" {
		return errors.NewValidationError("country_code", "is required")
	}
	if country.CountryName == "" {
		return errors.NewValidationError("country_name", "is required")
	}
	if len(country.CountryName) > 100 {
		return errors.NewValidationError("country_name", "cannot exceed 100 characters")
	}
	
	// Validate ISO3 code if provided
	if country.ISO3Code != nil && *country.ISO3Code != "" {
		if len(*country.ISO3Code) != 3 {
			return errors.NewValidationError("iso3_code", "must be exactly 3 characters")
		}
	}
	
	// Validate continent code if provided
	if country.ContinentCode != nil && *country.ContinentCode != "" {
		validContinents := []string{"AF", "AS", "EU", "NA", "SA", "OC", "AN"}
		valid := false
		for _, continent := range validContinents {
			if *country.ContinentCode == continent {
				valid = true
				break
			}
		}
		if !valid {
			return errors.NewValidationError("continent_code", "must be one of: AF, AS, EU, NA, SA, OC, AN")
		}
	}
	
	return nil
}

// Error handling
func (h *ValidatedCountriesHandler) handleError(c *gin.Context, err error) {
	if layerErr, ok := err.(*errors.LayerError); ok {
		h.respondWithError(c, layerErr)
		return
	}
	
	// Convert unknown errors to internal server error
	h.respondWithError(c, errors.NewPresentationError("INTERNAL_ERROR", "Internal server error", err))
}

func (h *ValidatedCountriesHandler) respondWithError(c *gin.Context, err *errors.LayerError) {
	h.logger.Error(c.Request.Context(), "Request failed", err,
		logging.Field{Key: "layer", Value: err.Layer},
		logging.Field{Key: "code", Value: err.Code},
		logging.Field{Key: "path", Value: c.Request.URL.Path})
	
	response := errors.NewErrorResponse(err)
	c.JSON(err.HTTPStatus(), response)
}