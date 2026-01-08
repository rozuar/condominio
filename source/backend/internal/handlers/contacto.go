package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
)

type ContactoHandler struct {
	service *services.ContactoService
}

func NewContactoHandler(service *services.ContactoService) *ContactoHandler {
	return &ContactoHandler{service: service}
}

func (h *ContactoHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.MensajeContactoFilter{
		Page:    1,
		PerPage: 20,
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
		filter.Status = models.ContactoStatus(status)
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list mensajes")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ContactoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mensaje, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMensajeNotFound) {
			writeError(w, http.StatusNotFound, "Mensaje not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get mensaje")
		return
	}

	writeJSON(w, http.StatusOK, mensaje)
}

func (h *ContactoHandler) GetMisMensajes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	mensajes, err := h.service.GetMisMensajes(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get mensajes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensajes": mensajes,
		"total":    len(mensajes),
	})
}

func (h *ContactoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMensajeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Email = strings.TrimSpace(req.Email)
	req.Asunto = strings.TrimSpace(req.Asunto)
	req.Mensaje = strings.TrimSpace(req.Mensaje)

	if req.Nombre == "" {
		writeError(w, http.StatusBadRequest, "nombre is required")
		return
	}
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Asunto == "" {
		writeError(w, http.StatusBadRequest, "asunto is required")
		return
	}
	if req.Mensaje == "" {
		writeError(w, http.StatusBadRequest, "mensaje is required")
		return
	}

	// Get user ID if authenticated
	var userID *string
	if uid, ok := r.Context().Value("user_id").(string); ok && uid != "" {
		userID = &uid
	}

	mensaje, err := h.service.Create(r.Context(), &req, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create mensaje")
		return
	}

	writeJSON(w, http.StatusCreated, mensaje)
}

func (h *ContactoHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	readBy := r.Context().Value("user_id").(string)

	mensaje, err := h.service.MarkAsRead(r.Context(), id, readBy)
	if err != nil {
		if errors.Is(err, services.ErrMensajeNotFound) {
			writeError(w, http.StatusNotFound, "Mensaje not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to mark as read")
		return
	}

	writeJSON(w, http.StatusOK, mensaje)
}

func (h *ContactoHandler) Reply(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.ReplyMensajeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Respuesta = strings.TrimSpace(req.Respuesta)
	if req.Respuesta == "" {
		writeError(w, http.StatusBadRequest, "respuesta is required")
		return
	}

	repliedBy := r.Context().Value("user_id").(string)

	mensaje, err := h.service.Reply(r.Context(), id, &req, repliedBy)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMensajeNotFound):
			writeError(w, http.StatusNotFound, "Mensaje not found")
		case errors.Is(err, services.ErrAlreadyReplied):
			writeError(w, http.StatusBadRequest, "Mensaje already replied")
		default:
			writeError(w, http.StatusInternalServerError, "Failed to reply")
		}
		return
	}

	writeJSON(w, http.StatusOK, mensaje)
}

func (h *ContactoHandler) Archive(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mensaje, err := h.service.Archive(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMensajeNotFound) {
			writeError(w, http.StatusNotFound, "Mensaje not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to archive")
		return
	}

	writeJSON(w, http.StatusOK, mensaje)
}

func (h *ContactoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrMensajeNotFound) {
			writeError(w, http.StatusNotFound, "Mensaje not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ContactoHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
