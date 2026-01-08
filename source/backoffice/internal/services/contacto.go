package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/pkg/email"
)

var (
	ErrMensajeNotFound = errors.New("mensaje not found")
	ErrAlreadyReplied  = errors.New("mensaje already replied")
)

type ContactoService struct {
	db    *database.DB
	email *email.Service
}

func NewContactoService(db *database.DB, emailSvc *email.Service) *ContactoService {
	return &ContactoService{db: db, email: emailSvc}
}

func (s *ContactoService) List(ctx context.Context, filter models.MensajeContactoFilter) (*models.MensajeContactoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT m.id, m.user_id, COALESCE(u.name, '') as user_name,
		       m.nombre, m.email, m.asunto, m.mensaje, m.status,
		       m.read_at, m.read_by, COALESCE(ur.name, '') as read_by_name,
		       m.replied_at, m.replied_by, COALESCE(urp.name, '') as replied_by_name,
		       COALESCE(m.respuesta, ''),
		       m.created_at, m.updated_at
		FROM mensajes_contacto m
		LEFT JOIN users u ON m.user_id = u.id
		LEFT JOIN users ur ON m.read_by = ur.id
		LEFT JOIN users urp ON m.replied_by = urp.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM mensajes_contacto WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(` AND m.status = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND status = $%d`, argCount)
		args = append(args, filter.Status)
	}

	if filter.UserID != "" {
		argCount++
		query += fmt.Sprintf(` AND m.user_id = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND user_id = $%d`, argCount)
		args = append(args, filter.UserID)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY
		CASE m.status
			WHEN 'pending' THEN 1
			WHEN 'read' THEN 2
			WHEN 'replied' THEN 3
			WHEN 'archived' THEN 4
		END,
		m.created_at DESC`

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

	mensajes := []models.MensajeContacto{}
	for rows.Next() {
		var m models.MensajeContacto
		err := rows.Scan(
			&m.ID, &m.UserID, &m.UserName,
			&m.Nombre, &m.Email, &m.Asunto, &m.Mensaje, &m.Status,
			&m.ReadAt, &m.ReadBy, &m.ReadByName,
			&m.RepliedAt, &m.RepliedBy, &m.RepliedByName,
			&m.Respuesta,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		mensajes = append(mensajes, m)
	}

	return &models.MensajeContactoListResponse{
		Mensajes: mensajes,
		Total:    total,
		Page:     filter.Page,
		PerPage:  filter.PerPage,
	}, nil
}

func (s *ContactoService) GetByID(ctx context.Context, id string) (*models.MensajeContacto, error) {
	var m models.MensajeContacto
	err := s.db.Pool.QueryRow(ctx, `
		SELECT m.id, m.user_id, COALESCE(u.name, '') as user_name,
		       m.nombre, m.email, m.asunto, m.mensaje, m.status,
		       m.read_at, m.read_by, COALESCE(ur.name, '') as read_by_name,
		       m.replied_at, m.replied_by, COALESCE(urp.name, '') as replied_by_name,
		       COALESCE(m.respuesta, ''),
		       m.created_at, m.updated_at
		FROM mensajes_contacto m
		LEFT JOIN users u ON m.user_id = u.id
		LEFT JOIN users ur ON m.read_by = ur.id
		LEFT JOIN users urp ON m.replied_by = urp.id
		WHERE m.id = $1`, id).Scan(
		&m.ID, &m.UserID, &m.UserName,
		&m.Nombre, &m.Email, &m.Asunto, &m.Mensaje, &m.Status,
		&m.ReadAt, &m.ReadBy, &m.ReadByName,
		&m.RepliedAt, &m.RepliedBy, &m.RepliedByName,
		&m.Respuesta,
		&m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, ErrMensajeNotFound
	}
	return &m, nil
}

func (s *ContactoService) GetMisMensajes(ctx context.Context, userID string) ([]models.MensajeContacto, error) {
	rows, err := s.db.Pool.Query(ctx, `
		SELECT m.id, m.user_id, '',
		       m.nombre, m.email, m.asunto, m.mensaje, m.status,
		       m.read_at, m.read_by, '',
		       m.replied_at, m.replied_by, '',
		       COALESCE(m.respuesta, ''),
		       m.created_at, m.updated_at
		FROM mensajes_contacto m
		WHERE m.user_id = $1
		ORDER BY m.created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mensajes := []models.MensajeContacto{}
	for rows.Next() {
		var m models.MensajeContacto
		err := rows.Scan(
			&m.ID, &m.UserID, &m.UserName,
			&m.Nombre, &m.Email, &m.Asunto, &m.Mensaje, &m.Status,
			&m.ReadAt, &m.ReadBy, &m.ReadByName,
			&m.RepliedAt, &m.RepliedBy, &m.RepliedByName,
			&m.Respuesta,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		mensajes = append(mensajes, m)
	}
	return mensajes, nil
}

func (s *ContactoService) Create(ctx context.Context, req *models.CreateMensajeRequest, userID *string) (*models.MensajeContacto, error) {
	var id string
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO mensajes_contacto (user_id, nombre, email, asunto, mensaje)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		userID, req.Nombre, req.Email, req.Asunto, req.Mensaje).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Send confirmation email
	if s.email != nil && req.Email != "" {
		go func() {
			err := s.email.Send(email.Email{
				To:       []string{req.Email},
				Subject:  "Mensaje Recibido - Comunidad Viña Pelvin",
				Template: email.TemplateContactoRecibido,
				Data: map[string]string{
					"Nombre":  req.Nombre,
					"Asunto":  req.Asunto,
					"Mensaje": req.Mensaje,
				},
			})
			if err != nil {
				log.Printf("[EMAIL] Failed to send contacto confirmation: %v", err)
			}
		}()
	}

	return s.GetByID(ctx, id)
}

