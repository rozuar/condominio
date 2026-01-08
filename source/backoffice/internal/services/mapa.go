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
	ErrMapaAreaNotFound  = errors.New("mapa area not found")
	ErrMapaPuntoNotFound = errors.New("mapa punto not found")
)

type MapaService struct {
	db *database.DB
}

func NewMapaService(db *database.DB) *MapaService {
	return &MapaService{db: db}
}

// GetAllMapData returns all map data for the frontend
func (s *MapaService) GetAllMapData(ctx context.Context) (*models.MapaData, error) {
	// Get all areas
	areasRows, err := s.db.Pool.Query(ctx, `
		SELECT id, parcela_id, type, name, COALESCE(description, ''), coordinates,
		       center_lat, center_lng, fill_color, stroke_color, is_clickable,
		       created_at, updated_at
		FROM mapa_areas
		ORDER BY type, name`)
	if err != nil {
		return nil, err
	}
	defer areasRows.Close()

	areas := []models.MapaArea{}
	for areasRows.Next() {
		var a models.MapaArea
		err := areasRows.Scan(&a.ID, &a.ParcelaID, &a.Type, &a.Name, &a.Description,
			&a.Coordinates, &a.CenterLat, &a.CenterLng, &a.FillColor, &a.StrokeColor,
			&a.IsClickable, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}

	// Get all public puntos
	puntosRows, err := s.db.Pool.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), lat, lng, icon, type, is_public,
		       created_at, updated_at
		FROM mapa_puntos
		WHERE is_public = TRUE
		ORDER BY type, name`)
	if err != nil {
		return nil, err
	}
	defer puntosRows.Close()

	puntos := []models.MapaPunto{}
	for puntosRows.Next() {
		var p models.MapaPunto
		err := puntosRows.Scan(&p.ID, &p.Name, &p.Description, &p.Lat, &p.Lng,
			&p.Icon, &p.Type, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		puntos = append(puntos, p)
	}

	return &models.MapaData{
		Areas:  areas,
		Puntos: puntos,
	}, nil
}

// Area operations

func (s *MapaService) ListAreas(ctx context.Context, filter models.MapaAreaFilter) (*models.MapaAreaListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 100
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT id, parcela_id, type, name, COALESCE(description, ''), coordinates,
		       center_lat, center_lng, fill_color, stroke_color, is_clickable,
		       created_at, updated_at
		FROM mapa_areas
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM mapa_areas WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Type != "" {
		argCount++
		query += ` AND type = $` + strconv.Itoa(argCount)
		countQuery += ` AND type = $` + strconv.Itoa(argCount)
		args = append(args, filter.Type)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY type, name`
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

	areas := []models.MapaArea{}
	for rows.Next() {
		var a models.MapaArea
		err := rows.Scan(&a.ID, &a.ParcelaID, &a.Type, &a.Name, &a.Description,
			&a.Coordinates, &a.CenterLat, &a.CenterLng, &a.FillColor, &a.StrokeColor,
			&a.IsClickable, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}

	return &models.MapaAreaListResponse{
		Areas:   areas,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}, nil
}

func (s *MapaService) GetAreaByID(ctx context.Context, id string) (*models.MapaArea, error) {
	var a models.MapaArea
	err := s.db.Pool.QueryRow(ctx, `
		SELECT id, parcela_id, type, name, COALESCE(description, ''), coordinates,
		       center_lat, center_lng, fill_color, stroke_color, is_clickable,
		       created_at, updated_at
		FROM mapa_areas
		WHERE id = $1`, id).Scan(
		&a.ID, &a.ParcelaID, &a.Type, &a.Name, &a.Description, &a.Coordinates,
		&a.CenterLat, &a.CenterLng, &a.FillColor, &a.StrokeColor, &a.IsClickable,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, ErrMapaAreaNotFound
	}
	return &a, nil
}

func (s *MapaService) CreateArea(ctx context.Context, req *models.CreateMapaAreaRequest) (*models.MapaArea, error) {
	fillColor := "#4A7C23"
	if req.FillColor != "" {
		fillColor = req.FillColor
	}
	strokeColor := "#2D5016"
	if req.StrokeColor != "" {
		strokeColor = req.StrokeColor
	}
	isClickable := true
	if req.IsClickable != nil {
		isClickable = *req.IsClickable
	}

	var a models.MapaArea
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO mapa_areas (parcela_id, type, name, description, coordinates, center_lat, center_lng,
		                        fill_color, stroke_color, is_clickable)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, parcela_id, type, name, COALESCE(description, ''), coordinates,
		          center_lat, center_lng, fill_color, stroke_color, is_clickable, created_at, updated_at`,
		req.ParcelaID, req.Type, req.Name, req.Description, req.Coordinates,
		req.CenterLat, req.CenterLng, fillColor, strokeColor, isClickable).Scan(
		&a.ID, &a.ParcelaID, &a.Type, &a.Name, &a.Description, &a.Coordinates,
		&a.CenterLat, &a.CenterLng, &a.FillColor, &a.StrokeColor, &a.IsClickable,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *MapaService) UpdateArea(ctx context.Context, id string, req *models.UpdateMapaAreaRequest) (*models.MapaArea, error) {
	_, err := s.GetAreaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `UPDATE mapa_areas SET updated_at = NOW()`
	args := []interface{}{}
	argCount := 0

	if req.ParcelaID != nil {
		argCount++
		query += fmt.Sprintf(`, parcela_id = $%d`, argCount)
		args = append(args, *req.ParcelaID)
	}
	if req.Type != nil {
		argCount++
		query += fmt.Sprintf(`, type = $%d`, argCount)
		args = append(args, *req.Type)
	}
	if req.Name != nil {
		argCount++
		query += fmt.Sprintf(`, name = $%d`, argCount)
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		argCount++
		query += fmt.Sprintf(`, description = $%d`, argCount)
		args = append(args, *req.Description)
	}
	if req.Coordinates != nil {
		argCount++
		query += fmt.Sprintf(`, coordinates = $%d`, argCount)
		args = append(args, *req.Coordinates)
	}
	if req.CenterLat != nil {
		argCount++
		query += fmt.Sprintf(`, center_lat = $%d`, argCount)
		args = append(args, *req.CenterLat)
	}
	if req.CenterLng != nil {
		argCount++
		query += fmt.Sprintf(`, center_lng = $%d`, argCount)
		args = append(args, *req.CenterLng)
	}
	if req.FillColor != nil {
		argCount++
		query += fmt.Sprintf(`, fill_color = $%d`, argCount)
		args = append(args, *req.FillColor)
	}
	if req.StrokeColor != nil {
		argCount++
		query += fmt.Sprintf(`, stroke_color = $%d`, argCount)
		args = append(args, *req.StrokeColor)
	}
	if req.IsClickable != nil {
		argCount++
		query += fmt.Sprintf(`, is_clickable = $%d`, argCount)
		args = append(args, *req.IsClickable)
	}

	argCount++
	query += fmt.Sprintf(` WHERE id = $%d`, argCount)
	args = append(args, id)

	_, err = s.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetAreaByID(ctx, id)
}

func (s *MapaService) DeleteArea(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM mapa_areas WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrMapaAreaNotFound
	}
	return nil
}

// Punto operations

func (s *MapaService) ListPuntos(ctx context.Context, filter models.MapaPuntoFilter) (*models.MapaPuntoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 100
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT id, name, COALESCE(description, ''), lat, lng, icon, type, is_public,
		       created_at, updated_at
		FROM mapa_puntos
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM mapa_puntos WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Type != "" {
		argCount++
		query += ` AND type = $` + strconv.Itoa(argCount)
		countQuery += ` AND type = $` + strconv.Itoa(argCount)
		args = append(args, filter.Type)
	}
	if filter.IsPublic != nil {
		argCount++
		query += ` AND is_public = $` + strconv.Itoa(argCount)
		countQuery += ` AND is_public = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsPublic)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY type, name`
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

	puntos := []models.MapaPunto{}
	for rows.Next() {
		var p models.MapaPunto
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Lat, &p.Lng,
			&p.Icon, &p.Type, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		puntos = append(puntos, p)
	}

	return &models.MapaPuntoListResponse{
		Puntos:  puntos,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}, nil
}

func (s *MapaService) GetPuntoByID(ctx context.Context, id string) (*models.MapaPunto, error) {
	var p models.MapaPunto
	err := s.db.Pool.QueryRow(ctx, `
		SELECT id, name, COALESCE(description, ''), lat, lng, icon, type, is_public,
		       created_at, updated_at
		FROM mapa_puntos
		WHERE id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Lat, &p.Lng, &p.Icon, &p.Type, &p.IsPublic,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, ErrMapaPuntoNotFound
	}
	return &p, nil
}

