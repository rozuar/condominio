package models

import "time"

type DocumentoCategory string

const (
	DocumentoReglamento DocumentoCategory = "reglamento"
	DocumentoProtocolo  DocumentoCategory = "protocolo"
	DocumentoFormulario DocumentoCategory = "formulario"
	DocumentoOtro       DocumentoCategory = "otro"
)

type Documento struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	FileURL     string            `json:"file_url,omitempty"`
	Category    DocumentoCategory `json:"category"`
	IsPublic    bool              `json:"is_public"`
	CreatedBy   *string           `json:"created_by,omitempty"`
	CreatorName string            `json:"creator_name,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type CreateDocumentoRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	FileURL     string            `json:"file_url"`
	Category    DocumentoCategory `json:"category"`
	IsPublic    bool              `json:"is_public"`
}

type DocumentoListResponse struct {
	Documentos []Documento `json:"documentos"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
}

type DocumentoFilter struct {
	Category DocumentoCategory
	IsPublic *bool
	Page     int
	PerPage  int
}
