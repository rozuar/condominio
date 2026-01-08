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

type DocumentoHandler struct {
	service *services.DocumentoService
}

func NewDocumentoHandler(service *services.DocumentoService) *DocumentoHandler {
	return &DocumentoHandler{service: service}
}

func (h *DocumentoHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.DocumentoFilter{
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
	if cat := r.URL.Query().Get("category"); cat != "" {
		filter.Category = models.DocumentoCategory(cat)
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list documentos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DocumentoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	documento, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrDocumentoNotFound) {
			writeError(w, http.StatusNotFound, "Documento not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get documento")
		return
	}

	writeJSON(w, http.StatusOK, documento)
}

func (h *DocumentoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDocumentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.Category == "" {
		req.Category = models.DocumentoOtro
	}

	createdBy := r.Context().Value("user_id").(string)

	documento, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create documento")
		return
	}

	writeJSON(w, http.StatusCreated, documento)
}
