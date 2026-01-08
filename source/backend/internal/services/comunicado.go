package services

import (
	"context"
	"errors"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var (
	ErrComunicadoNotFound = errors.New("comunicado not found")
)

type ComunicadoService struct {
	db *database.DB
}

func NewComunicadoService(db *database.DB) *ComunicadoService {
	return &ComunicadoService{db: db}
}

func (s *ComunicadoService) List(ctx context.Context, filter models.ComunicadoFilter) (*models.ComunicadoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	// Build query
	query := `
		SELECT c.id, c.title, c.content, c.type, c.is_public, c.author_id,
		       COALESCE(u.name, '') as author_name, c.published_at, c.created_at, c.updated_at
		FROM comunicados c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM comunicados WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.IsPublic != nil {
		argCount++
		query += ` AND c.is_public = $` + string(rune('0'+argCount))
		countQuery += ` AND is_public = $` + string(rune('0'+argCount))
		args = append(args, *filter.IsPublic)
	}

	if filter.Type != "" {
		argCount++
		query += ` AND c.type = $` + string(rune('0'+argCount))
		countQuery += ` AND type = $` + string(rune('0'+argCount))
		args = append(args, filter.Type)
	}

	// Get total count
	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Add ordering and pagination
	query += ` ORDER BY c.published_at DESC NULLS LAST, c.created_at DESC`
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

	comunicados := []models.Comunicado{}
	for rows.Next() {
		var c models.Comunicado
		err := rows.Scan(&c.ID, &c.Title, &c.Content, &c.Type, &c.IsPublic,
			&c.AuthorID, &c.AuthorName, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comunicados = append(comunicados, c)
	}

	return &models.ComunicadoListResponse{
		Comunicados: comunicados,
		Total:       total,
		Page:        filter.Page,
		PerPage:     filter.PerPage,
	}, nil
}

func (s *ComunicadoService) GetByID(ctx context.Context, id string) (*models.Comunicado, error) {
	var c models.Comunicado
	err := s.db.Pool.QueryRow(ctx, `
		SELECT c.id, c.title, c.content, c.type, c.is_public, c.author_id,
		       COALESCE(u.name, '') as author_name, c.published_at, c.created_at, c.updated_at
		FROM comunicados c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.id = $1`, id).Scan(
		&c.ID, &c.Title, &c.Content, &c.Type, &c.IsPublic,
		&c.AuthorID, &c.AuthorName, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, ErrComunicadoNotFound
	}
	return &c, nil
}

func (s *ComunicadoService) Create(ctx context.Context, req *models.CreateComunicadoRequest, authorID string) (*models.Comunicado, error) {
	var publishedAt *time.Time
	if req.Publish {
		now := time.Now()
		publishedAt = &now
	}

	var c models.Comunicado
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO comunicados (title, content, type, is_public, author_id, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, content, type, is_public, author_id, published_at, created_at, updated_at`,
		req.Title, req.Content, req.Type, req.IsPublic, authorID, publishedAt).Scan(
		&c.ID, &c.Title, &c.Content, &c.Type, &c.IsPublic, &c.AuthorID, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *ComunicadoService) Update(ctx context.Context, id string, req *models.UpdateComunicadoRequest) (*models.Comunicado, error) {
	// Get current comunicado
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Content != nil {
		current.Content = *req.Content
	}
	if req.Type != nil {
		current.Type = *req.Type
	}
	if req.IsPublic != nil {
		current.IsPublic = *req.IsPublic
	}
	if req.Publish != nil && *req.Publish && current.PublishedAt == nil {
		now := time.Now()
		current.PublishedAt = &now
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE comunicados
		SET title = $1, content = $2, type = $3, is_public = $4, published_at = $5, updated_at = NOW()
		WHERE id = $6`,
		current.Title, current.Content, current.Type, current.IsPublic, current.PublishedAt, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *ComunicadoService) Delete(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM comunicados WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrComunicadoNotFound
	}
	return nil
}

func (s *ComunicadoService) GetLatest(ctx context.Context, limit int, publicOnly bool) ([]models.Comunicado, error) {
	query := `
		SELECT c.id, c.title, c.content, c.type, c.is_public, c.author_id,
		       COALESCE(u.name, '') as author_name, c.published_at, c.created_at, c.updated_at
		FROM comunicados c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.published_at IS NOT NULL`

	args := []interface{}{}
	if publicOnly {
		query += ` AND c.is_public = true`
	}
	query += ` ORDER BY c.published_at DESC LIMIT $1`
	args = append(args, limit)

	rows, err := s.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comunicados := []models.Comunicado{}
	for rows.Next() {
		var c models.Comunicado
		err := rows.Scan(&c.ID, &c.Title, &c.Content, &c.Type, &c.IsPublic,
			&c.AuthorID, &c.AuthorName, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comunicados = append(comunicados, c)
	}
	return comunicados, nil
}
