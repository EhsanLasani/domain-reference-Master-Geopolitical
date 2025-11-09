package v1

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/repositories-daos"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

// RegionsHandler handles region endpoints
type RegionsHandler struct {
	repo *repositories.RegionRepository
}

func NewRegionsHandler(repo *repositories.RegionRepository) *RegionsHandler {
	return &RegionsHandler{repo: repo}
}

func (h *RegionsHandler) GetAll(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	regions, err := h.repo.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"regions": regions, "count": len(regions)})
}

func (h *RegionsHandler) GetByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	region, err := h.repo.GetByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if region == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}
	c.JSON(http.StatusOK, region)
}

func (h *RegionsHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), tenantID, &region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, region)
}

func (h *RegionsHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	region.RegionCode = code
	if err := h.repo.Update(c.Request.Context(), tenantID, &region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, region)
}

func (h *RegionsHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	if err := h.repo.Delete(c.Request.Context(), tenantID, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "region deleted"})
}

// LanguagesHandler handles language endpoints
type LanguagesHandler struct {
	repo *repositories.LanguageRepository
}

func NewLanguagesHandler(repo *repositories.LanguageRepository) *LanguagesHandler {
	return &LanguagesHandler{repo: repo}
}

func (h *LanguagesHandler) GetAll(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	languages, err := h.repo.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"languages": languages, "count": len(languages)})
}

func (h *LanguagesHandler) GetByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	language, err := h.repo.GetByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if language == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "language not found"})
		return
	}
	c.JSON(http.StatusOK, language)
}

func (h *LanguagesHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var language models.Language
	if err := c.ShouldBindJSON(&language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), tenantID, &language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, language)
}

func (h *LanguagesHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	var language models.Language
	if err := c.ShouldBindJSON(&language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	language.LanguageCode = code
	if err := h.repo.Update(c.Request.Context(), tenantID, &language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, language)
}

func (h *LanguagesHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	if err := h.repo.Delete(c.Request.Context(), tenantID, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "language deleted"})
}

// TimezonesHandler handles timezone endpoints
type TimezonesHandler struct {
	repo *repositories.TimezoneRepository
}

func NewTimezonesHandler(repo *repositories.TimezoneRepository) *TimezonesHandler {
	return &TimezonesHandler{repo: repo}
}

func (h *TimezonesHandler) GetAll(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	timezones, err := h.repo.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"timezones": timezones, "count": len(timezones)})
}

func (h *TimezonesHandler) GetByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	timezone, err := h.repo.GetByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if timezone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "timezone not found"})
		return
	}
	c.JSON(http.StatusOK, timezone)
}

func (h *TimezonesHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var timezone models.Timezone
	if err := c.ShouldBindJSON(&timezone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), tenantID, &timezone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, timezone)
}

func (h *TimezonesHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	var timezone models.Timezone
	if err := c.ShouldBindJSON(&timezone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	timezone.TimezoneCode = code
	if err := h.repo.Update(c.Request.Context(), tenantID, &timezone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, timezone)
}

func (h *TimezonesHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	if err := h.repo.Delete(c.Request.Context(), tenantID, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "timezone deleted"})
}

// SubdivisionsHandler handles subdivision endpoints
type SubdivisionsHandler struct {
	repo *repositories.SubdivisionRepository
}

func NewSubdivisionsHandler(repo *repositories.SubdivisionRepository) *SubdivisionsHandler {
	return &SubdivisionsHandler{repo: repo}
}

func (h *SubdivisionsHandler) GetAll(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	subdivisions, err := h.repo.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subdivisions": subdivisions, "count": len(subdivisions)})
}

func (h *SubdivisionsHandler) GetByCountry(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	countryIDStr := c.Param("countryId")
	countryID, err := uuid.Parse(countryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country ID"})
		return
	}
	subdivisions, err := h.repo.GetByCountry(c.Request.Context(), tenantID, countryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subdivisions": subdivisions, "count": len(subdivisions)})
}

func (h *SubdivisionsHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var subdivision models.CountrySubdivision
	if err := c.ShouldBindJSON(&subdivision); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), tenantID, &subdivision); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subdivision)
}

func (h *SubdivisionsHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var subdivision models.CountrySubdivision
	if err := c.ShouldBindJSON(&subdivision); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subdivision.SubdivisionID = id
	if err := h.repo.Update(c.Request.Context(), tenantID, &subdivision); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subdivision)
}

func (h *SubdivisionsHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), tenantID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "subdivision deleted"})
}

// LocalesHandler handles locale endpoints
type LocalesHandler struct {
	repo *repositories.LocaleRepository
}

func NewLocalesHandler(repo *repositories.LocaleRepository) *LocalesHandler {
	return &LocalesHandler{repo: repo}
}

func (h *LocalesHandler) GetAll(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	locales, err := h.repo.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"locales": locales, "count": len(locales)})
}

func (h *LocalesHandler) GetByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	locale, err := h.repo.GetByCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if locale == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "locale not found"})
		return
	}
	c.JSON(http.StatusOK, locale)
}

func (h *LocalesHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var locale models.Locales
	if err := c.ShouldBindJSON(&locale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), tenantID, &locale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, locale)
}

func (h *LocalesHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	var locale models.Locales
	if err := c.ShouldBindJSON(&locale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	locale.LocaleCode = code
	if err := h.repo.Update(c.Request.Context(), tenantID, &locale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locale)
}

func (h *LocalesHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")
	if err := h.repo.Delete(c.Request.Context(), tenantID, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "locale deleted"})
}