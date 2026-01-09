package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/condominio/backend/internal/config"
	"github.com/condominio/backend/internal/database"
	"github.com/condominio/backend/internal/router"
	"github.com/condominio/backend/internal/services"
	"github.com/condominio/backend/pkg/email"
	"github.com/condominio/backend/pkg/jwt"
	"github.com/condominio/backend/pkg/oauth"
)

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting server in %s mode...", cfg.Env)

	// Connect to database
	db, err := database.New(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Run migrations
	if err := db.RunMigrations(context.Background()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTExpiryHours,
		cfg.JWTRefreshExpiryHours,
	)

	// Initialize email service
	emailSvc := email.NewService(email.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
		FromName: cfg.SMTPFromName,
		Enabled:  cfg.SMTPEnabled,
	})

	// Register email templates
	if err := email.RegisterAllTemplates(emailSvc); err != nil {
		log.Fatalf("Failed to register email templates: %v", err)
	}

	if emailSvc.IsEnabled() {
		log.Println("Email service enabled")
	} else {
		log.Println("Email service disabled (SMTP not configured)")
	}

	// Initialize services
	svc := &router.Services{
		Auth:         services.NewAuthService(db, jwtManager),
		Comunicado:   services.NewComunicadoService(db),
		Evento:       services.NewEventoService(db),
		Tesoreria:    services.NewTesoreriaService(db),
		Acta:         services.NewActaService(db),
		Documento:    services.NewDocumentoService(db),
		Emergencia:   services.NewEmergenciaService(db),
		Votacion:     services.NewVotacionService(db),
		GastoComun:   services.NewGastoComunService(db),
		Contacto:     services.NewContactoService(db, emailSvc),
		Galeria:      services.NewGaleriaService(db),
		Mapa:         services.NewMapaService(db),
		Notificacion: services.NewNotificacionService(db, emailSvc),
	}

	// Initialize Google OAuth service
	googleSvc := oauth.NewGoogleService(oauth.GoogleConfig{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
	})

	if googleSvc.IsConfigured() {
		log.Println("Google OAuth enabled")
	} else {
		log.Println("Google OAuth disabled (credentials not configured)")
	}

	// Initialize router
	oauthCfg := &router.OAuthConfig{
		GoogleService: googleSvc,
		FrontendURL:   cfg.FrontendURL,
	}
	r := router.New(svc, jwtManager, oauthCfg)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
