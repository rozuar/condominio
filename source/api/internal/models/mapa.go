package models

import (
	"encoding/json"
	"time"
)

type AreaType string

const (
	AreaTypeParcela   AreaType = "parcela"
	AreaTypeAreaComun AreaType = "area_comun"
	AreaTypeAcceso    AreaType = "acceso"
	AreaTypeCanal     AreaType = "canal"
	AreaTypeCamino    AreaType = "camino"
)

type MapaArea struct {
	ID          string          `json:"id"`
	ParcelaID   *int            `json:"parcela_id,omitempty"`
	Type        AreaType        `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Coordinates json.RawMessage `json:"coordinates"` // GeoJSON polygon coordinates
	CenterLat   *float64        `json:"center_lat,omitempty"`
	CenterLng   *float64        `json:"center_lng,omitempty"`
	FillColor   string          `json:"fill_color"`
	StrokeColor string          `json:"stroke_color"`
	IsClickable bool            `json:"is_clickable"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type MapaPunto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Icon        string    `json:"icon"`
	Type        string    `json:"type"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMapaAreaRequest struct {
	ParcelaID   *int            `json:"parcela_id,omitempty"`
	Type        AreaType        `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Coordinates json.RawMessage `json:"coordinates"`
	CenterLat   *float64        `json:"center_lat,omitempty"`
	CenterLng   *float64        `json:"center_lng,omitempty"`
	FillColor   string          `json:"fill_color,omitempty"`
	StrokeColor string          `json:"stroke_color,omitempty"`
	IsClickable *bool           `json:"is_clickable,omitempty"`
}

type UpdateMapaAreaRequest struct {
	ParcelaID   *int             `json:"parcela_id,omitempty"`
	Type        *AreaType        `json:"type,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Coordinates *json.RawMessage `json:"coordinates,omitempty"`
	CenterLat   *float64         `json:"center_lat,omitempty"`
	CenterLng   *float64         `json:"center_lng,omitempty"`
	FillColor   *string          `json:"fill_color,omitempty"`
	StrokeColor *string          `json:"stroke_color,omitempty"`
	IsClickable *bool            `json:"is_clickable,omitempty"`
}

type CreateMapaPuntoRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Icon        string  `json:"icon,omitempty"`
	Type        string  `json:"type"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

type UpdateMapaPuntoRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Lat         *float64 `json:"lat,omitempty"`
	Lng         *float64 `json:"lng,omitempty"`
	Icon        *string  `json:"icon,omitempty"`
	Type        *string  `json:"type,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
}

type MapaAreaListResponse struct {
	Areas   []MapaArea `json:"areas"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}

type MapaPuntoListResponse struct {
	Puntos  []MapaPunto `json:"puntos"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

type MapaAreaFilter struct {
	Type    AreaType
	Page    int
	PerPage int
}

type MapaPuntoFilter struct {
	Type     string
	IsPublic *bool
	Page     int
	PerPage  int
}

// MapaData represents all map data for the frontend
type MapaData struct {
	Areas  []MapaArea  `json:"areas"`
	Puntos []MapaPunto `json:"puntos"`
}

func (at AreaType) IsValid() bool {
	switch at {
	case AreaTypeParcela, AreaTypeAreaComun, AreaTypeAcceso, AreaTypeCanal, AreaTypeCamino:
		return true
	}
	return false
}
