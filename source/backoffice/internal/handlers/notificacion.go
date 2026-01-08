package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
)

type NotificacionHandler struct {
	service *services.NotificacionService
}

func NewNotificacionHandler(service *services.NotificacionService) *NotificacionHandler {
	return &NotificacionHandler{service: service}
}

// List returns user's notifications
func (h *NotificacionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	filter := models.NotificacionFilter{
		UserID:  userID,
		Page:    1,
		PerPage: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}
	if perPage := r.URL.Query().Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil {
			filter.PerPage = pp
		}
	}
	if notifType := r.URL.Query().Get("type"); notifType != "" {
		filter.Type = models.NotificationType(notifType)
	}
	if isRead := r.URL.Query().Get("is_read"); isRead != "" {
		b := isRead == "true"
		filter.IsRead = &b
	}

	resp, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list notifications")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetByID returns a specific notification
func (h *NotificacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	id := chi.URLParam(r, "id")

	notif, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotificacionNotFound) {
			writeError(w, http.StatusNotFound, "Notification not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}

	writeJSON(w, http.StatusOK, notif)
}

// GetStats returns notification statistics for the user
func (h *NotificacionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	stats, err := h.service.GetStats(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// MarkAsRead marks a notification as read
func (h *NotificacionHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	id := chi.URLParam(r, "id")

	notif, err := h.service.MarkAsRead(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotificacionNotFound) {
			writeError(w, http.StatusNotFound, "Notification not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to mark as read")
		return
	}

	writeJSON(w, http.StatusOK, notif)
}

// MarkAllAsRead marks all notifications as read
func (h *NotificacionHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	count, err := h.service.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to mark all as read")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Notifications marked as read",
		"count":   count,
	})
}

// Delete removes a notification
func (h *NotificacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotificacionNotFound) {
			writeError(w, http.StatusNotFound, "Notification not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Notification deleted"})
}

// DeleteAll removes all notifications for the user
func (h *NotificacionHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	count, err := h.service.DeleteAll(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete notifications")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "All notifications deleted",
		"count":   count,
	})
}

// DeleteRead removes all read notifications for the user
func (h *NotificacionHandler) DeleteRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	count, err := h.service.DeleteRead(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete read notifications")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Read notifications deleted",
		"count":   count,
	})
}

// Admin endpoints

// Create creates a notification for a specific user (admin only)
func (h *NotificacionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "User ID is required")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.Body == "" {
		writeError(w, http.StatusBadRequest, "Body is required")
		return
	}
	if !req.Type.IsValid() {
		req.Type = models.NotificationTypeSistema
	}

	notif, err := h.service.Create(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create notification")
		return
	}

	writeJSON(w, http.StatusCreated, notif)
}

// CreateBulk creates notifications for multiple users (admin only)
func (h *NotificacionHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBulkNotificacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.UserIDs) == 0 {
		writeError(w, http.StatusBadRequest, "User IDs are required")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.Body == "" {
		writeError(w, http.StatusBadRequest, "Body is required")
		return
	}
	if !req.Type.IsValid() {
		req.Type = models.NotificationTypeSistema
	}

	count, err := h.service.CreateBulk(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create notifications")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Notifications created",
		"count":   count,
	})
}

// CreateBroadcast creates notifications for users with specific roles (admin only)
func (h *NotificacionHandler) CreateBroadcast(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBroadcastNotificacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.Body == "" {
		writeError(w, http.StatusBadRequest, "Body is required")
		return
	}
	if !req.Type.IsValid() {
		req.Type = models.NotificationTypeSistema
	}

	count, err := h.service.CreateBroadcast(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to broadcast notifications")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Broadcast notifications created",
		"count":   count,
	})
}
