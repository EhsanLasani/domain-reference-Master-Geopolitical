// PRESENTATION LAYER - Handler
package handlers

import (
	"net/http"
	"reference-master-geopolitical/services"

	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	service *services.CountryService
}

func NewCountryHandler(service *services.CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

func (h *CountryHandler) GetCountries(c *gin.Context) {
	countries, err := h.service.GetActiveCountries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   countries,
		"count":  len(countries),
	})
}