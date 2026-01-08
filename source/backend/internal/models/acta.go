package models

import "time"

type ActaType string

const (
	ActaOrdinaria      ActaType = "ordinaria"
	ActaExtraordinaria ActaType = "extraordinaria"
)

type Acta struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	MeetingDate    time.Time `json:"meeting_date"`
	Type           ActaType  `json:"type"`
	AttendeesCount *int      `json:"attendees_count,omitempty"`
	CreatedBy      *string   `json:"created_by,omitempty"`
	CreatorName    string    `json:"creator_name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateActaRequest struct {
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	MeetingDate    time.Time `json:"meeting_date"`
	Type           ActaType  `json:"type"`
	AttendeesCount *int      `json:"attendees_count"`
}

type ActaListResponse struct {
	Actas   []Acta `json:"actas"`
	Total   int    `json:"total"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type ActaFilter struct {
	Type    ActaType
	Year    int
	Page    int
	PerPage int
}
