package models

import "time"

type ComunicadoType string

const (
	ComunicadoInformativo ComunicadoType = "informativo"
	ComunicadoSeguridad   ComunicadoType = "seguridad"
	ComunicadoTesoreria   ComunicadoType = "tesoreria"
	ComunicadoAsamblea    ComunicadoType = "asamblea"
)

type Comunicado struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	Type        ComunicadoType `json:"type"`
	IsPublic    bool           `json:"is_public"`
	AuthorID    *string        `json:"author_id,omitempty"`
	AuthorName  string         `json:"author_name,omitempty"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateComunicadoRequest struct {
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	Type     ComunicadoType `json:"type"`
	IsPublic bool           `json:"is_public"`
	Publish  bool           `json:"publish"`
}

type UpdateComunicadoRequest struct {
	Title    *string         `json:"title,omitempty"`
	Content  *string         `json:"content,omitempty"`
	Type     *ComunicadoType `json:"type,omitempty"`
	IsPublic *bool           `json:"is_public,omitempty"`
	Publish  *bool           `json:"publish,omitempty"`
}

type ComunicadoListResponse struct {
	Comunicados []Comunicado `json:"comunicados"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
}

type ComunicadoFilter struct {
	Type     ComunicadoType
	IsPublic *bool
	Page     int
	PerPage  int
}
