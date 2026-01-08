package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
)

type EmergenciaHandler struct {
	service *services.EmergenciaService
}

func NewEmergenciaHandler(service *services.EmergenciaService) *EmergenciaHandler {
	return &EmergenciaHandler{service: service}
}

func (h *EmergenciaHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.EmergenciaFilter{
		Page:    1,
		PerPage: 10,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}
	if perPage := r.URL.Query().Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil {
			filter.PerPage = pp
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.EmergenciaStatus(status)
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		filter.Priority = models.EmergenciaPriority(priority)
	}
	if active := r.URL.Query().Get("active"); active == "true" {
		filter.Active = true
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list emergencias")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *EmergenciaHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	emergencias, err := h.service.GetActive(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get active emergencias")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"emergencias": emergencias,
		"total":       len(emergencias),
	})
}

func (h *EmergenciaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	emergencia, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrEmergenciaNotFound) {
			writeError(w, http.StatusNotFound, "Emergencia not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get emergencia")
		return
	}

	writeJSON(w, http.StatusOK, emergencia)
}

func (h *EmergenciaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmergenciaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	if !req.Priority.IsValid() {
		req.Priority = models.EmergenciaPriorityMedium
	}

	createdBy := r.Context().Value("user_id").(string)

	emergencia, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create emergencia")
		return
	}

	writeJSON(w, http.StatusCreated, emergencia)
}

func (h *EmergenciaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateEmergenciaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	emergencia, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrEmergenciaNotFound) {
			writeError(w, http.StatusNotFound, "Emergencia not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update emergencia")
		return
	}

	writeJSON(w, http.StatusOK, emergencia)
}

func (h *EmergenciaHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resolvedBy := r.Context().Value("user_id").(string)

	emergencia, err := h.service.Resolve(r.Context(), id, resolvedBy)
	if err != nil {
		if errors.Is(err, services.ErrEmergenciaNotFound) {
			writeError(w, http.StatusNotFound, "Emergencia not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to resolve emergencia")
		return
	}

	writeJSON(w, http.StatusOK, emergencia)
}

func (h *EmergenciaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrEmergenciaNotFound) {
			writeError(w, http.StatusNotFound, "Emergencia not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete emergencia")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
