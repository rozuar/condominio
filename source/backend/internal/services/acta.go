package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var ErrActaNotFound = errors.New("acta not found")

type ActaService struct {
	db *database.DB
}

func NewActaService(db *database.DB) *ActaService {
	return &ActaService{db: db}
}

func (s *ActaService) List(ctx context.Context, filter models.ActaFilter) (*models.ActaListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT a.id, a.title, a.content, a.meeting_date, a.type, a.attendees_count,
		       a.created_by, COALESCE(u.name, ''), a.created_at, a.updated_at
		FROM actas a
		LEFT JOIN users u ON a.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM actas WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Type != "" {
		argCount++
		query += ` AND a.type = $` + strconv.Itoa(argCount)
		countQuery += ` AND type = $` + strconv.Itoa(argCount)
		args = append(args, filter.Type)
	}

	if filter.Year > 0 {
		argCount++
		query += ` AND EXTRACT(YEAR FROM a.meeting_date) = $` + strconv.Itoa(argCount)
		countQuery += ` AND EXTRACT(YEAR FROM meeting_date) = $` + strconv.Itoa(argCount)
		args = append(args, filter.Year)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY a.meeting_date DESC`
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

	actas := []models.Acta{}
	for rows.Next() {
		var a models.Acta
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.MeetingDate, &a.Type, &a.AttendeesCount,
			&a.CreatedBy, &a.CreatorName, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actas = append(actas, a)
	}

	return &models.ActaListResponse{
		Actas:   actas,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}, nil
}

func (s *ActaService) GetByID(ctx context.Context, id string) (*models.Acta, error) {
	var a models.Acta
	err := s.db.Pool.QueryRow(ctx, `
		SELECT a.id, a.title, a.content, a.meeting_date, a.type, a.attendees_count,
		       a.created_by, COALESCE(u.name, ''), a.created_at, a.updated_at
		FROM actas a
		LEFT JOIN users u ON a.created_by = u.id
		WHERE a.id = $1`, id).Scan(
		&a.ID, &a.Title, &a.Content, &a.MeetingDate, &a.Type, &a.AttendeesCount,
		&a.CreatedBy, &a.CreatorName, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, ErrActaNotFound
	}
	return &a, nil
}

func (s *ActaService) Create(ctx context.Context, req *models.CreateActaRequest, createdBy string) (*models.Acta, error) {
	var a models.Acta
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO actas (title, content, meeting_date, type, attendees_count, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, content, meeting_date, type, attendees_count, created_by, created_at, updated_at`,
		req.Title, req.Content, req.MeetingDate, req.Type, req.AttendeesCount, createdBy).Scan(
		&a.ID, &a.Title, &a.Content, &a.MeetingDate, &a.Type, &a.AttendeesCount, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
