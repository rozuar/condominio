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

type ComunicadoHandler struct {
	service *services.ComunicadoService
}

func NewComunicadoHandler(service *services.ComunicadoService) *ComunicadoHandler {
	return &ComunicadoHandler{service: service}
}

func (h *ComunicadoHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.ComunicadoFilter{
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
	if t := r.URL.Query().Get("type"); t != "" {
		filter.Type = models.ComunicadoType(t)
	}

	// Check if user is authenticated
	userRole := r.Context().Value("user_role")
	if userRole == nil {
		// Public access - only show public comunicados
		isPublic := true
		filter.IsPublic = &isPublic
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list comunicados")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ComunicadoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	comunicado, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrComunicadoNotFound) {
			writeError(w, http.StatusNotFound, "Comunicado not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get comunicado")
		return
	}

	// Check access
	userRole := r.Context().Value("user_role")
	if !comunicado.IsPublic && userRole == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	writeJSON(w, http.StatusOK, comunicado)
}

func (h *ComunicadoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateComunicadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	if req.Type == "" {
		req.Type = models.ComunicadoInformativo
	}

	authorID := r.Context().Value("user_id").(string)

	comunicado, err := h.service.Create(r.Context(), &req, authorID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create comunicado")
		return
	}

	writeJSON(w, http.StatusCreated, comunicado)
}

func (h *ComunicadoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateComunicadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	comunicado, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrComunicadoNotFound) {
			writeError(w, http.StatusNotFound, "Comunicado not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update comunicado")
		return
	}

	writeJSON(w, http.StatusOK, comunicado)
}

func (h *ComunicadoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrComunicadoNotFound) {
			writeError(w, http.StatusNotFound, "Comunicado not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete comunicado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ComunicadoHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	limit := 3
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10 {
			limit = parsed
		}
	}

	userRole := r.Context().Value("user_role")
	publicOnly := userRole == nil

	comunicados, err := h.service.GetLatest(r.Context(), limit, publicOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get latest comunicados")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"comunicados": comunicados,
	})
}
