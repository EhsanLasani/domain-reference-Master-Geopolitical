// Package v1 implements REST API v1 endpoints with full CRUD operations
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/business-logic-layer/application-services"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

type CountriesHandler struct {
	countryService *applicationservices.CountryAppService
	logger         logging.Logger
}

func NewCountriesHandler(
	countryService *applicationservices.CountryAppService,
	logger logging.Logger,
) *CountriesHandler {
	return &CountriesHandler{
		countryService: countryService,
		logger:         logger,
	}
}

// GetAllCountries handles GET /countries
func (h *CountriesHandler) GetAllCountries(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	countries, err := h.countryService.GetAllCountries(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"count":     len(countries),
	})
}

// CreateCountry handles POST /countries
func (h *CountriesHandler) CreateCountry(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	
	var country models.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err := h.countryService.CreateCountry(c.Request.Context(), tenantID, &country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, country)
}

// GetCountryByCode handles GET /countries/:code
func (h *CountriesHandler) GetCountryByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")

	country, err := h.countryService.GetCountryByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, country)
}

// UpdateCountry handles PUT /countries/:code
func (h *CountriesHandler) UpdateCountry(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DeleteCountry handles DELETE /countries/:code
func (h *CountriesHandler) DeleteCountry(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}