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

type MapaHandler struct {
	service *services.MapaService
}

func NewMapaHandler(service *services.MapaService) *MapaHandler {
	return &MapaHandler{service: service}
}

// GetAllMapData returns all map data for the frontend (public)
func (h *MapaHandler) GetAllMapData(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAllMapData(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get map data")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// Area handlers

func (h *MapaHandler) ListAreas(w http.ResponseWriter, r *http.Request) {
	filter := models.MapaAreaFilter{
		Page:    1,
		PerPage: 100,
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
	if areaType := r.URL.Query().Get("type"); areaType != "" {
		filter.Type = models.AreaType(areaType)
	}

	resp, err := h.service.ListAreas(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list areas")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *MapaHandler) GetAreaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	area, err := h.service.GetAreaByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMapaAreaNotFound) {
			writeError(w, http.StatusNotFound, "Area not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get area")
		return
	}

	writeJSON(w, http.StatusOK, area)
}

func (h *MapaHandler) CreateArea(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMapaAreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if !req.Type.IsValid() {
		writeError(w, http.StatusBadRequest, "Invalid area type")
		return
	}
	if len(req.Coordinates) == 0 {
		writeError(w, http.StatusBadRequest, "Coordinates are required")
		return
	}

	area, err := h.service.CreateArea(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create area")
		return
	}

	writeJSON(w, http.StatusCreated, area)
}

func (h *MapaHandler) UpdateArea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateMapaAreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Type != nil && !req.Type.IsValid() {
		writeError(w, http.StatusBadRequest, "Invalid area type")
		return
	}

	area, err := h.service.UpdateArea(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrMapaAreaNotFound) {
			writeError(w, http.StatusNotFound, "Area not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update area")
		return
	}

	writeJSON(w, http.StatusOK, area)
}

func (h *MapaHandler) DeleteArea(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteArea(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMapaAreaNotFound) {
			writeError(w, http.StatusNotFound, "Area not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete area")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Area deleted successfully"})
}

// Punto handlers

func (h *MapaHandler) ListPuntos(w http.ResponseWriter, r *http.Request) {
	filter := models.MapaPuntoFilter{
		Page:    1,
		PerPage: 100,
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
	if puntoType := r.URL.Query().Get("type"); puntoType != "" {
		filter.Type = puntoType
	}
	if isPublic := r.URL.Query().Get("is_public"); isPublic != "" {
		b := isPublic == "true"
		filter.IsPublic = &b
	}

	resp, err := h.service.ListPuntos(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list puntos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *MapaHandler) GetPuntoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	punto, err := h.service.GetPuntoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMapaPuntoNotFound) {
			writeError(w, http.StatusNotFound, "Punto not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get punto")
		return
	}

	writeJSON(w, http.StatusOK, punto)
}

func (h *MapaHandler) CreatePunto(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMapaPuntoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Type == "" {
		writeError(w, http.StatusBadRequest, "Type is required")
		return
	}

	punto, err := h.service.CreatePunto(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create punto")
		return
	}

	writeJSON(w, http.StatusCreated, punto)
}

func (h *MapaHandler) UpdatePunto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateMapaPuntoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	punto, err := h.service.UpdatePunto(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrMapaPuntoNotFound) {
			writeError(w, http.StatusNotFound, "Punto not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update punto")
		return
	}

	writeJSON(w, http.StatusOK, punto)
}

func (h *MapaHandler) DeletePunto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeletePunto(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMapaPuntoNotFound) {
			writeError(w, http.StatusNotFound, "Punto not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete punto")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Punto deleted successfully"})
}
