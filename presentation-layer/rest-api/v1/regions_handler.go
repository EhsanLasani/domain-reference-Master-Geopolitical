package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	applicationservices "github.com/EhsanLasani/domain-reference-Master-Geopolitical/business-logic-layer/application-services"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/models"
)

type RegionsHandler struct {
	regionService *applicationservices.RegionAppService
}

func NewRegionsHandler(regionService *applicationservices.RegionAppService) *RegionsHandler {
	return &RegionsHandler{regionService: regionService}
}

func (h *RegionsHandler) CreateRegion(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	var input models.RegionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	region, err := h.regionService.CreateRegion(r.Context(), tenantID, &input)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	h.writeJSON(w, http.StatusCreated, region)
}

func (h *RegionsHandler) GetRegion(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid region ID")
		return
	}
	
	region, err := h.regionService.GetRegion(r.Context(), tenantID, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Region not found")
		return
	}
	
	h.writeJSON(w, http.StatusOK, region)
}

func (h *RegionsHandler) ListRegions(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	
	regions, total, err := h.regionService.ListRegions(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response := map[string]interface{}{
		"regions": regions,
		"count":   total,
		"limit":   limit,
		"offset":  offset,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

func (h *RegionsHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid region ID")
		return
	}
	
	var input models.RegionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	region, err := h.regionService.UpdateRegion(r.Context(), tenantID, id, &input)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	h.writeJSON(w, http.StatusOK, region)
}

func (h *RegionsHandler) DeleteRegion(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid region ID")
		return
	}
	
	if err := h.regionService.DeleteRegion(r.Context(), tenantID, id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *RegionsHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *RegionsHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}