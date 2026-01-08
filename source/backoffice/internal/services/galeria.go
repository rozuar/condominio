package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var (
	ErrGaleriaNotFound     = errors.New("galeria not found")
	ErrGaleriaItemNotFound = errors.New("galeria item not found")
)

type GaleriaService struct {
	db *database.DB
}

func NewGaleriaService(db *database.DB) *GaleriaService {
	return &GaleriaService{db: db}
}

func (s *GaleriaService) List(ctx context.Context, filter models.GaleriaFilter) (*models.GaleriaListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT g.id, g.title, COALESCE(g.description, ''), g.event_date,
		       g.is_public, COALESCE(g.cover_image_url, ''),
		       (SELECT COUNT(*) FROM galeria_items gi WHERE gi.galeria_id = g.id) as items_count,
		       g.created_by, COALESCE(u.name, ''), g.created_at, g.updated_at
		FROM galerias g
		LEFT JOIN users u ON g.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM galerias WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.IsPublic != nil {
		argCount++
		query += ` AND g.is_public = $` + strconv.Itoa(argCount)
		countQuery += ` AND is_public = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsPublic)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY g.event_date DESC NULLS LAST, g.created_at DESC`
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

	galerias := []models.Galeria{}
	for rows.Next() {
		var g models.Galeria
		err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.EventDate,
			&g.IsPublic, &g.CoverImageURL, &g.ItemsCount,
			&g.CreatedBy, &g.CreatorName, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		galerias = append(galerias, g)
	}

	return &models.GaleriaListResponse{
		Galerias: galerias,
		Total:    total,
		Page:     filter.Page,
		PerPage:  filter.PerPage,
	}, nil
}

func (s *GaleriaService) GetByID(ctx context.Context, id string) (*models.Galeria, error) {
	var g models.Galeria
	err := s.db.Pool.QueryRow(ctx, `
		SELECT g.id, g.title, COALESCE(g.description, ''), g.event_date,
		       g.is_public, COALESCE(g.cover_image_url, ''),
		       (SELECT COUNT(*) FROM galeria_items gi WHERE gi.galeria_id = g.id) as items_count,
		       g.created_by, COALESCE(u.name, ''), g.created_at, g.updated_at
		FROM galerias g
		LEFT JOIN users u ON g.created_by = u.id
		WHERE g.id = $1`, id).Scan(
		&g.ID, &g.Title, &g.Description, &g.EventDate, &g.IsPublic, &g.CoverImageURL,
		&g.ItemsCount, &g.CreatedBy, &g.CreatorName, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, ErrGaleriaNotFound
	}
	return &g, nil
}

func (s *GaleriaService) GetWithItems(ctx context.Context, id string) (*models.GaleriaWithItems, error) {
	galeria, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Pool.Query(ctx, `
		SELECT id, galeria_id, file_url, COALESCE(thumbnail_url, ''), file_type,
		       COALESCE(caption, ''), order_index, created_at
		FROM galeria_items
		WHERE galeria_id = $1
		ORDER BY order_index ASC, created_at ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.GaleriaItem{}
	for rows.Next() {
		var item models.GaleriaItem
		err := rows.Scan(&item.ID, &item.GaleriaID, &item.FileURL, &item.ThumbnailURL,
			&item.FileType, &item.Caption, &item.OrderIndex, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &models.GaleriaWithItems{
		Galeria: *galeria,
		Items:   items,
	}, nil
}

func (s *GaleriaService) Create(ctx context.Context, req *models.CreateGaleriaRequest, createdBy string) (*models.Galeria, error) {
	var g models.Galeria
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO galerias (title, description, event_date, is_public, cover_image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, COALESCE(description, ''), event_date, is_public,
		          COALESCE(cover_image_url, ''), created_by, created_at, updated_at`,
		req.Title, req.Description, req.EventDate, req.IsPublic, req.CoverImageURL, createdBy).Scan(
		&g.ID, &g.Title, &g.Description, &g.EventDate, &g.IsPublic,
		&g.CoverImageURL, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	g.ItemsCount = 0
	return &g, nil
}

func (s *GaleriaService) Update(ctx context.Context, id string, req *models.UpdateGaleriaRequest) (*models.Galeria, error) {
	// First check if galeria exists
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `UPDATE galerias SET updated_at = NOW()`
	args := []interface{}{}
	argCount := 0

	if req.Title != nil {
		argCount++
		query += fmt.Sprintf(`, title = $%d`, argCount)
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		argCount++
		query += fmt.Sprintf(`, description = $%d`, argCount)
		args = append(args, *req.Description)
	}
	if req.EventDate != nil {
		argCount++
		query += fmt.Sprintf(`, event_date = $%d`, argCount)
		args = append(args, *req.EventDate)
	}
	if req.IsPublic != nil {
		argCount++
		query += fmt.Sprintf(`, is_public = $%d`, argCount)
		args = append(args, *req.IsPublic)
	}
	if req.CoverImageURL != nil {
		argCount++
		query += fmt.Sprintf(`, cover_image_url = $%d`, argCount)
		args = append(args, *req.CoverImageURL)
	}

	argCount++
	query += fmt.Sprintf(` WHERE id = $%d`, argCount)
	args = append(args, id)

	_, err = s.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *GaleriaService) Delete(ctx context.Context, id string) error {
	// Items are deleted via CASCADE
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM galerias WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrGaleriaNotFound
	}
	return nil
}

// Item operations

func (s *GaleriaService) AddItem(ctx context.Context, galeriaID string, req *models.AddGaleriaItemRequest) (*models.GaleriaItem, error) {
	// Check if galeria exists
	_, err := s.GetByID(ctx, galeriaID)
	if err != nil {
		return nil, err
	}

	var item models.GaleriaItem
	err = s.db.Pool.QueryRow(ctx, `
		INSERT INTO galeria_items (galeria_id, file_url, thumbnail_url, file_type, caption, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, galeria_id, file_url, COALESCE(thumbnail_url, ''), file_type,
		          COALESCE(caption, ''), order_index, created_at`,
		galeriaID, req.FileURL, req.ThumbnailURL, req.FileType, req.Caption, req.OrderIndex).Scan(
		&item.ID, &item.GaleriaID, &item.FileURL, &item.ThumbnailURL,
		&item.FileType, &item.Caption, &item.OrderIndex, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *GaleriaService) GetItem(ctx context.Context, itemID string) (*models.GaleriaItem, error) {
	var item models.GaleriaItem
	err := s.db.Pool.QueryRow(ctx, `
		SELECT id, galeria_id, file_url, COALESCE(thumbnail_url, ''), file_type,
		       COALESCE(caption, ''), order_index, created_at
		FROM galeria_items
		WHERE id = $1`, itemID).Scan(
		&item.ID, &item.GaleriaID, &item.FileURL, &item.ThumbnailURL,
		&item.FileType, &item.Caption, &item.OrderIndex, &item.CreatedAt)
	if err != nil {
		return nil, ErrGaleriaItemNotFound
	}
	return &item, nil
}

func (s *GaleriaService) UpdateItem(ctx context.Context, itemID string, req *models.UpdateGaleriaItemRequest) (*models.GaleriaItem, error) {
	// First check if item exists
	_, err := s.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	query := `UPDATE galeria_items SET id = id`
	args := []interface{}{}
	argCount := 0

	if req.FileURL != nil {
		argCount++
		query += fmt.Sprintf(`, file_url = $%d`, argCount)
		args = append(args, *req.FileURL)
	}
	if req.ThumbnailURL != nil {
		argCount++
		query += fmt.Sprintf(`, thumbnail_url = $%d`, argCount)
		args = append(args, *req.ThumbnailURL)
	}
	if req.FileType != nil {
		argCount++
		query += fmt.Sprintf(`, file_type = $%d`, argCount)
		args = append(args, *req.FileType)
	}
	if req.Caption != nil {
		argCount++
		query += fmt.Sprintf(`, caption = $%d`, argCount)
		args = append(args, *req.Caption)
	}
	if req.OrderIndex != nil {
		argCount++
		query += fmt.Sprintf(`, order_index = $%d`, argCount)
		args = append(args, *req.OrderIndex)
	}

	argCount++
	query += fmt.Sprintf(` WHERE id = $%d`, argCount)
	args = append(args, itemID)

	_, err = s.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetItem(ctx, itemID)
}

func (s *GaleriaService) DeleteItem(ctx context.Context, itemID string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM galeria_items WHERE id = $1`, itemID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrGaleriaItemNotFound
	}
	return nil
}

func (s *GaleriaService) ReorderItems(ctx context.Context, galeriaID string, itemIDs []string) error {
	// Check if galeria exists
	_, err := s.GetByID(ctx, galeriaID)
	if err != nil {
		return err
	}

	// Update order for each item
	for i, itemID := range itemIDs {
		_, err := s.db.Pool.Exec(ctx, `
			UPDATE galeria_items SET order_index = $1
			WHERE id = $2 AND galeria_id = $3`,
			i, itemID, galeriaID)
		if err != nil {
			return err
		}
	}

	return nil
}
