package models

import "time"

type EmergenciaPriority string

const (
	EmergenciaPriorityLow      EmergenciaPriority = "low"
	EmergenciaPriorityMedium   EmergenciaPriority = "medium"
	EmergenciaPriorityHigh     EmergenciaPriority = "high"
	EmergenciaPriorityCritical EmergenciaPriority = "critical"
)

type EmergenciaStatus string

const (
	EmergenciaStatusActive   EmergenciaStatus = "active"
	EmergenciaStatusResolved EmergenciaStatus = "resolved"
	EmergenciaStatusExpired  EmergenciaStatus = "expired"
)

type Emergencia struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Content      string             `json:"content"`
	Priority     EmergenciaPriority `json:"priority"`
	Status       EmergenciaStatus   `json:"status"`
	ExpiresAt    *time.Time         `json:"expires_at,omitempty"`
	NotifyEmail  bool               `json:"notify_email"`
	NotifyPush   bool               `json:"notify_push"`
	CreatedBy    *string            `json:"created_by,omitempty"`
	CreatorName  string             `json:"creator_name,omitempty"`
	ResolvedAt   *time.Time         `json:"resolved_at,omitempty"`
	ResolvedBy   *string            `json:"resolved_by,omitempty"`
	ResolverName string             `json:"resolver_name,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type CreateEmergenciaRequest struct {
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Priority    EmergenciaPriority `json:"priority"`
	ExpiresAt   *time.Time         `json:"expires_at,omitempty"`
	NotifyEmail bool               `json:"notify_email"`
	NotifyPush  bool               `json:"notify_push"`
}

type UpdateEmergenciaRequest struct {
	Title       *string             `json:"title,omitempty"`
	Content     *string             `json:"content,omitempty"`
	Priority    *EmergenciaPriority `json:"priority,omitempty"`
	Status      *EmergenciaStatus   `json:"status,omitempty"`
	ExpiresAt   *time.Time          `json:"expires_at,omitempty"`
	NotifyEmail *bool               `json:"notify_email,omitempty"`
	NotifyPush  *bool               `json:"notify_push,omitempty"`
}

type ResolveEmergenciaRequest struct {
	Resolution string `json:"resolution,omitempty"`
}

type EmergenciaListResponse struct {
	Emergencias []Emergencia `json:"emergencias"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
}

type EmergenciaFilter struct {
	Status   EmergenciaStatus
	Priority EmergenciaPriority
	Active   bool
	Page     int
	PerPage  int
}

func (p EmergenciaPriority) IsValid() bool {
	switch p {
	case EmergenciaPriorityLow, EmergenciaPriorityMedium, EmergenciaPriorityHigh, EmergenciaPriorityCritical:
		return true
	}
	return false
}

func (s EmergenciaStatus) IsValid() bool {
	switch s {
	case EmergenciaStatusActive, EmergenciaStatusResolved, EmergenciaStatusExpired:
		return true
	}
	return false
}

func (p EmergenciaPriority) Order() int {
	switch p {
	case EmergenciaPriorityCritical:
		return 4
	case EmergenciaPriorityHigh:
		return 3
	case EmergenciaPriorityMedium:
		return 2
	case EmergenciaPriorityLow:
		return 1
	}
	return 0
}
