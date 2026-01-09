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

type EventoHandler struct {
	service *services.EventoService
}

func NewEventoHandler(service *services.EventoService) *EventoHandler {
	return &EventoHandler{service: service}
}

func (h *EventoHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.EventoFilter{
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
		filter.Type = models.EventoType(t)
	}
	if r.URL.Query().Get("upcoming") == "true" {
		filter.Upcoming = true
	}

	userRole := r.Context().Value("user_role")
	if userRole == nil {
		isPublic := true
		filter.IsPublic = &isPublic
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list eventos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *EventoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	evento, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrEventoNotFound) {
			writeError(w, http.StatusNotFound, "Evento not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get evento")
		return
	}

	userRole := r.Context().Value("user_role")
	if !evento.IsPublic && userRole == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	writeJSON(w, http.StatusOK, evento)
}

func (h *EventoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.EventDate.IsZero() {
		writeError(w, http.StatusBadRequest, "Event date is required")
		return
	}

	if req.Type == "" {
		req.Type = models.EventoReunion
	}

	createdBy := r.Context().Value("user_id").(string)

	evento, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create evento")
		return
	}

	writeJSON(w, http.StatusCreated, evento)
}

func (h *EventoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateEventoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	evento, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrEventoNotFound) {
			writeError(w, http.StatusNotFound, "Evento not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update evento")
		return
	}

	writeJSON(w, http.StatusOK, evento)
}

func (h *EventoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrEventoNotFound) {
			writeError(w, http.StatusNotFound, "Evento not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete evento")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventoHandler) GetUpcoming(w http.ResponseWriter, r *http.Request) {
	limit := 3
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10 {
			limit = parsed
		}
	}

	userRole := r.Context().Value("user_role")
	publicOnly := userRole == nil

	eventos, err := h.service.GetUpcoming(r.Context(), limit, publicOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get upcoming eventos")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"eventos": eventos,
	})
}
