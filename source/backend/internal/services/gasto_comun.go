package services

import (
	"context"
	"errors"
	"time"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
)

var (
	ErrPeriodoNotFound     = errors.New("periodo not found")
	ErrPeriodoExists       = errors.New("periodo already exists for this year/month")
	ErrGastoComunNotFound  = errors.New("gasto comun not found")
	ErrGastoAlreadyPaid    = errors.New("gasto already paid")
	ErrInvalidPaymentAmount = errors.New("invalid payment amount")
)

type GastoComunService struct {
	db *database.DB
}

func NewGastoComunService(db *database.DB) *GastoComunService {
	return &GastoComunService{db: db}
}

// ============================================
// PERIODOS
// ============================================

func (s *GastoComunService) ListPeriodos(ctx context.Context, filter models.PeriodoFilter) (*models.PeriodoListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 12
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT p.id, p.year, p.month, p.monto_base, p.fecha_vencimiento,
		       COALESCE(p.descripcion, ''), p.created_at, p.updated_at,
		       COUNT(g.id) as total_parcelas,
		       COUNT(g.id) FILTER (WHERE g.status = 'paid') as total_pagados,
		       COUNT(g.id) FILTER (WHERE g.status IN ('pending', 'overdue')) as total_pendientes,
		       COALESCE(SUM(g.monto_pagado), 0) as monto_recaudado,
		       COALESCE(SUM(g.monto) - SUM(g.monto_pagado), 0) as monto_pendiente
		FROM periodos_gasto p
		LEFT JOIN gastos_comunes g ON p.id = g.periodo_id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM periodos_gasto WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Year > 0 {
		argCount++
		query += ` AND p.year = $` + string(rune('0'+argCount))
		countQuery += ` AND year = $` + string(rune('0'+argCount))
		args = append(args, filter.Year)
	}

	query += ` GROUP BY p.id ORDER BY p.year DESC, p.month DESC`

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

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

	periodos := []models.PeriodoGasto{}
	for rows.Next() {
		var p models.PeriodoGasto
		err := rows.Scan(
			&p.ID, &p.Year, &p.Month, &p.MontoBase, &p.FechaVencimiento,
			&p.Descripcion, &p.CreatedAt, &p.UpdatedAt,
			&p.TotalParcelas, &p.TotalPagados, &p.TotalPendientes,
			&p.MontoRecaudado, &p.MontoPendiente)
		if err != nil {
			return nil, err
		}
		periodos = append(periodos, p)
	}

	return &models.PeriodoListResponse{
		Periodos: periodos,
		Total:    total,
		Page:     filter.Page,
		PerPage:  filter.PerPage,
	}, nil
}

func (s *GastoComunService) GetPeriodo(ctx context.Context, id string) (*models.PeriodoGasto, error) {
	var p models.PeriodoGasto
	err := s.db.Pool.QueryRow(ctx, `
		SELECT p.id, p.year, p.month, p.monto_base, p.fecha_vencimiento,
		       COALESCE(p.descripcion, ''), p.created_at, p.updated_at,
		       COUNT(g.id) as total_parcelas,
		       COUNT(g.id) FILTER (WHERE g.status = 'paid') as total_pagados,
		       COUNT(g.id) FILTER (WHERE g.status IN ('pending', 'overdue')) as total_pendientes,
		       COALESCE(SUM(g.monto_pagado), 0) as monto_recaudado,
		       COALESCE(SUM(g.monto) - SUM(g.monto_pagado), 0) as monto_pendiente
		FROM periodos_gasto p
		LEFT JOIN gastos_comunes g ON p.id = g.periodo_id
		WHERE p.id = $1
		GROUP BY p.id`, id).Scan(
		&p.ID, &p.Year, &p.Month, &p.MontoBase, &p.FechaVencimiento,
		&p.Descripcion, &p.CreatedAt, &p.UpdatedAt,
		&p.TotalParcelas, &p.TotalPagados, &p.TotalPendientes,
		&p.MontoRecaudado, &p.MontoPendiente)
	if err != nil {
		return nil, ErrPeriodoNotFound
	}
	return &p, nil
}

