package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var (
	ErrEmergenciaNotFound = errors.New("emergencia not found")
)

type EmergenciaService struct {
	db *database.DB
}

func NewEmergenciaService(db *database.DB) *EmergenciaService {
	return &EmergenciaService{db: db}
}

func (s *EmergenciaService) List(ctx context.Context, filter models.EmergenciaFilter) (*models.EmergenciaListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT e.id, e.title, e.content, e.priority, e.status, e.expires_at,
		       e.notify_email, e.notify_push, e.created_by,
		       COALESCE(u.name, '') as creator_name,
		       e.resolved_at, e.resolved_by,
		       COALESCE(r.name, '') as resolver_name,
		       e.created_at, e.updated_at
		FROM emergencias e
		LEFT JOIN users u ON e.created_by = u.id
		LEFT JOIN users r ON e.resolved_by = r.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM emergencias WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(` AND e.status = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND status = $%d`, argCount)
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		argCount++
		query += fmt.Sprintf(` AND e.priority = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND priority = $%d`, argCount)
		args = append(args, filter.Priority)
	}

	if filter.Active {
		query += ` AND e.status = 'active'`
		countQuery += ` AND status = 'active'`
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Order by priority (critical first) then by created_at
	query += ` ORDER BY
		CASE e.priority
			WHEN 'critical' THEN 1
			WHEN 'high' THEN 2
			WHEN 'medium' THEN 3
			WHEN 'low' THEN 4
		END,
		e.created_at DESC`

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

	emergencias := []models.Emergencia{}
	for rows.Next() {
		var e models.Emergencia
		err := rows.Scan(
			&e.ID, &e.Title, &e.Content, &e.Priority, &e.Status, &e.ExpiresAt,
			&e.NotifyEmail, &e.NotifyPush, &e.CreatedBy, &e.CreatorName,
			&e.ResolvedAt, &e.ResolvedBy, &e.ResolverName,
			&e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		emergencias = append(emergencias, e)
	}

	return &models.EmergenciaListResponse{
		Emergencias: emergencias,
		Total:       total,
		Page:        filter.Page,
		PerPage:     filter.PerPage,
	}, nil
}

func (s *EmergenciaService) GetByID(ctx context.Context, id string) (*models.Emergencia, error) {
	var e models.Emergencia
	err := s.db.Pool.QueryRow(ctx, `
		SELECT e.id, e.title, e.content, e.priority, e.status, e.expires_at,
		       e.notify_email, e.notify_push, e.created_by,
		       COALESCE(u.name, '') as creator_name,
		       e.resolved_at, e.resolved_by,
		       COALESCE(r.name, '') as resolver_name,
		       e.created_at, e.updated_at
		FROM emergencias e
		LEFT JOIN users u ON e.created_by = u.id
		LEFT JOIN users r ON e.resolved_by = r.id
		WHERE e.id = $1`, id).Scan(
		&e.ID, &e.Title, &e.Content, &e.Priority, &e.Status, &e.ExpiresAt,
		&e.NotifyEmail, &e.NotifyPush, &e.CreatedBy, &e.CreatorName,
		&e.ResolvedAt, &e.ResolvedBy, &e.ResolverName,
		&e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, ErrEmergenciaNotFound
	}
	return &e, nil
}

func (s *EmergenciaService) GetActive(ctx context.Context) ([]models.Emergencia, error) {
	rows, err := s.db.Pool.Query(ctx, `
		SELECT e.id, e.title, e.content, e.priority, e.status, e.expires_at,
		       e.notify_email, e.notify_push, e.created_by,
		       COALESCE(u.name, '') as creator_name,
		       e.resolved_at, e.resolved_by,
		       COALESCE(r.name, '') as resolver_name,
		       e.created_at, e.updated_at
		FROM emergencias e
		LEFT JOIN users u ON e.created_by = u.id
		LEFT JOIN users r ON e.resolved_by = r.id
		WHERE e.status = 'active'
		  AND (e.expires_at IS NULL OR e.expires_at > NOW())
		ORDER BY
			CASE e.priority
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				WHEN 'low' THEN 4
			END,
			e.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	emergencias := []models.Emergencia{}
	for rows.Next() {
		var e models.Emergencia
		err := rows.Scan(
			&e.ID, &e.Title, &e.Content, &e.Priority, &e.Status, &e.ExpiresAt,
			&e.NotifyEmail, &e.NotifyPush, &e.CreatedBy, &e.CreatorName,
			&e.ResolvedAt, &e.ResolvedBy, &e.ResolverName,
			&e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		emergencias = append(emergencias, e)
	}
	return emergencias, nil
}

func (s *EmergenciaService) Create(ctx context.Context, req *models.CreateEmergenciaRequest, createdBy string) (*models.Emergencia, error) {
	if !req.Priority.IsValid() {
		req.Priority = models.EmergenciaPriorityMedium
	}

	var e models.Emergencia
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO emergencias (title, content, priority, status, expires_at, notify_email, notify_push, created_by)
		VALUES ($1, $2, $3, 'active', $4, $5, $6, $7)
		RETURNING id, title, content, priority, status, expires_at, notify_email, notify_push, created_by, created_at, updated_at`,
		req.Title, req.Content, req.Priority, req.ExpiresAt, req.NotifyEmail, req.NotifyPush, createdBy).Scan(
		&e.ID, &e.Title, &e.Content, &e.Priority, &e.Status, &e.ExpiresAt,
		&e.NotifyEmail, &e.NotifyPush, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// TODO: Send notifications if notify_email or notify_push is true

	return &e, nil
}

func (s *EmergenciaService) Update(ctx context.Context, id string, req *models.UpdateEmergenciaRequest) (*models.Emergencia, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Content != nil {
		current.Content = *req.Content
	}
	if req.Priority != nil {
		current.Priority = *req.Priority
	}
	if req.Status != nil {
		current.Status = *req.Status
	}
	if req.ExpiresAt != nil {
		current.ExpiresAt = req.ExpiresAt
	}
	if req.NotifyEmail != nil {
		current.NotifyEmail = *req.NotifyEmail
	}
	if req.NotifyPush != nil {
		current.NotifyPush = *req.NotifyPush
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE emergencias
		SET title = $1, content = $2, priority = $3, status = $4, expires_at = $5,
		    notify_email = $6, notify_push = $7, updated_at = NOW()
		WHERE id = $8`,
		current.Title, current.Content, current.Priority, current.Status, current.ExpiresAt,
		current.NotifyEmail, current.NotifyPush, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *EmergenciaService) Resolve(ctx context.Context, id string, resolvedBy string) (*models.Emergencia, error) {
	now := time.Now()
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE emergencias
		SET status = 'resolved', resolved_at = $1, resolved_by = $2, updated_at = NOW()
		WHERE id = $3 AND status = 'active'`,
		now, resolvedBy, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *EmergenciaService) Delete(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM emergencias WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrEmergenciaNotFound
	}
	return nil
}

// ExpireOld marks expired emergencies as expired
func (s *EmergenciaService) ExpireOld(ctx context.Context) (int, error) {
	result, err := s.db.Pool.Exec(ctx, `
		UPDATE emergencias
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'active' AND expires_at IS NOT NULL AND expires_at < NOW()`)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}
