package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
)

type VotacionHandler struct {
	service *services.VotacionService
}

func NewVotacionHandler(service *services.VotacionService) *VotacionHandler {
	return &VotacionHandler{service: service}
}

func (h *VotacionHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.VotacionFilter{
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
		filter.Status = models.VotacionStatus(status)
	}
	if active := r.URL.Query().Get("active"); active == "true" {
		filter.Active = true
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list votaciones")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *VotacionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	var userID *string
	if uid, ok := r.Context().Value("user_id").(string); ok {
		userID = &uid
	}

	votaciones, err := h.service.GetActive(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get active votaciones")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"votaciones": votaciones,
		"total":      len(votaciones),
	})
}

func (h *VotacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userID *string
	if uid, ok := r.Context().Value("user_id").(string); ok {
		userID = &uid
	}

	votacion, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get votacion")
		return
	}

	writeJSON(w, http.StatusOK, votacion)
}

func (h *VotacionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVotacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if len(req.Opciones) < 2 {
		writeError(w, http.StatusBadRequest, "At least 2 options are required")
		return
	}

	createdBy := r.Context().Value("user_id").(string)

	votacion, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create votacion")
		return
	}

	writeJSON(w, http.StatusCreated, votacion)
}

func (h *VotacionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateVotacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	votacion, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, votacion)
}

type PublishRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

func (h *VotacionHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req PublishRequest
	json.NewDecoder(r.Body).Decode(&req) // Optional body

	votacion, err := h.service.Publish(r.Context(), id, req.StartDate, req.EndDate)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, votacion)
}

func (h *VotacionHandler) Close(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	votacion, err := h.service.Close(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to close votacion")
		return
	}

	writeJSON(w, http.StatusOK, votacion)
}

func (h *VotacionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	votacion, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to cancel votacion")
		return
	}

	writeJSON(w, http.StatusOK, votacion)
}

func (h *VotacionHandler) EmitirVoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.EmitirVotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !req.IsAbstention && req.OpcionID == nil {
		writeError(w, http.StatusBadRequest, "opcion_id is required unless abstaining")
		return
	}

	userID := r.Context().Value("user_id").(string)

	err := h.service.EmitirVoto(r.Context(), id, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVotacionNotFound):
			writeError(w, http.StatusNotFound, "Votacion not found")
		case errors.Is(err, services.ErrVotacionNotActive):
			writeError(w, http.StatusBadRequest, "Votacion is not active")
		case errors.Is(err, services.ErrAlreadyVoted):
			writeError(w, http.StatusConflict, "You have already voted")
		case errors.Is(err, services.ErrInvalidOpcion):
			writeError(w, http.StatusBadRequest, "Invalid option")
		case errors.Is(err, services.ErrAbstentionNotAllowed):
			writeError(w, http.StatusBadRequest, "Abstention is not allowed")
		default:
			writeError(w, http.StatusInternalServerError, "Failed to emit vote")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Vote registered successfully"})
}

func (h *VotacionHandler) GetResultados(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resultado, err := h.service.GetResultados(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrVotacionNotFound) {
			writeError(w, http.StatusNotFound, "Votacion not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get results")
		return
	}

	writeJSON(w, http.StatusOK, resultado)
}

func (h *VotacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
