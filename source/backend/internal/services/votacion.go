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
	ErrVotacionNotFound    = errors.New("votacion not found")
	ErrVotacionNotActive   = errors.New("votacion is not active")
	ErrAlreadyVoted        = errors.New("user has already voted")
	ErrInvalidOpcion       = errors.New("invalid option for this votacion")
	ErrAbstentionNotAllowed = errors.New("abstention is not allowed")
)

type VotacionService struct {
	db *database.DB
}

func NewVotacionService(db *database.DB) *VotacionService {
	return &VotacionService{db: db}
}

func (s *VotacionService) List(ctx context.Context, filter models.VotacionFilter) (*models.VotacionListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	offset := (filter.Page - 1) * filter.PerPage

	query := `
		SELECT v.id, v.title, v.description, v.status, v.start_date, v.end_date,
		       v.requires_quorum, v.quorum_percentage, v.allow_abstention,
		       v.created_by, COALESCE(u.name, '') as creator_name,
		       v.created_at, v.updated_at,
		       (SELECT COUNT(*) FROM votos WHERE votacion_id = v.id) as total_votos
		FROM votaciones v
		LEFT JOIN users u ON v.created_by = u.id
		WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM votaciones WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(` AND v.status = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND status = $%d`, argCount)
		args = append(args, filter.Status)
	}

	if filter.Active {
		query += ` AND v.status = 'active'`
		countQuery += ` AND status = 'active'`
	}

	var total int
	err := s.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += ` ORDER BY
		CASE v.status
			WHEN 'active' THEN 1
			WHEN 'draft' THEN 2
			WHEN 'closed' THEN 3
			WHEN 'cancelled' THEN 4
		END,
		v.created_at DESC`

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

	votaciones := []models.Votacion{}
	for rows.Next() {
		var v models.Votacion
		err := rows.Scan(
			&v.ID, &v.Title, &v.Description, &v.Status, &v.StartDate, &v.EndDate,
			&v.RequiresQuorum, &v.QuorumPercentage, &v.AllowAbstention,
			&v.CreatedBy, &v.CreatorName,
			&v.CreatedAt, &v.UpdatedAt, &v.TotalVotos)
		if err != nil {
			return nil, err
		}
		votaciones = append(votaciones, v)
	}

	return &models.VotacionListResponse{
		Votaciones: votaciones,
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
	}, nil
}

func (s *VotacionService) GetByID(ctx context.Context, id string, userID *string) (*models.Votacion, error) {
	var v models.Votacion
	err := s.db.Pool.QueryRow(ctx, `
		SELECT v.id, v.title, v.description, v.status, v.start_date, v.end_date,
		       v.requires_quorum, v.quorum_percentage, v.allow_abstention,
		       v.created_by, COALESCE(u.name, '') as creator_name,
		       v.created_at, v.updated_at
		FROM votaciones v
		LEFT JOIN users u ON v.created_by = u.id
		WHERE v.id = $1`, id).Scan(
		&v.ID, &v.Title, &v.Description, &v.Status, &v.StartDate, &v.EndDate,
		&v.RequiresQuorum, &v.QuorumPercentage, &v.AllowAbstention,
		&v.CreatedBy, &v.CreatorName,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, ErrVotacionNotFound
	}

	// Get opciones
	opciones, err := s.getOpciones(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Opciones = opciones

	// Get total votos
	err = s.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM votos WHERE votacion_id = $1`, id).Scan(&v.TotalVotos)
	if err != nil {
		return nil, err
	}

	// Check if user has voted
	if userID != nil {
		var count int
		err = s.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM votos WHERE votacion_id = $1 AND user_id = $2`, id, *userID).Scan(&count)
		if err != nil {
			return nil, err
		}
		v.HasVoted = count > 0
	}

	return &v, nil
}

func (s *VotacionService) getOpciones(ctx context.Context, votacionID string) ([]models.VotacionOpcion, error) {
	rows, err := s.db.Pool.Query(ctx, `
		SELECT o.id, o.votacion_id, o.label, COALESCE(o.description, ''), o.order_index,
		       (SELECT COUNT(*) FROM votos WHERE opcion_id = o.id) as votos_count
		FROM votacion_opciones o
		WHERE o.votacion_id = $1
		ORDER BY o.order_index`, votacionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	opciones := []models.VotacionOpcion{}
	for rows.Next() {
		var o models.VotacionOpcion
		err := rows.Scan(&o.ID, &o.VotacionID, &o.Label, &o.Description, &o.OrderIndex, &o.VotosCount)
		if err != nil {
			return nil, err
		}
		opciones = append(opciones, o)
	}
	return opciones, nil
}

func (s *VotacionService) GetActive(ctx context.Context, userID *string) ([]models.Votacion, error) {
	rows, err := s.db.Pool.Query(ctx, `
		SELECT v.id, v.title, v.description, v.status, v.start_date, v.end_date,
		       v.requires_quorum, v.quorum_percentage, v.allow_abstention,
		       v.created_by, COALESCE(u.name, '') as creator_name,
		       v.created_at, v.updated_at,
		       (SELECT COUNT(*) FROM votos WHERE votacion_id = v.id) as total_votos
		FROM votaciones v
		LEFT JOIN users u ON v.created_by = u.id
		WHERE v.status = 'active'
		  AND (v.start_date IS NULL OR v.start_date <= NOW())
		  AND (v.end_date IS NULL OR v.end_date > NOW())
		ORDER BY v.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votaciones := []models.Votacion{}
	for rows.Next() {
		var v models.Votacion
		err := rows.Scan(
			&v.ID, &v.Title, &v.Description, &v.Status, &v.StartDate, &v.EndDate,
			&v.RequiresQuorum, &v.QuorumPercentage, &v.AllowAbstention,
			&v.CreatedBy, &v.CreatorName,
			&v.CreatedAt, &v.UpdatedAt, &v.TotalVotos)
		if err != nil {
			return nil, err
		}

		// Check if user has voted
		if userID != nil {
			var count int
			s.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM votos WHERE votacion_id = $1 AND user_id = $2`, v.ID, *userID).Scan(&count)
			v.HasVoted = count > 0
		}

		votaciones = append(votaciones, v)
	}
	return votaciones, nil
}

func (s *VotacionService) Create(ctx context.Context, req *models.CreateVotacionRequest, createdBy string) (*models.Votacion, error) {
	if req.QuorumPercentage <= 0 {
		req.QuorumPercentage = 50
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var votacionID string
	err = tx.QueryRow(ctx, `
		INSERT INTO votaciones (title, description, status, requires_quorum, quorum_percentage, allow_abstention, created_by)
		VALUES ($1, $2, 'draft', $3, $4, $5, $6)
		RETURNING id`,
		req.Title, req.Description, req.RequiresQuorum, req.QuorumPercentage, req.AllowAbstention, createdBy).Scan(&votacionID)
	if err != nil {
		return nil, err
	}

	// Insert opciones
	for i, label := range req.Opciones {
		_, err = tx.Exec(ctx, `
			INSERT INTO votacion_opciones (votacion_id, label, order_index)
			VALUES ($1, $2, $3)`, votacionID, label, i)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, votacionID, nil)
}

func (s *VotacionService) Update(ctx context.Context, id string, req *models.UpdateVotacionRequest) (*models.Votacion, error) {
	current, err := s.GetByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	if current.Status != models.VotacionStatusDraft {
		return nil, errors.New("can only update draft votaciones")
	}

	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.RequiresQuorum != nil {
		current.RequiresQuorum = *req.RequiresQuorum
	}
	if req.QuorumPercentage != nil {
		current.QuorumPercentage = *req.QuorumPercentage
	}
	if req.AllowAbstention != nil {
		current.AllowAbstention = *req.AllowAbstention
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE votaciones
		SET title = $1, description = $2, requires_quorum = $3, quorum_percentage = $4, allow_abstention = $5, updated_at = NOW()
		WHERE id = $6`,
		current.Title, current.Description, current.RequiresQuorum, current.QuorumPercentage, current.AllowAbstention, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, nil)
}

func (s *VotacionService) Publish(ctx context.Context, id string, startDate, endDate *time.Time) (*models.Votacion, error) {
	current, err := s.GetByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	if current.Status != models.VotacionStatusDraft {
		return nil, errors.New("can only publish draft votaciones")
	}

	if len(current.Opciones) < 2 {
		return nil, errors.New("votacion must have at least 2 options")
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE votaciones
		SET status = 'active', start_date = $1, end_date = $2, updated_at = NOW()
		WHERE id = $3`, startDate, endDate, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, nil)
}

func (s *VotacionService) Close(ctx context.Context, id string) (*models.Votacion, error) {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE votaciones
		SET status = 'closed', updated_at = NOW()
		WHERE id = $1 AND status = 'active'`, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, nil)
}

func (s *VotacionService) Cancel(ctx context.Context, id string) (*models.Votacion, error) {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE votaciones
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1 AND status IN ('draft', 'active')`, id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, nil)
}

func (s *VotacionService) EmitirVoto(ctx context.Context, votacionID string, userID string, req *models.EmitirVotoRequest) error {
	votacion, err := s.GetByID(ctx, votacionID, &userID)
	if err != nil {
		return err
	}

	if !votacion.CanVote() {
		return ErrVotacionNotActive
	}

	if votacion.HasVoted {
		return ErrAlreadyVoted
	}

	if req.IsAbstention {
		if !votacion.AllowAbstention {
			return ErrAbstentionNotAllowed
		}
		_, err = s.db.Pool.Exec(ctx, `
			INSERT INTO votos (votacion_id, user_id, is_abstention, voted_at)
			VALUES ($1, $2, true, NOW())`, votacionID, userID)
	} else {
		// Validate opcion
		validOpcion := false
		for _, o := range votacion.Opciones {
			if o.ID == *req.OpcionID {
				validOpcion = true
				break
			}
		}
		if !validOpcion {
			return ErrInvalidOpcion
		}

		_, err = s.db.Pool.Exec(ctx, `
			INSERT INTO votos (votacion_id, user_id, opcion_id, is_abstention, voted_at)
			VALUES ($1, $2, $3, false, NOW())`, votacionID, userID, req.OpcionID)
	}

	return err
}

func (s *VotacionService) GetResultados(ctx context.Context, id string) (*models.VotacionResultado, error) {
	votacion, err := s.GetByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	// Get total vecinos (users with role vecino or directiva)
	var totalVecinos int
	err = s.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role IN ('vecino', 'directiva')`).Scan(&totalVecinos)
	if err != nil {
		return nil, err
	}

	// Get total votos and abstenciones
	var totalVotos, totalAbstenciones int
	err = s.db.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_abstention = true) as abstenciones
		FROM votos WHERE votacion_id = $1`, id).Scan(&totalVotos, &totalAbstenciones)
	if err != nil {
		return nil, err
	}

	// Build resultados
	resultados := []models.OpcionResultado{}
	votosEfectivos := totalVotos - totalAbstenciones
	for _, o := range votacion.Opciones {
		percentage := float64(0)
		if votosEfectivos > 0 {
			percentage = float64(o.VotosCount) / float64(votosEfectivos) * 100
		}
		resultados = append(resultados, models.OpcionResultado{
			OpcionID:   o.ID,
			Label:      o.Label,
			Count:      o.VotosCount,
			Percentage: percentage,
		})
	}

	// Check quorum
	participacion := float64(0)
	if totalVecinos > 0 {
		participacion = float64(totalVotos) / float64(totalVecinos) * 100
	}
	quorumAlcanzado := !votacion.RequiresQuorum || participacion >= float64(votacion.QuorumPercentage)

	return &models.VotacionResultado{
		Votacion:         *votacion,
		TotalVotos:       totalVotos,
		TotalAbstenciones: totalAbstenciones,
		Resultados:       resultados,
		QuorumAlcanzado:  quorumAlcanzado,
		TotalVecinos:     totalVecinos,
		Participacion:    participacion,
	}, nil
}

func (s *VotacionService) Delete(ctx context.Context, id string) error {
	result, err := s.db.Pool.Exec(ctx, `DELETE FROM votaciones WHERE id = $1 AND status = 'draft'`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("can only delete draft votaciones")
	}
	return nil
}
