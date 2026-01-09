package services

import (
	"context"
	"errors"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var (
	ErrEventoNotFound = errors.New("evento not found")
)

type EventoService struct {
	db *database.DB
}

func NewEventoService(db *database.DB) *EventoService {
	return &EventoService{db: db}
}

func (s *EventoService) List(ctx context.Context, filter models.EventoFilter) (*models.EventoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT e.id, e.title, COALESCE(e.description, ''), e.event_date, e.event_end_date,
		       COALESCE(e.location, ''), e.type, e.is_public, e.created_by,
		       COALESCE(u.name, '') as creator_name, e.created_at, e.updated_at
		FROM eventos e
		LEFT JOIN users u ON e.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM eventos WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.IsPublic != nil {
		argCount++
		query += ` AND e.is_public = $` + string(rune('0'+argCount))
		countQuery += ` AND is_public = $` + string(rune('0'+argCount))
		args = append(args, *filter.IsPublic)
	}

	if filter.Type != "" {
		argCount++
		query += ` AND e.type = $` + string(rune('0'+argCount))
		countQuery += ` AND type = $` + string(rune('0'+argCount))
		args = append(args, filter.Type)
	}

	if filter.Upcoming {
		argCount++
		query += ` AND e.event_date >= $` + string(rune('0'+argCount))
		countQuery += ` AND event_date >= $` + string(rune('0'+argCount))
		args = append(args, time.Now())
	}

	if filter.FromDate != nil {
		argCount++
		query += ` AND e.event_date >= $` + string(rune('0'+argCount))
		countQuery += ` AND event_date >= $` + string(rune('0'+argCount))
		args = append(args, *filter.FromDate)
	}

	if filter.ToDate != nil {
		argCount++
		query += ` AND e.event_date <= $` + string(rune('0'+argCount))
		countQuery += ` AND event_date <= $` + string(rune('0'+argCount))
		args = append(args, *filter.ToDate)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY e.event_date ASC`
	argCount++
	query += ` LIMIT $` + string(rune('0'+argCount))
	args = append(args, filter.PerPage)
	argCount++
	query += ` OFFSET $` + string(rune('0'+argCount))
	args = append(args, offset)

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventos := []models.Evento{}
	for rows.Next() {
		var e models.Evento
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate, &e.EventEndDate,
			&e.Location, &e.Type, &e.IsPublic, &e.CreatedBy, &e.CreatorName, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		eventos = append(eventos, e)
	}

	return &models.EventoListResponse{
		Eventos: eventos,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}, nil
}

func (s *EventoService) GetByID(ctx context.Context, id string) (*models.Evento, error) {
	var e models.Evento
	err := s.db.Pool.QueryRow(ctx, `
		SELECT e.id, e.title, COALESCE(e.description, ''), e.event_date, e.event_end_date,
		       COALESCE(e.location, ''), e.type, e.is_public, e.created_by,
		       COALESCE(u.name, '') as creator_name, e.created_at, e.updated_at
		FROM eventos e
		LEFT JOIN users u ON e.created_by = u.id
		WHERE e.id = $1`, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventDate, &e.EventEndDate,
		&e.Location, &e.Type, &e.IsPublic, &e.CreatedBy, &e.CreatorName, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, ErrEventoNotFound
	}
	return &e, nil
}

func (s *EventoService) Create(ctx context.Context, req *models.CreateEventoRequest, createdBy string) (*models.Evento, error) {
	var e models.Evento
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO eventos (title, description, event_date, event_end_date, location, type, is_public, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, event_date, event_end_date, location, type, is_public, created_by, created_at, updated_at`,
		req.Title, req.Description, req.EventDate, req.EventEndDate, req.Location, req.Type, req.IsPublic, createdBy).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventDate, &e.EventEndDate, &e.Location, &e.Type, &e.IsPublic, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventoService) Update(ctx context.Context, id string, req *models.UpdateEventoRequest) (*models.Evento, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.EventDate != nil {
		current.EventDate = *req.EventDate
	}
	if req.EventEndDate != nil {
		current.EventEndDate = req.EventEndDate
	}
	if req.Location != nil {
		current.Location = *req.Location
	}
	if req.Type != nil {
		current.Type = *req.Type
	}
	if req.IsPublic != nil {
		current.IsPublic = *req.IsPublic
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE eventos
		SET title = $1, description = $2, event_date = $3, event_end_date = $4,
		    location = $5, type = $6, is_public = $7, updated_at = NOW()
		WHERE id = $8`,
		current.Title, current.Description, current.EventDate, current.EventEndDate,
		current.Location, current.Type, current.IsPublic, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *EventoService) Delete(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM eventos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrEventoNotFound
	}
	return nil
}

func (s *EventoService) GetUpcoming(ctx context.Context, limit int, publicOnly bool) ([]models.Evento, error) {
	query := `
		SELECT e.id, e.title, COALESCE(e.description, ''), e.event_date, e.event_end_date,
		       COALESCE(e.location, ''), e.type, e.is_public, e.created_by,
		       COALESCE(u.name, '') as creator_name, e.created_at, e.updated_at
		FROM eventos e
		LEFT JOIN users u ON e.created_by = u.id
		WHERE e.event_date >= $1`

	args := []interface{}{time.Now()}
	if publicOnly {
		query += ` AND e.is_public = true`
	}
	query += ` ORDER BY e.event_date ASC LIMIT $2`
	args = append(args, limit)

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventos := []models.Evento{}
	for rows.Next() {
		var e models.Evento
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate, &e.EventEndDate,
			&e.Location, &e.Type, &e.IsPublic, &e.CreatedBy, &e.CreatorName, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		eventos = append(eventos, e)
	}
	return eventos, nil
}
