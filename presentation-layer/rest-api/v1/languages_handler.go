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

type LanguagesHandler struct {
	languageService *applicationservices.LanguageAppService
}

func NewLanguagesHandler(languageService *applicationservices.LanguageAppService) *LanguagesHandler {
	return &LanguagesHandler{languageService: languageService}
}

func (h *LanguagesHandler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	var input models.LanguageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	language, err := h.languageService.CreateLanguage(r.Context(), tenantID, &input)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	h.writeJSON(w, http.StatusCreated, language)
}

func (h *LanguagesHandler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid language ID")
		return
	}
	
	language, err := h.languageService.GetLanguage(r.Context(), tenantID, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Language not found")
		return
	}
	
	h.writeJSON(w, http.StatusOK, language)
}

func (h *LanguagesHandler) ListLanguages(w http.ResponseWriter, r *http.Request) {
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
	
	languages, total, err := h.languageService.ListLanguages(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response := map[string]interface{}{
		"languages": languages,
		"count":     total,
		"limit":     limit,
		"offset":    offset,
	}
	
	h.writeJSON(w, http.StatusOK, response)
}

func (h *LanguagesHandler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid language ID")
		return
	}
	
	var input models.LanguageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	language, err := h.languageService.UpdateLanguage(r.Context(), tenantID, id, &input)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	h.writeJSON(w, http.StatusOK, language)
}

func (h *LanguagesHandler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.writeError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}
	
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid language ID")
		return
	}
	
	if err := h.languageService.DeleteLanguage(r.Context(), tenantID, id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *LanguagesHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *LanguagesHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}