package models

import (
	"time"
)

type Role string

const (
	RoleVisitor   Role = "visitor"
	RoleVecino    Role = "vecino"
	RoleDirectiva Role = "directiva"
)

type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Name          string     `json:"name"`
	Role          Role       `json:"role"`
	ParcelaID     *int       `json:"parcela_id,omitempty"`
	GoogleID      *string    `json:"-"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (u *User) IsDirectiva() bool {
	return u.Role == RoleDirectiva
}

func (u *User) IsVecino() bool {
	return u.Role == RoleVecino || u.Role == RoleDirectiva
}