func (s *MapaService) CreatePunto(ctx context.Context, req *models.CreateMapaPuntoRequest) (*models.MapaPunto, error) {
	icon := "marker"
	if req.Icon != "" {
		icon = req.Icon
	}
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	var p models.MapaPunto
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO mapa_puntos (name, description, lat, lng, icon, type, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, COALESCE(description, ''), lat, lng, icon, type, is_public, created_at, updated_at`,
		req.Name, req.Description, req.Lat, req.Lng, icon, req.Type, isPublic).Scan(
		&p.ID, &p.Name, &p.Description, &p.Lat, &p.Lng, &p.Icon, &p.Type, &p.IsPublic,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *MapaService) UpdatePunto(ctx context.Context, id string, req *models.UpdateMapaPuntoRequest) (*models.MapaPunto, error) {
	_, err := s.GetPuntoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `UPDATE mapa_puntos SET updated_at = NOW()`
	args := []interface{}{}
	argCount := 0

	if req.Name != nil {
		argCount++
		query += fmt.Sprintf(`, name = $%d`, argCount)
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		argCount++
		query += fmt.Sprintf(`, description = $%d`, argCount)
		args = append(args, *req.Description)
	}
	if req.Lat != nil {
		argCount++
		query += fmt.Sprintf(`, lat = $%d`, argCount)
		args = append(args, *req.Lat)
	}
	if req.Lng != nil {
		argCount++
		query += fmt.Sprintf(`, lng = $%d`, argCount)
		args = append(args, *req.Lng)
	}
	if req.Icon != nil {
		argCount++
		query += fmt.Sprintf(`, icon = $%d`, argCount)
		args = append(args, *req.Icon)
	}
	if req.Type != nil {
		argCount++
		query += fmt.Sprintf(`, type = $%d`, argCount)
		args = append(args, *req.Type)
	}
	if req.IsPublic != nil {
		argCount++
		query += fmt.Sprintf(`, is_public = $%d`, argCount)
		args = append(args, *req.IsPublic)
	}

	argCount++
	query += fmt.Sprintf(` WHERE id = $%d`, argCount)
	args = append(args, id)

	_, err = s.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPuntoByID(ctx, id)
}

func (s *MapaService) DeletePunto(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM mapa_puntos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrMapaPuntoNotFound
	}
	return nil
}
