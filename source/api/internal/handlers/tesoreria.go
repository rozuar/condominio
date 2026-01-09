package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
)

type TesoreriaHandler struct {
	service *services.TesoreriaService
}

func NewTesoreriaHandler(service *services.TesoreriaService) *TesoreriaHandler {
	return &TesoreriaHandler{service: service}
}

func (h *TesoreriaHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.MovimientoFilter{
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
	if t := r.URL.Query().Get("type"); t != "" {
		filter.Type = models.MovimientoType(t)
	}
	if year := r.URL.Query().Get("year"); year != "" {
		if y, err := strconv.Atoi(year); err == nil {
			filter.Year = y
		}
	}
	if month := r.URL.Query().Get("month"); month != "" {
		if m, err := strconv.Atoi(month); err == nil {
			filter.Month = m
		}
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list movimientos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TesoreriaHandler) GetResumen(w http.ResponseWriter, r *http.Request) {
	resumen, err := h.service.GetResumen(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get resumen")
		return
	}

	writeJSON(w, http.StatusOK, resumen)
}

func (h *TesoreriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMovimientoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Description == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "Description and positive amount are required")
		return
	}
	if req.Type != models.MovimientoIngreso && req.Type != models.MovimientoEgreso {
		writeError(w, http.StatusBadRequest, "Type must be 'ingreso' or 'egreso'")
		return
	}

	createdBy := r.Context().Value("user_id").(string)

	movimiento, err := h.service.Create(r.Context(), &req, createdBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create movimiento")
		return
	}

	writeJSON(w, http.StatusCreated, movimiento)
}
