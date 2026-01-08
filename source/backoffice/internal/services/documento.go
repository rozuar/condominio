package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var ErrDocumentoNotFound = errors.New("documento not found")

type DocumentoService struct {
	db *database.DB
}

func NewDocumentoService(db *database.DB) *DocumentoService {
	return &DocumentoService{db: db}
}

func (s *DocumentoService) List(ctx context.Context, filter models.DocumentoFilter) (*models.DocumentoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT d.id, d.title, COALESCE(d.description, ''), COALESCE(d.file_url, ''),
		       d.category, d.is_public, d.created_by, COALESCE(u.name, ''), d.created_at, d.updated_at
		FROM documentos d
		LEFT JOIN users u ON d.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM documentos WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Category != "" {
		argCount++
		query += ` AND d.category = $` + strconv.Itoa(argCount)
		countQuery += ` AND category = $` + strconv.Itoa(argCount)
		args = append(args, filter.Category)
	}

	if filter.IsPublic != nil {
		argCount++
		query += ` AND d.is_public = $` + strconv.Itoa(argCount)
		countQuery += ` AND is_public = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsPublic)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY d.created_at DESC`
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

	documentos := []models.Documento{}
	for rows.Next() {
		var d models.Documento
		err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.FileURL,
			&d.Category, &d.IsPublic, &d.CreatedBy, &d.CreatorName, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		documentos = append(documentos, d)
	}

	return &models.DocumentoListResponse{
		Documentos: documentos,
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
	}, nil
}

func (s *DocumentoService) GetByID(ctx context.Context, id string) (*models.Documento, error) {
	var d models.Documento
	err := s.db.Pool.QueryRow(ctx, `
		SELECT d.id, d.title, COALESCE(d.description, ''), COALESCE(d.file_url, ''),
		       d.category, d.is_public, d.created_by, COALESCE(u.name, ''), d.created_at, d.updated_at
		FROM documentos d
		LEFT JOIN users u ON d.created_by = u.id
		WHERE d.id = $1`, id).Scan(
		&d.ID, &d.Title, &d.Description, &d.FileURL, &d.Category, &d.IsPublic,
		&d.CreatedBy, &d.CreatorName, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, ErrDocumentoNotFound
	}
	return &d, nil
}

func (s *DocumentoService) Create(ctx context.Context, req *models.CreateDocumentoRequest, createdBy string) (*models.Documento, error) {
	var d models.Documento
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO documentos (title, description, file_url, category, is_public, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, file_url, category, is_public, created_by, created_at, updated_at`,
		req.Title, req.Description, req.FileURL, req.Category, req.IsPublic, createdBy).Scan(
		&d.ID, &d.Title, &d.Description, &d.FileURL, &d.Category, &d.IsPublic, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
