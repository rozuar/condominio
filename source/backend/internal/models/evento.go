package models

import "time"

type EventoType string

const (
	EventoReunion  EventoType = "reunion"
	EventoAsamblea EventoType = "asamblea"
	EventoTrabajo  EventoType = "trabajo"
	EventoSocial   EventoType = "social"
)

type Evento struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	EventDate    time.Time  `json:"event_date"`
	EventEndDate *time.Time `json:"event_end_date,omitempty"`
	Location     string     `json:"location,omitempty"`
	Type         EventoType `json:"type"`
	IsPublic     bool       `json:"is_public"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatorName  string     `json:"creator_name,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateEventoRequest struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	EventDate    time.Time  `json:"event_date"`
	EventEndDate *time.Time `json:"event_end_date,omitempty"`
	Location     string     `json:"location"`
	Type         EventoType `json:"type"`
	IsPublic     bool       `json:"is_public"`
}

type UpdateEventoRequest struct {
	Title        *string     `json:"title,omitempty"`
	Description  *string     `json:"description,omitempty"`
	EventDate    *time.Time  `json:"event_date,omitempty"`
	EventEndDate *time.Time  `json:"event_end_date,omitempty"`
	Location     *string     `json:"location,omitempty"`
	Type         *EventoType `json:"type,omitempty"`
	IsPublic     *bool       `json:"is_public,omitempty"`
}

type EventoListResponse struct {
	Eventos []Evento `json:"eventos"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}

type EventoFilter struct {
	Type      EventoType
	IsPublic  *bool
	FromDate  *time.Time
	ToDate    *time.Time
	Page      int
	PerPage   int
	Upcoming  bool
}