func (s *ContactoService) MarkAsRead(ctx context.Context, id string, readBy string) (*models.MensajeContacto, error) {
	msg, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if msg.Status != models.ContactoStatusPending {
		return msg, nil // Already read or replied
	}

	now := time.Now()
	_, err = s.db.Pool.Exec(ctx, `
		UPDATE mensajes_contacto
		SET status = 'read', read_at = $1, read_by = $2, updated_at = NOW()
		WHERE id = $3`,
		now, readBy, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *ContactoService) Reply(ctx context.Context, id string, req *models.ReplyMensajeRequest, repliedBy string) (*models.MensajeContacto, error) {
	msg, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if msg.Status == models.ContactoStatusReplied {
		return nil, ErrAlreadyReplied
	}

	now := time.Now()
	_, err = s.db.Pool.Exec(ctx, `
		UPDATE mensajes_contacto
		SET status = 'replied', replied_at = $1, replied_by = $2, respuesta = $3, updated_at = NOW()
		WHERE id = $4`,
		now, repliedBy, req.Respuesta, id)
	if err != nil {
		return nil, err
	}

	// Send response email
	if s.email != nil && msg.Email != "" {
		go func() {
			err := s.email.Send(email.Email{
				To:       []string{msg.Email},
				Subject:  "Respuesta a tu mensaje - Comunidad Viña Pelvin",
				Template: email.TemplateContactoRespuesta,
				Data: map[string]string{
					"Nombre":          msg.Nombre,
					"Asunto":          msg.Asunto,
					"MensajeOriginal": msg.Mensaje,
					"Respuesta":       req.Respuesta,
				},
			})
			if err != nil {
				log.Printf("[EMAIL] Failed to send contacto response: %v", err)
			}
		}()
	}

	return s.GetByID(ctx, id)
}

func (s *ContactoService) Archive(ctx context.Context, id string) (*models.MensajeContacto, error) {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE mensajes_contacto
		SET status = 'archived', updated_at = NOW()
		WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *ContactoService) Delete(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM mensajes_contacto WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrMensajeNotFound
	}
	return nil
}

func (s *ContactoService) GetStats(ctx context.Context) (*models.ContactoStats, error) {
	var stats models.ContactoStats
	err := s.db.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'read') as read,
			COUNT(*) FILTER (WHERE status = 'replied') as replied,
			COUNT(*) FILTER (WHERE status = 'archived') as archived,
			COUNT(*) as total
		FROM mensajes_contacto`).Scan(
		&stats.TotalPending, &stats.TotalRead, &stats.TotalReplied,
		&stats.TotalArchived, &stats.Total)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
