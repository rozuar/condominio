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

type GaleriaHandler struct {
	service *services.GaleriaService
}

func NewGaleriaHandler(service *services.GaleriaService) *GaleriaHandler {
	return &GaleriaHandler{service: service}
}

func (h *GaleriaHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.GaleriaFilter{
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
	if isPublic := r.URL.Query().Get("is_public"); isPublic != "" {
		b := isPublic == "true"
		filter.IsPublic = &b
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list galerias")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *GaleriaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	galeria, err := h.service.GetWithItems(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaNotFound) {
			writeError(w, http.StatusNotFound, "Galeria not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get galeria")
		return
	}

	writeJSON(w, http.StatusOK, galeria)
}

func (h *GaleriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGaleriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	createdBy := r.Context().Value("user_id").(string)

	galeria, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create galeria")
		return
	}

	writeJSON(w, http.StatusCreated, galeria)
}

func (h *GaleriaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateGaleriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	galeria, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaNotFound) {
			writeError(w, http.StatusNotFound, "Galeria not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update galeria")
		return
	}

	writeJSON(w, http.StatusOK, galeria)
}

func (h *GaleriaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaNotFound) {
			writeError(w, http.StatusNotFound, "Galeria not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete galeria")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Galeria deleted successfully"})
}

// Item handlers

func (h *GaleriaHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	galeriaID := chi.URLParam(r, "id")

	var req models.AddGaleriaItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FileURL == "" {
		writeError(w, http.StatusBadRequest, "File URL is required")
		return
	}
	if !req.FileType.IsValid() {
		writeError(w, http.StatusBadRequest, "Invalid file type (must be 'image' or 'video')")
		return
	}

	item, err := h.service.AddItem(r.Context(), galeriaID, &req)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaNotFound) {
			writeError(w, http.StatusNotFound, "Galeria not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to add item")
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *GaleriaHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateGaleriaItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FileType != nil && !req.FileType.IsValid() {
		writeError(w, http.StatusBadRequest, "Invalid file type (must be 'image' or 'video')")
		return
	}

	item, err := h.service.UpdateItem(r.Context(), itemID, &req)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaItemNotFound) {
			writeError(w, http.StatusNotFound, "Item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *GaleriaHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")

	err := h.service.DeleteItem(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaItemNotFound) {
			writeError(w, http.StatusNotFound, "Item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Item deleted successfully"})
}

func (h *GaleriaHandler) ReorderItems(w http.ResponseWriter, r *http.Request) {
	galeriaID := chi.URLParam(r, "id")

	var req struct {
		ItemIDs []string `json:"item_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.ItemIDs) == 0 {
		writeError(w, http.StatusBadRequest, "Item IDs are required")
		return
	}

	err := h.service.ReorderItems(r.Context(), galeriaID, req.ItemIDs)
	if err != nil {
		if errors.Is(err, services.ErrGaleriaNotFound) {
			writeError(w, http.StatusNotFound, "Galeria not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to reorder items")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Items reordered successfully"})
}
