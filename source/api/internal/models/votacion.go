package models

import "time"

type VotacionStatus string

const (
	VotacionStatusDraft     VotacionStatus = "draft"
	VotacionStatusActive    VotacionStatus = "active"
	VotacionStatusClosed    VotacionStatus = "closed"
	VotacionStatusCancelled VotacionStatus = "cancelled"
)

type Votacion struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description,omitempty"`
	Status           VotacionStatus `json:"status"`
	StartDate        *time.Time     `json:"start_date,omitempty"`
	EndDate          *time.Time     `json:"end_date,omitempty"`
	RequiresQuorum   bool           `json:"requires_quorum"`
	QuorumPercentage int            `json:"quorum_percentage"`
	AllowAbstention  bool           `json:"allow_abstention"`
	Opciones         []VotacionOpcion `json:"opciones,omitempty"`
	CreatedBy        *string        `json:"created_by,omitempty"`
	CreatorName      string         `json:"creator_name,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	// Computed fields
	TotalVotos       int            `json:"total_votos,omitempty"`
	HasVoted         bool           `json:"has_voted,omitempty"`
}

type VotacionOpcion struct {
	ID          string `json:"id"`
	VotacionID  string `json:"votacion_id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	OrderIndex  int    `json:"order_index"`
	VotosCount  int    `json:"votos_count,omitempty"`
}

type Voto struct {
	ID          string    `json:"id"`
	VotacionID  string    `json:"votacion_id"`
	UserID      string    `json:"user_id"`
	OpcionID    *string   `json:"opcion_id,omitempty"`
	IsAbstention bool     `json:"is_abstention"`
	VotedAt     time.Time `json:"voted_at"`
}

type CreateVotacionRequest struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	RequiresQuorum   bool     `json:"requires_quorum"`
	QuorumPercentage int      `json:"quorum_percentage"`
	AllowAbstention  bool     `json:"allow_abstention"`
	Opciones         []string `json:"opciones"` // Labels for options
}

type UpdateVotacionRequest struct {
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	RequiresQuorum   *bool   `json:"requires_quorum,omitempty"`
	QuorumPercentage *int    `json:"quorum_percentage,omitempty"`
	AllowAbstention  *bool   `json:"allow_abstention,omitempty"`
}

type AddOpcionRequest struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

type EmitirVotoRequest struct {
	OpcionID     *string `json:"opcion_id,omitempty"`
	IsAbstention bool    `json:"is_abstention"`
}

type VotacionListResponse struct {
	Votaciones []Votacion `json:"votaciones"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
}

type VotacionResultado struct {
	Votacion         Votacion           `json:"votacion"`
	TotalVotos       int                `json:"total_votos"`
	TotalAbstenciones int               `json:"total_abstenciones"`
	Resultados       []OpcionResultado  `json:"resultados"`
	QuorumAlcanzado  bool               `json:"quorum_alcanzado"`
	TotalVecinos     int                `json:"total_vecinos"`
	Participacion    float64            `json:"participacion"`
}

type OpcionResultado struct {
	OpcionID   string  `json:"opcion_id"`
	Label      string  `json:"label"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type VotacionFilter struct {
	Status  VotacionStatus
	Active  bool
	Page    int
	PerPage int
}

func (s VotacionStatus) IsValid() bool {
	switch s {
	case VotacionStatusDraft, VotacionStatusActive, VotacionStatusClosed, VotacionStatusCancelled:
		return true
	}
	return false
}

func (v *Votacion) CanVote() bool {
	if v.Status != VotacionStatusActive {
		return false
	}
	now := time.Now()
	if v.StartDate != nil && now.Before(*v.StartDate) {
		return false
	}
	if v.EndDate != nil && now.After(*v.EndDate) {
		return false
	}
	return true
}

func (v *Votacion) IsActive() bool {
	return v.Status == VotacionStatusActive
}
