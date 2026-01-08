package models

import "time"

type MovimientoType string

const (
	MovimientoIngreso MovimientoType = "ingreso"
	MovimientoEgreso  MovimientoType = "egreso"
)

type Movimiento struct {
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Amount      float64        `json:"amount"`
	Type        MovimientoType `json:"type"`
	Category    string         `json:"category,omitempty"`
	Date        time.Time      `json:"date"`
	CreatedBy   *string        `json:"created_by,omitempty"`
	CreatorName string         `json:"creator_name,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateMovimientoRequest struct {
	Description string         `json:"description"`
	Amount      float64        `json:"amount"`
	Type        MovimientoType `json:"type"`
	Category    string         `json:"category"`
	Date        time.Time      `json:"date"`
}

type MovimientoListResponse struct {
	Movimientos []Movimiento `json:"movimientos"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
}

type ResumenTesoreria struct {
	TotalIngresos float64 `json:"total_ingresos"`
	TotalEgresos  float64 `json:"total_egresos"`
	Balance       float64 `json:"balance"`
}

type MovimientoFilter struct {
	Type    MovimientoType
	Year    int
	Month   int
	Page    int
	PerPage int
}
