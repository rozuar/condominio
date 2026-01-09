package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

type TesoreriaService struct {
	db *database.DB
}

func NewTesoreriaService(db *database.DB) *TesoreriaService {
	return &TesoreriaService{db: db}
}

func (s *TesoreriaService) List(ctx context.Context, filter models.MovimientoFilter) (*models.MovimientoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT m.id, m.description, m.amount, m.type, COALESCE(m.category, ''),
		       m.date, m.created_by, COALESCE(u.name, ''), m.created_at, m.updated_at
		FROM movimientos_tesoreria m
		LEFT JOIN users u ON m.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM movimientos_tesoreria WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Type != "" {
		argCount++
		query += ` AND m.type = $` + strconv.Itoa(argCount)
		countQuery += ` AND type = $` + strconv.Itoa(argCount)
		args = append(args, filter.Type)
	}

	if filter.Year > 0 {
		argCount++
		query += ` AND EXTRACT(YEAR FROM m.date) = $` + strconv.Itoa(argCount)
		countQuery += ` AND EXTRACT(YEAR FROM date) = $` + strconv.Itoa(argCount)
		args = append(args, filter.Year)
	}

	if filter.Month > 0 {
		argCount++
		query += ` AND EXTRACT(MONTH FROM m.date) = $` + strconv.Itoa(argCount)
		countQuery += ` AND EXTRACT(MONTH FROM date) = $` + strconv.Itoa(argCount)
		args = append(args, filter.Month)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY m.date DESC`
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

	movimientos := []models.Movimiento{}
	for rows.Next() {
		var m models.Movimiento
		err := rows.Scan(&m.ID, &m.Description, &m.Amount, &m.Type, &m.Category,
			&m.Date, &m.CreatedBy, &m.CreatorName, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}

	return &models.MovimientoListResponse{
		Movimientos: movimientos,
		Total:       total,
		Page:        filter.Page,
		PerPage:     filter.PerPage,
	}, nil
}

func (s *TesoreriaService) GetResumen(ctx context.Context) (*models.ResumenTesoreria, error) {
	var ingresos, egresos float64

	err := s.db.Pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM movimientos_tesoreria WHERE type = 'ingreso'`).Scan(&ingresos)
	if err != nil {
		return nil, err
	}

	err = s.db.Pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM movimientos_tesoreria WHERE type = 'egreso'`).Scan(&egresos)
	if err != nil {
		return nil, err
	}

	return &models.ResumenTesoreria{
		TotalIngresos: ingresos,
		TotalEgresos:  egresos,
		Balance:       ingresos - egresos,
	}, nil
}

func (s *TesoreriaService) Create(ctx context.Context, req *models.CreateMovimientoRequest, createdBy string) (*models.Movimiento, error) {
	var m models.Movimiento
	err := s.db.Pool.QueryRow(ctx, `
		INSERT INTO movimientos_tesoreria (description, amount, type, category, date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, description, amount, type, category, date, created_by, created_at, updated_at`,
		req.Description, req.Amount, req.Type, req.Category, req.Date, createdBy).Scan(
		&m.ID, &m.Description, &m.Amount, &m.Type, &m.Category, &m.Date, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
