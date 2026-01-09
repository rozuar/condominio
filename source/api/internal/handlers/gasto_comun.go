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

type GastoComunHandler struct {
	service *services.GastoComunService
}

func NewGastoComunHandler(service *services.GastoComunService) *GastoComunHandler {
	return &GastoComunHandler{service: service}
}

// ============================================
// PERIODOS
// ============================================

func (h *GastoComunHandler) ListPeriodos(w http.ResponseWriter, r *http.Request) {
	filter := models.PeriodoFilter{
		Page:    1,
		PerPage: 12,
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
	if year := r.URL.Query().Get("year"); year != "" {
		if y, err := strconv.Atoi(year); err == nil {
			filter.Year = y
		}
	}

	resp, err := h.service.ListPeriodos(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list periodos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *GastoComunHandler) GetPeriodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	periodo, err := h.service.GetPeriodo(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPeriodoNotFound) {
			writeError(w, http.StatusNotFound, "Periodo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get periodo")
		return
	}

	writeJSON(w, http.StatusOK, periodo)
}

func (h *GastoComunHandler) GetPeriodoActual(w http.ResponseWriter, r *http.Request) {
	periodo, err := h.service.GetPeriodoActual(r.Context())
	if err != nil {
		if errors.Is(err, services.ErrPeriodoNotFound) {
			writeError(w, http.StatusNotFound, "No periodo for current month")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get periodo actual")
		return
	}

	writeJSON(w, http.StatusOK, periodo)
}

func (h *GastoComunHandler) CreatePeriodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePeriodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Year < 2020 || req.Year > 2100 {
		writeError(w, http.StatusBadRequest, "Invalid year")
		return
	}
	if req.Month < 1 || req.Month > 12 {
		writeError(w, http.StatusBadRequest, "Invalid month")
		return
	}
	if req.MontoBase <= 0 {
		writeError(w, http.StatusBadRequest, "monto_base must be positive")
		return
	}
	if req.FechaVencimiento == "" {
		writeError(w, http.StatusBadRequest, "fecha_vencimiento is required")
		return
	}

	periodo, err := h.service.CreatePeriodo(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrPeriodoExists) {
			writeError(w, http.StatusConflict, "Periodo already exists for this year/month")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, periodo)
}

func (h *GastoComunHandler) UpdatePeriodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdatePeriodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	periodo, err := h.service.UpdatePeriodo(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrPeriodoNotFound) {
			writeError(w, http.StatusNotFound, "Periodo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, periodo)
}

func (h *GastoComunHandler) GetResumen(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resumen, err := h.service.GetResumen(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPeriodoNotFound) {
			writeError(w, http.StatusNotFound, "Periodo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get resumen")
		return
	}

	writeJSON(w, http.StatusOK, resumen)
}

// ============================================
// GASTOS COMUNES
// ============================================

func (h *GastoComunHandler) ListGastos(w http.ResponseWriter, r *http.Request) {
	filter := models.GastoComunFilter{
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
	if periodoID := r.URL.Query().Get("periodo_id"); periodoID != "" {
		filter.PeriodoID = periodoID
	}
	if parcelaID := r.URL.Query().Get("parcela_id"); parcelaID != "" {
		if pid, err := strconv.Atoi(parcelaID); err == nil {
			filter.ParcelaID = pid
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.PagoStatus(status)
	}

	resp, err := h.service.ListGastos(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list gastos")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *GastoComunHandler) GetGasto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	gasto, err := h.service.GetGasto(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrGastoComunNotFound) {
			writeError(w, http.StatusNotFound, "Gasto not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get gasto")
		return
	}

	writeJSON(w, http.StatusOK, gasto)
}

func (h *GastoComunHandler) GetMiEstadoCuenta(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	estado, err := h.service.GetMiEstadoCuenta(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, estado)
}

func (h *GastoComunHandler) RegistrarPago(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.RegistrarPagoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Monto <= 0 {
		writeError(w, http.StatusBadRequest, "monto must be positive")
		return
	}
	if req.Metodo == "" {
		writeError(w, http.StatusBadRequest, "metodo is required")
		return
	}

	gasto, err := h.service.RegistrarPago(r.Context(), id, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGastoComunNotFound):
			writeError(w, http.StatusNotFound, "Gasto not found")
		case errors.Is(err, services.ErrGastoAlreadyPaid):
			writeError(w, http.StatusBadRequest, "Gasto already paid")
		case errors.Is(err, services.ErrInvalidPaymentAmount):
			writeError(w, http.StatusBadRequest, "Invalid payment amount")
		default:
			writeError(w, http.StatusInternalServerError, "Failed to register payment")
		}
		return
	}

	writeJSON(w, http.StatusOK, gasto)
}

func (h *GastoComunHandler) MarcarVencidos(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.MarcarVencidos(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to mark overdue")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"marked_overdue": count})
}
