package models

import "time"

type ContactoStatus string

const (
	ContactoStatusPending  ContactoStatus = "pending"
	ContactoStatusRead     ContactoStatus = "read"
	ContactoStatusReplied  ContactoStatus = "replied"
	ContactoStatusArchived ContactoStatus = "archived"
)

func (s ContactoStatus) IsValid() bool {
	switch s {
	case ContactoStatusPending, ContactoStatusRead, ContactoStatusReplied, ContactoStatusArchived:
		return true
	}
	return false
}

type MensajeContacto struct {
	ID         string         `json:"id"`
	UserID     *string        `json:"user_id,omitempty"`
	UserName   string         `json:"user_name,omitempty"`
	Nombre     string         `json:"nombre"`
	Email      string         `json:"email"`
	Asunto     string         `json:"asunto"`
	Mensaje    string         `json:"mensaje"`
	Status     ContactoStatus `json:"status"`
	ReadAt     *time.Time     `json:"read_at,omitempty"`
	ReadBy     *string        `json:"read_by,omitempty"`
	ReadByName string         `json:"read_by_name,omitempty"`
	RepliedAt  *time.Time     `json:"replied_at,omitempty"`
	RepliedBy  *string        `json:"replied_by,omitempty"`
	RepliedByName string      `json:"replied_by_name,omitempty"`
	Respuesta  string         `json:"respuesta,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type CreateMensajeRequest struct {
	Nombre  string `json:"nombre"`
	Email   string `json:"email"`
	Asunto  string `json:"asunto"`
	Mensaje string `json:"mensaje"`
}

type ReplyMensajeRequest struct {
	Respuesta string `json:"respuesta"`
}

type MensajeContactoListResponse struct {
	Mensajes []MensajeContacto `json:"mensajes"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
}

type MensajeContactoFilter struct {
	Status  ContactoStatus
	UserID  string
	Page    int
	PerPage int
}

type ContactoStats struct {
	TotalPending  int `json:"total_pending"`
	TotalRead     int `json:"total_read"`
	TotalReplied  int `json:"total_replied"`
	TotalArchived int `json:"total_archived"`
	Total         int `json:"total"`
}
