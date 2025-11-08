// ============================================================================
// FILE: country_handler_proper.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - Web/Mobile Applications
// PURPOSE: REST API handlers for country management
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	countryService CountryService
}

type CountryService interface {
	GetCountries(filter CountryFilter) ([]Country, error)
	GetCountryByCode(code string) (*Country, error)
}

type Country struct {
	ID           string `json:"country_id"`
	Code         string `json:"country_code"`
	Name         string `json:"country_name"`
	ISO3Code     string `json:"iso3_code"`
	OfficialName string `json:"official_name"`
	IsActive     bool   `json:"is_active"`
}

type CountryFilter struct {
	IsActive      *bool  `form:"is_active"`
	ContinentCode string `form:"continent_code"`
	Limit         int    `form:"limit"`
	Offset        int    `form:"offset"`
}

func NewCountryHandler(service CountryService) *CountryHandler {
	return &CountryHandler{
		countryService: service,
	}
}

func (h *CountryHandler) GetCountries(c *gin.Context) {
	var filter CountryFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	countries, err := h.countryService.GetCountries(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     countries,
		"total":    len(countries),
		"has_more": len(countries) == filter.Limit,
	})
}

func (h *CountryHandler) GetCountryByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Country code is required"})
		return
	}

	country, err := h.countryService.GetCountryByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, country)
}

func (h *CountryHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/countries", h.GetCountries)
		v1.GET("/countries/:code", h.GetCountryByCode)
	}
}