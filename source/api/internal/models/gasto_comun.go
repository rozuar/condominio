package models

import "time"

type PagoStatus string

const (
	PagoStatusPending   PagoStatus = "pending"
	PagoStatusPaid      PagoStatus = "paid"
	PagoStatusOverdue   PagoStatus = "overdue"
	PagoStatusCancelled PagoStatus = "cancelled"
)

func (s PagoStatus) IsValid() bool {
	switch s {
	case PagoStatusPending, PagoStatusPaid, PagoStatusOverdue, PagoStatusCancelled:
		return true
	}
	return false
}

type PeriodoGasto struct {
	ID               string    `json:"id"`
	Year             int       `json:"year"`
	Month            int       `json:"month"`
	MontoBase        float64   `json:"monto_base"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Descripcion      string    `json:"descripcion,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	// Computed fields
	TotalParcelas   int     `json:"total_parcelas,omitempty"`
	TotalPagados    int     `json:"total_pagados,omitempty"`
	TotalPendientes int     `json:"total_pendientes,omitempty"`
	MontoRecaudado  float64 `json:"monto_recaudado,omitempty"`
	MontoPendiente  float64 `json:"monto_pendiente,omitempty"`
}

type GastoComun struct {
	ID             string     `json:"id"`
	PeriodoID      string     `json:"periodo_id"`
	ParcelaID      int        `json:"parcela_id"`
	ParcelaNumero  string     `json:"parcela_numero,omitempty"`
	UserID         *string    `json:"user_id,omitempty"`
	UserName       string     `json:"user_name,omitempty"`
	Monto          float64    `json:"monto"`
	MontoPagado    float64    `json:"monto_pagado"`
	Status         PagoStatus `json:"status"`
	FechaPago      *time.Time `json:"fecha_pago,omitempty"`
	MetodoPago     string     `json:"metodo_pago,omitempty"`
	ReferenciaPago string     `json:"referencia_pago,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Related data
	Periodo *PeriodoGasto `json:"periodo,omitempty"`
}

type Pago struct {
	ID                string    `json:"id"`
	GastoComunID      string    `json:"gasto_comun_id"`
	Monto             float64   `json:"monto"`
	Metodo            string    `json:"metodo"` // transbank, mercadopago, transferencia, efectivo
	ReferenciaExterna string    `json:"referencia_externa,omitempty"`
	Estado            string    `json:"estado"` // pending, approved, rejected
	Detalles          string    `json:"detalles,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

// Request/Response types

type CreatePeriodoRequest struct {
	Year             int     `json:"year"`
	Month            int     `json:"month"`
	MontoBase        float64 `json:"monto_base"`
	FechaVencimiento string  `json:"fecha_vencimiento"` // YYYY-MM-DD
	Descripcion      string  `json:"descripcion,omitempty"`
}

type UpdatePeriodoRequest struct {
	MontoBase        *float64 `json:"monto_base,omitempty"`
	FechaVencimiento *string  `json:"fecha_vencimiento,omitempty"`
	Descripcion      *string  `json:"descripcion,omitempty"`
}

type RegistrarPagoRequest struct {
	Monto             float64 `json:"monto"`
	Metodo            string  `json:"metodo"` // transbank, mercadopago, transferencia, efectivo
	ReferenciaExterna string  `json:"referencia_externa,omitempty"`
}

type PeriodoListResponse struct {
	Periodos []PeriodoGasto `json:"periodos"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PerPage  int            `json:"per_page"`
}

type GastoComunListResponse struct {
	Gastos  []GastoComun `json:"gastos"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	PerPage int          `json:"per_page"`
}

type PeriodoFilter struct {
	Year    int
	Page    int
	PerPage int
}

type GastoComunFilter struct {
	PeriodoID string
	ParcelaID int
	Status    PagoStatus
	Page      int
	PerPage   int
}

type ResumenGastos struct {
	Periodo            PeriodoGasto `json:"periodo"`
	TotalParcelas      int          `json:"total_parcelas"`
	TotalPagados       int          `json:"total_pagados"`
	TotalPendientes    int          `json:"total_pendientes"`
	TotalVencidos      int          `json:"total_vencidos"`
	MontoTotal         float64      `json:"monto_total"`
	MontoRecaudado     float64      `json:"monto_recaudado"`
	MontoPendiente     float64      `json:"monto_pendiente"`
	PorcentajeRecaudo  float64      `json:"porcentaje_recaudo"`
}

type MiEstadoCuenta struct {
	ParcelaID      int            `json:"parcela_id"`
	ParcelaNumero  string         `json:"parcela_numero"`
	GastosPendientes []GastoComun `json:"gastos_pendientes"`
	GastosPagados    []GastoComun `json:"gastos_pagados"`
	TotalPendiente   float64      `json:"total_pendiente"`
	TotalPagado      float64      `json:"total_pagado"`
}