func (s *GastoComunService) GetPeriodoActual(ctx context.Context) (*models.PeriodoGasto, error) {
	now := time.Now()
	var p models.PeriodoGasto
	err := s.db.Pool.QueryRow(ctx, `
		SELECT p.id, p.year, p.month, p.monto_base, p.fecha_vencimiento,
		       COALESCE(p.descripcion, ''), p.created_at, p.updated_at,
		       COUNT(g.id) as total_parcelas,
		       COUNT(g.id) FILTER (WHERE g.status = 'paid') as total_pagados,
		       COUNT(g.id) FILTER (WHERE g.status IN ('pending', 'overdue')) as total_pendientes,
		       COALESCE(SUM(g.monto_pagado), 0) as monto_recaudado,
		       COALESCE(SUM(g.monto) - SUM(g.monto_pagado), 0) as monto_pendiente
		FROM periodos_gasto p
		LEFT JOIN gastos_comunes g ON p.id = g.periodo_id
		WHERE p.year = $1 AND p.month = $2
		GROUP BY p.id`, now.Year(), int(now.Month())).Scan(
		&p.ID, &p.Year, &p.Month, &p.MontoBase, &p.FechaVencimiento,
		&p.Descripcion, &p.CreatedAt, &p.UpdatedAt,
		&p.TotalParcelas, &p.TotalPagados, &p.TotalPendientes,
		&p.MontoRecaudado, &p.MontoPendiente)
	if err != nil {
		return nil, ErrPeriodoNotFound
	}
	return &p, nil
}

func (s *GastoComunService) CreatePeriodo(ctx context.Context, req *models.CreatePeriodoRequest) (*models.PeriodoGasto, error) {
	fechaVenc, err := time.Parse("2006-01-02", req.FechaVencimiento)
	if err != nil {
		return nil, errors.New("invalid fecha_vencimiento format")
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Check if periodo exists
	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM periodos_gasto WHERE year = $1 AND month = $2)`,
		req.Year, req.Month).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPeriodoExists
	}

	// Create periodo
	var periodoID string
	err = tx.QueryRow(ctx, `
		INSERT INTO periodos_gasto (year, month, monto_base, fecha_vencimiento, descripcion)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		req.Year, req.Month, req.MontoBase, fechaVenc, req.Descripcion).Scan(&periodoID)
	if err != nil {
		return nil, err
	}

	// Generate gastos_comunes for all parcelas
	_, err = tx.Exec(ctx, `
		INSERT INTO gastos_comunes (periodo_id, parcela_id, user_id, monto, status)
		SELECT $1, p.id, u.id, $2, 'pending'
		FROM parcelas p
		LEFT JOIN users u ON u.parcela_id = p.id AND u.role IN ('vecino', 'directiva')`,
		periodoID, req.MontoBase)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetPeriodo(ctx, periodoID)
}

