package services

import (
	"context"
	"errors"
	"strings"

	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrRegistrationClosed = errors.New("registration is disabled")
)

type AuthService struct {
	db         *database.DB
	jwtManager *jwt.JWTManager
}

func NewAuthService(db *database.DB, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if !isValidEmail(req.Email) {
		return nil, ErrInvalidEmail
	}

	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	// Check if user exists
	var exists bool
	err := s.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{}
	err = s.db.Pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, email, name, role, parcela_id, email_verified, created_at, updated_at`,
		req.Email, string(hashedPassword), req.Name, models.RoleVecino,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.ParcelaID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user := &models.User{}
	err := s.db.Pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, parcela_id, email_verified, created_at, updated_at
		 FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.ParcelaID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := s.db.Pool.QueryRow(ctx,
		`SELECT id, email, name, role, parcela_id, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.ParcelaID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, googleID, email, name string) (*models.AuthResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	// First, try to find user by Google ID
	user := &models.User{}
	err := s.db.Pool.QueryRow(ctx,
		`SELECT id, email, name, role, parcela_id, email_verified, created_at, updated_at
		 FROM users WHERE google_id = $1`,
		googleID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.ParcelaID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err == nil {
		// User found by Google ID, generate tokens
		return s.generateAuthResponse(user)
	}

	// Not found by Google ID, try by email
	err = s.db.Pool.QueryRow(ctx,
		`SELECT id, email, name, role, parcela_id, email_verified, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.ParcelaID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err == nil {
		// User exists with this email, link Google ID
		_, err = s.db.Pool.Exec(ctx,
			`UPDATE users SET google_id = $1, email_verified = true, updated_at = NOW() WHERE id = $2`,
			googleID, user.ID)
		if err != nil {
			return nil, err
		}
		user.EmailVerified = true
		return s.generateAuthResponse(user)
	}

	// User doesn't exist -> registration is closed. Accounts must be provisioned by directiva/backoffice.
	return nil, ErrRegistrationClosed
}

func (s *AuthService) generateAuthResponse(user *models.User) (*models.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
