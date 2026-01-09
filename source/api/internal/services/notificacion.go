package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/pkg/email"
)

var ErrNotificacionNotFound = errors.New("notificacion not found")

type NotificacionService struct {
	db    *database.DB
	email *email.Service
}

func NewNotificacionService(db *database.DB, emailSvc *email.Service) *NotificacionService {
	return &NotificacionService{db: db, email: emailSvc}
}

func (s *NotificacionService) List(ctx context.Context, filter models.NotificacionFilter) (*models.NotificacionListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT id, user_id, title, body, type, reference_id, is_read, read_at, created_at
		FROM notificaciones
		WHERE user_id = $1`
	countQuery := `SELECT COUNT(*) FROM notificaciones WHERE user_id = $1`
	unreadQuery := `SELECT COUNT(*) FROM notificaciones WHERE user_id = $1 AND is_read = FALSE`
	args := []interface{}{filter.UserID}
	argCount := 1

	if filter.Type != "" {
		argCount++
		query += ` AND type = $` + strconv.Itoa(argCount)
		countQuery += ` AND type = $` + strconv.Itoa(argCount)
		args = append(args, filter.Type)
	}
	if filter.IsRead != nil {
		argCount++
		query += ` AND is_read = $` + strconv.Itoa(argCount)
		countQuery += ` AND is_read = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsRead)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	var unread int
	err = s.db.Pool.QueryRow(ctx, unreadQuery, filter.UserID).Scan(&unread)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY created_at DESC`
	argCount++
	query += fmt.Sprintf(` LIMIT $%d`, argCount)
	args = append(args, filter.PerPage)
	argCount++
	query += fmt.Sprintf(` OFFSET $%d`, argCount)
	args = append(args, offset)

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notificaciones := []models.Notificacion{}
	for rows.Next() {
		var n models.Notificacion
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Type,
			&n.ReferenceID, &n.IsRead, &n.ReadAt, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notificaciones = append(notificaciones, n)
	}

	return &models.NotificacionListResponse{
		Notificaciones: notificaciones,
		Total:          total,
		Unread:         unread,
		Page:           filter.Page,
		PerPage:        filter.PerPage,
	}, nil
}

func (s *NotificacionService) GetByID(ctx context.Context, id string, userID string) (*models.Notificacion, error) {
	var n models.Notificacion
	err := s.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, title, body, type, reference_id, is_read, read_at, created_at
		FROM notificaciones
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&n.ID, &n.UserID, &n.Title, &n.Body, &n.Type,
		&n.ReferenceID, &n.IsRead, &n.ReadAt, &n.CreatedAt)
	if err != nil {
		return nil, ErrNotificacionNotFound
	}
	return &n, nil
}

func (s *NotificacionService) GetStats(ctx context.Context, userID string) (*models.NotificacionStats, error) {
	var stats models.NotificacionStats

	err := s.db.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_read = FALSE) as unread,
			COUNT(*) FILTER (WHERE is_read = TRUE) as read
		FROM notificaciones
		WHERE user_id = $1`, userID).Scan(&stats.Total, &stats.Unread, &stats.Read)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *NotificacionService) Create(ctx context.Context, req *models.CreateNotificacionRequest) (*models.Notificacion, error) {
	var n models.Notificacion
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO notificaciones (user_id, title, body, type, reference_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, body, type, reference_id, is_read, read_at, created_at`,
		req.UserID, req.Title, req.Body, req.Type, req.ReferenceID).Scan(
		&n.ID, &n.UserID, &n.Title, &n.Body, &n.Type,
		&n.ReferenceID, &n.IsRead, &n.ReadAt, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// CreateBulk creates notifications for multiple users
func (s *NotificacionService) CreateBulk(ctx context.Context, req *models.CreateBulkNotificacionRequest) (int, error) {
	if len(req.UserIDs) == 0 {
		return 0, nil
	}

	// Use batch insert
	count := 0
	for _, userID := range req.UserIDs {
		_, err := s.db.Pool.Exec(ctx, `
			INSERT INTO notificaciones (user_id, title, body, type, reference_id)
			VALUES ($1, $2, $3, $4, $5)`,
			userID, req.Title, req.Body, req.Type, req.ReferenceID)
		if err != nil {
			// Skip failed inserts but continue
			continue
		}
		count++
	}

	return count, nil
}

// CreateBroadcast creates notifications for all users matching roles
func (s *NotificacionService) CreateBroadcast(ctx context.Context, req *models.CreateBroadcastNotificacionRequest) (int, error) {
	// Get all user IDs and emails matching roles
	var query string
	var args []interface{}

	if len(req.Roles) == 0 {
		// All users
		query = `SELECT id, email FROM users WHERE role != 'visitor'`
	} else {
		// Specific roles
		query = `SELECT id, email FROM users WHERE role = ANY($1)`
		args = append(args, req.Roles)
	}

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	userIDs := []string{}
	userEmails := []string{}
	for rows.Next() {
		var id, userEmail string
		if err := rows.Scan(&id, &userEmail); err != nil {
			continue
		}
		userIDs = append(userIDs, id)
		if userEmail != "" {
			userEmails = append(userEmails, userEmail)
		}
	}

	if len(userIDs) == 0 {
		return 0, nil
	}

	// Send email for emergency notifications
	if s.email != nil && req.SendEmail && len(userEmails) > 0 {
		go s.sendBroadcastEmails(userEmails, req.Title, req.Body, req.Type)
	}

	return s.CreateBulk(ctx, &models.CreateBulkNotificacionRequest{
		UserIDs:     userIDs,
		Title:       req.Title,
		Body:        req.Body,
		Type:        req.Type,
		ReferenceID: req.ReferenceID,
	})
}

// sendBroadcastEmails sends notification emails to multiple users
func (s *NotificacionService) sendBroadcastEmails(emails []string, title, body string, notifType models.NotificationType) {
	var templateName string
	var subject string

	switch notifType {
	case models.NotificationTypeEmergencia:
		templateName = email.TemplateEmergencia
		subject = "ALERTA - " + title
	default:
		templateName = email.TemplateNotificacion
		subject = title + " - Comunidad Viña Pelvin"
	}

	// Send in batches to avoid overloading SMTP
	for _, userEmail := range emails {
		var data interface{}
		if notifType == models.NotificationTypeEmergencia {
			data = map[string]string{
				"Title":         title,
				"Content":       body,
				"Priority":      "Alta",
				"PriorityClass": "high",
				"URL":           "https://vinapelvin.cl",
			}
		} else {
			data = map[string]string{
				"Title": title,
				"Body":  body,
				"URL":   "https://vinapelvin.cl",
			}
		}

		err := s.email.Send(email.Email{
			To:       []string{userEmail},
			Subject:  subject,
			Template: templateName,
			Data:     data,
		})
		if err != nil {
			log.Printf("[EMAIL] Failed to send broadcast to %s: %v", userEmail, err)
		}
	}
}

func (s *NotificacionService) MarkAsRead(ctx context.Context, id string, userID string) (*models.Notificacion, error) {
	now := time.Now()
	result, err := s.db.Pool.Exec(ctx, `
		UPDATE notificaciones
		SET is_read = TRUE, read_at = $1
		WHERE id = $2 AND user_id = $3 AND is_read = FALSE`,
		now, id, userID)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		// Either not found or already read
		return s.GetByID(ctx, id, userID)
	}

	return s.GetByID(ctx, id, userID)
}

func (s *NotificacionService) MarkAllAsRead(ctx context.Context, userID string) (int, error) {
	now := time.Now()
	result, err := s.db.Pool.Exec(ctx, `
		UPDATE notificaciones
		SET is_read = TRUE, read_at = $1
		WHERE user_id = $2 AND is_read = FALSE`,
		now, userID)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (s *NotificacionService) Delete(ctx context.Context, id string, userID string) error {
	result, err := s.db.Pool.Exec(ctx, `
		DELETE FROM notificaciones
		WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotificacionNotFound
	}
	return nil
}

func (s *NotificacionService) DeleteAll(ctx context.Context, userID string) (int, error) {
	result, err := s.db.Pool.Exec(ctx, `
		DELETE FROM notificaciones
		WHERE user_id = $1`, userID)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (s *NotificacionService) DeleteRead(ctx context.Context, userID string) (int, error) {
	result, err := s.db.Pool.Exec(ctx, `
		DELETE FROM notificaciones
		WHERE user_id = $1 AND is_read = TRUE`, userID)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}