func (s *GastoComunService) UpdatePeriodo(ctx context.Context, id string, req *models.UpdatePeriodoRequest) (*models.PeriodoGasto, error) {
	current, err := s.GetPeriodo(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.MontoBase != nil {
		current.MontoBase = *req.MontoBase
	}
	if req.FechaVencimiento != nil {
		fechaVenc, err := time.Parse("2006-01-02", *req.FechaVencimiento)
		if err != nil {
			return nil, errors.New("invalid fecha_vencimiento format")
		}
		current.FechaVencimiento = fechaVenc
	}
	if req.Descripcion != nil {
		current.Descripcion = *req.Descripcion
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE periodos_gasto
		SET monto_base = $1, fecha_vencimiento = $2, descripcion = $3, updated_at = NOW()
		WHERE id = $4`,
		current.MontoBase, current.FechaVencimiento, current.Descripcion, id)
	if err != nil {
		return nil, err
	}

	return s.GetPeriodo(ctx, id)
}

func (s *GastoComunService) GetResumen(ctx context.Context, periodoID string) (*models.ResumenGastos, error) {
	periodo, err := s.GetPeriodo(ctx, periodoID)
	if err != nil {
		return nil, err
	}

	var totalVencidos int
	err = s.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM gastos_comunes WHERE periodo_id = $1 AND status = 'overdue'`,
		periodoID).Scan(&totalVencidos)
	if err != nil {
		return nil, err
	}

	montoTotal := float64(periodo.TotalParcelas) * periodo.MontoBase
	porcentaje := float64(0)
	if montoTotal > 0 {
		porcentaje = (periodo.MontoRecaudado / montoTotal) * 100
	}

	return &models.ResumenGastos{
		Periodo:           *periodo,
		TotalParcelas:     periodo.TotalParcelas,
		TotalPagados:      periodo.TotalPagados,
		TotalPendientes:   periodo.TotalPendientes,
		TotalVencidos:     totalVencidos,
		MontoTotal:        montoTotal,
		MontoRecaudado:    periodo.MontoRecaudado,
		MontoPendiente:    periodo.MontoPendiente,
		PorcentajeRecaudo: porcentaje,
	}, nil
}

// ============================================
// GASTOS COMUNES
// ============================================

func (s *GastoComunService) ListGastos(ctx context.Context, filter models.GastoComunFilter) (*models.GastoComunListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT g.id, g.periodo_id, g.parcela_id, p.numero as parcela_numero,
		       g.user_id, COALESCE(u.name, '') as user_name,
		       g.monto, g.monto_pagado, g.status,
		       g.fecha_pago, COALESCE(g.metodo_pago, ''), COALESCE(g.referencia_pago, ''),
		       g.created_at, g.updated_at
		FROM gastos_comunes g
		JOIN parcelas p ON g.parcela_id = p.id
		LEFT JOIN users u ON g.user_id = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM gastos_comunes WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.PeriodoID != "" {
		argCount++
		query += ` AND g.periodo_id = $` + string(rune('0'+argCount))
		countQuery += ` AND periodo_id = $` + string(rune('0'+argCount))
		args = append(args, filter.PeriodoID)
	}

	if filter.ParcelaID > 0 {
		argCount++
		query += ` AND g.parcela_id = $` + string(rune('0'+argCount))
		countQuery += ` AND parcela_id = $` + string(rune('0'+argCount))
		args = append(args, filter.ParcelaID)
	}

	if filter.Status != "" {
		argCount++
		query += ` AND g.status = $` + string(rune('0'+argCount))
		countQuery += ` AND status = $` + string(rune('0'+argCount))
		args = append(args, filter.Status)
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY p.numero::int`
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

	gastos := []models.GastoComun{}
	for rows.Next() {
		var g models.GastoComun
		err := rows.Scan(
			&g.ID, &g.PeriodoID, &g.ParcelaID, &g.ParcelaNumero,
			&g.UserID, &g.UserName,
			&g.Monto, &g.MontoPagado, &g.Status,
			&g.FechaPago, &g.MetodoPago, &g.ReferenciaPago,
			&g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		gastos = append(gastos, g)
	}

	return &models.GastoComunListResponse{
		Gastos:  gastos,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}, nil
}

func (s *GastoComunService) GetGasto(ctx context.Context, id string) (*models.GastoComun, error) {
	var g models.GastoComun
	err := s.db.Pool.QueryRow(ctx, `
		SELECT g.id, g.periodo_id, g.parcela_id, p.numero as parcela_numero,
		       g.user_id, COALESCE(u.name, '') as user_name,
		       g.monto, g.monto_pagado, g.status,
		       g.fecha_pago, COALESCE(g.metodo_pago, ''), COALESCE(g.referencia_pago, ''),
		       g.created_at, g.updated_at
		FROM gastos_comunes g
		JOIN parcelas p ON g.parcela_id = p.id
		LEFT JOIN users u ON g.user_id = u.id
		WHERE g.id = $1`, id).Scan(
		&g.ID, &g.PeriodoID, &g.ParcelaID, &g.ParcelaNumero,
		&g.UserID, &g.UserName,
		&g.Monto, &g.MontoPagado, &g.Status,
		&g.FechaPago, &g.MetodoPago, &g.ReferenciaPago,
		&g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, ErrGastoComunNotFound
	}
	return &g, nil
}

func (s *GastoComunService) GetMiEstadoCuenta(ctx context.Context, userID string) (*models.MiEstadoCuenta, error) {
	// Get user's parcela
	var parcelaID int
	var parcelaNumero string
	err := s.db.Pool.QueryRow(ctx, `
		SELECT p.id, p.numero FROM parcelas p
		JOIN users u ON u.parcela_id = p.id
		WHERE u.id = $1`, userID).Scan(&parcelaID, &parcelaNumero)
	if err != nil {
		return nil, errors.New("user has no associated parcela")
	}

	// Get pending gastos
	rowsPendientes, err := s.db.Pool.Query(ctx, `
		SELECT g.id, g.periodo_id, g.parcela_id, p.numero,
		       g.user_id, '',
		       g.monto, g.monto_pagado, g.status,
		       g.fecha_pago, COALESCE(g.metodo_pago, ''), COALESCE(g.referencia_pago, ''),
		       g.created_at, g.updated_at,
		       pg.year, pg.month
		FROM gastos_comunes g
		JOIN parcelas p ON g.parcela_id = p.id
		JOIN periodos_gasto pg ON g.periodo_id = pg.id
		WHERE g.parcela_id = $1 AND g.status IN ('pending', 'overdue')
		ORDER BY pg.year DESC, pg.month DESC`, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rowsPendientes.Close()

	gastosPendientes := []models.GastoComun{}
	var totalPendiente float64
	for rowsPendientes.Next() {
		var g models.GastoComun
		var year, month int
		err := rowsPendientes.Scan(
			&g.ID, &g.PeriodoID, &g.ParcelaID, &g.ParcelaNumero,
			&g.UserID, &g.UserName,
			&g.Monto, &g.MontoPagado, &g.Status,
			&g.FechaPago, &g.MetodoPago, &g.ReferenciaPago,
			&g.CreatedAt, &g.UpdatedAt,
			&year, &month)
		if err != nil {
			return nil, err
		}
		g.Periodo = &models.PeriodoGasto{Year: year, Month: month}
		gastosPendientes = append(gastosPendientes, g)
		totalPendiente += (g.Monto - g.MontoPagado)
	}

	// Get paid gastos (last 12)
	rowsPagados, err := s.db.Pool.Query(ctx, `
		SELECT g.id, g.periodo_id, g.parcela_id, p.numero,
		       g.user_id, '',
		       g.monto, g.monto_pagado, g.status,
		       g.fecha_pago, COALESCE(g.metodo_pago, ''), COALESCE(g.referencia_pago, ''),
		       g.created_at, g.updated_at,
		       pg.year, pg.month
		FROM gastos_comunes g
		JOIN parcelas p ON g.parcela_id = p.id
		JOIN periodos_gasto pg ON g.periodo_id = pg.id
		WHERE g.parcela_id = $1 AND g.status = 'paid'
		ORDER BY pg.year DESC, pg.month DESC
		LIMIT 12`, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rowsPagados.Close()

	gastosPagados := []models.GastoComun{}
	var totalPagado float64
	for rowsPagados.Next() {
		var g models.GastoComun
		var year, month int
		err := rowsPagados.Scan(
			&g.ID, &g.PeriodoID, &g.ParcelaID, &g.ParcelaNumero,
			&g.UserID, &g.UserName,
			&g.Monto, &g.MontoPagado, &g.Status,
			&g.FechaPago, &g.MetodoPago, &g.ReferenciaPago,
			&g.CreatedAt, &g.UpdatedAt,
			&year, &month)
		if err != nil {
			return nil, err
		}
		g.Periodo = &models.PeriodoGasto{Year: year, Month: month}
		gastosPagados = append(gastosPagados, g)
		totalPagado += g.MontoPagado
	}

	return &models.MiEstadoCuenta{
		ParcelaID:        parcelaID,
		ParcelaNumero:    parcelaNumero,
		GastosPendientes: gastosPendientes,
		GastosPagados:    gastosPagados,
		TotalPendiente:   totalPendiente,
		TotalPagado:      totalPagado,
	}, nil
}

func (s *GastoComunService) RegistrarPago(ctx context.Context, gastoID string, req *models.RegistrarPagoRequest) (*models.GastoComun, error) {
	gasto, err := s.GetGasto(ctx, gastoID)
	if err != nil {
		return nil, err
	}

	if gasto.Status == models.PagoStatusPaid {
		return nil, ErrGastoAlreadyPaid
	}

	pendiente := gasto.Monto - gasto.MontoPagado
	if req.Monto <= 0 || req.Monto > pendiente {
		return nil, ErrInvalidPaymentAmount
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Register payment
	_, err = tx.Exec(ctx, `
		INSERT INTO pagos (gasto_comun_id, monto, metodo, referencia_externa, estado)
		VALUES ($1, $2, $3, $4, 'approved')`,
		gastoID, req.Monto, req.Metodo, req.ReferenciaExterna)
	if err != nil {
		return nil, err
	}

	// Update gasto_comun
	newMontoPagado := gasto.MontoPagado + req.Monto
	newStatus := models.PagoStatusPending
	var fechaPago *time.Time
	var metodoPago, referenciaPago string

	if newMontoPagado >= gasto.Monto {
		newStatus = models.PagoStatusPaid
		now := time.Now()
		fechaPago = &now
		metodoPago = req.Metodo
		referenciaPago = req.ReferenciaExterna
	}

	_, err = tx.Exec(ctx, `
		UPDATE gastos_comunes
		SET monto_pagado = $1, status = $2, fecha_pago = $3, metodo_pago = $4, referencia_pago = $5, updated_at = NOW()
		WHERE id = $6`,
		newMontoPagado, newStatus, fechaPago, metodoPago, referenciaPago, gastoID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetGasto(ctx, gastoID)
}

func (s *GastoComunService) MarcarVencidos(ctx context.Context) (int, error) {
	result, err := s.db.Pool.Exec(ctx, `
		UPDATE gastos_comunes g
		SET status = 'overdue', updated_at = NOW()
		FROM periodos_gasto p
		WHERE g.periodo_id = p.id
		  AND g.status = 'pending'
		  AND p.fecha_vencimiento < NOW()`)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}
