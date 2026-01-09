package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/condominio/backend/internal/models"
	"github.com/condominio/backend/internal/services"
	"github.com/condominio/backend/pkg/oauth"
)

type AuthHandler struct {
	authService   *services.AuthService
	googleService *oauth.GoogleService
	frontendURL   string
}

func NewAuthHandler(authService *services.AuthService, googleService *oauth.GoogleService, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		googleService: googleService,
		frontendURL:   frontendURL,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "Email, password and name are required")
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			writeError(w, http.StatusConflict, "User already exists")
		case errors.Is(err, services.ErrInvalidEmail):
			writeError(w, http.StatusBadRequest, "Invalid email format")
		case errors.Is(err, services.ErrWeakPassword):
			writeError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		default:
			writeError(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "Login failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	resp, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// GoogleLogin initiates the Google OAuth flow
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.googleService == nil || !h.googleService.IsConfigured() {
		writeError(w, http.StatusServiceUnavailable, "Google OAuth not configured")
		return
	}

	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}

	// Set state in cookie for validation in callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	authURL := h.googleService.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the Google OAuth callback
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.googleService == nil || !h.googleService.IsConfigured() {
		h.redirectWithError(w, r, "Google OAuth not configured")
		return
	}

	// Verify state matches
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		h.redirectWithError(w, r, "Invalid state")
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" || state != stateCookie.Value {
		h.redirectWithError(w, r, "State mismatch")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Check for error from Google
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Printf("[OAUTH] Google returned error: %s", errMsg)
		h.redirectWithError(w, r, "Access denied")
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "No authorization code")
		return
	}

	// Exchange code for user info
	userInfo, err := h.googleService.Authenticate(r.Context(), code)
	if err != nil {
		log.Printf("[OAUTH] Failed to authenticate with Google: %v", err)
		h.redirectWithError(w, r, "Authentication failed")
		return
	}

	// Validate email
	if userInfo.Email == "" || !userInfo.VerifiedEmail {
		h.redirectWithError(w, r, "Invalid or unverified email")
		return
	}

	// Login or register user
	authResp, err := h.authService.LoginWithGoogle(r.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		log.Printf("[OAUTH] Failed to login/register user: %v", err)
		if errors.Is(err, services.ErrRegistrationClosed) {
			h.redirectWithError(w, r, "Cuenta no habilitada. Solicita acceso a la directiva.")
			return
		}
		h.redirectWithError(w, r, "Authentication failed")
		return
	}

	// Redirect to frontend with tokens
	h.redirectWithTokens(w, r, authResp)
}

// redirectWithError redirects to frontend with error message
func (h *AuthHandler) redirectWithError(w http.ResponseWriter, r *http.Request, message string) {
	redirectURL, _ := url.Parse(h.frontendURL)
	redirectURL.Path = "/auth/callback"
	q := redirectURL.Query()
	q.Set("error", message)
	redirectURL.RawQuery = q.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

// redirectWithTokens redirects to frontend with auth tokens
func (h *AuthHandler) redirectWithTokens(w http.ResponseWriter, r *http.Request, authResp *models.AuthResponse) {
	redirectURL, _ := url.Parse(h.frontendURL)
	redirectURL.Path = "/auth/callback"
	q := redirectURL.Query()
	q.Set("access_token", authResp.AccessToken)
	q.Set("refresh_token", authResp.RefreshToken)
	redirectURL.RawQuery = q.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
