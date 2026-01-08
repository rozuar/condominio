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

type ActaHandler struct {
	service *services.ActaService
}

func NewActaHandler(service *services.ActaService) *ActaHandler {
	return &ActaHandler{service: service}
}

func (h *ActaHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.ActaFilter{
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
		filter.Type = models.ActaType(t)
	}
	if year := r.URL.Query().Get("year"); year != "" {
		if y, err := strconv.Atoi(year); err == nil {
			filter.Year = y
		}
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list actas")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ActaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acta, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrActaNotFound) {
			writeError(w, http.StatusNotFound, "Acta not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get acta")
		return
	}

	writeJSON(w, http.StatusOK, acta)
}

func (h *ActaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateActaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}
	if req.Type == "" {
		req.Type = models.ActaOrdinaria
	}

	createdBy := r.Context().Value("user_id").(string)

	acta, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create acta")
		return
	}

	writeJSON(w, http.StatusCreated, acta)
}
