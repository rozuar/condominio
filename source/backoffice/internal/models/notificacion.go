package models

import "time"

type NotificationType string

const (
	NotificationTypeComunicado  NotificationType = "comunicado"
	NotificationTypeEmergencia  NotificationType = "emergencia"
	NotificationTypeVotacion    NotificationType = "votacion"
	NotificationTypePago        NotificationType = "pago"
	NotificationTypeEvento      NotificationType = "evento"
	NotificationTypeDocumento   NotificationType = "documento"
	NotificationTypeActa        NotificationType = "acta"
	NotificationTypeContacto    NotificationType = "contacto"
	NotificationTypeGastoComun  NotificationType = "gasto_comun"
	NotificationTypeSistema     NotificationType = "sistema"
)

type Notificacion struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	Title       string           `json:"title"`
	Body        string           `json:"body"`
	Type        NotificationType `json:"type"`
	ReferenceID *string          `json:"reference_id,omitempty"`
	IsRead      bool             `json:"is_read"`
	ReadAt      *time.Time       `json:"read_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type CreateNotificacionRequest struct {
	UserID      string           `json:"user_id"`
	Title       string           `json:"title"`
	Body        string           `json:"body"`
	Type        NotificationType `json:"type"`
	ReferenceID *string          `json:"reference_id,omitempty"`
}

// CreateBulkNotificacionRequest is for sending notifications to multiple users
type CreateBulkNotificacionRequest struct {
	UserIDs     []string         `json:"user_ids"`
	Title       string           `json:"title"`
	Body        string           `json:"body"`
	Type        NotificationType `json:"type"`
	ReferenceID *string          `json:"reference_id,omitempty"`
}

// CreateBroadcastNotificacionRequest is for sending to all users with specific roles
type CreateBroadcastNotificacionRequest struct {
	Roles       []string         `json:"roles"` // empty = all users
	Title       string           `json:"title"`
	Body        string           `json:"body"`
	Type        NotificationType `json:"type"`
	ReferenceID *string          `json:"reference_id,omitempty"`
	SendEmail   bool             `json:"send_email"` // also send email notification
}

type NotificacionListResponse struct {
	Notificaciones []Notificacion `json:"notificaciones"`
	Total          int            `json:"total"`
	Unread         int            `json:"unread"`
	Page           int            `json:"page"`
	PerPage        int            `json:"per_page"`
}

type NotificacionFilter struct {
	UserID   string
	Type     NotificationType
	IsRead   *bool
	Page     int
	PerPage  int
}

type NotificacionStats struct {
	Total    int `json:"total"`
	Unread   int `json:"unread"`
	Read     int `json:"read"`
}

func (nt NotificationType) IsValid() bool {
	switch nt {
	case NotificationTypeComunicado, NotificationTypeEmergencia, NotificationTypeVotacion,
		NotificationTypePago, NotificationTypeEvento, NotificationTypeDocumento,
		NotificationTypeActa, NotificationTypeContacto, NotificationTypeGastoComun,
		NotificationTypeSistema:
		return true
	}
	return false
}
