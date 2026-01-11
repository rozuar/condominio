package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/condominio/backend/internal/handlers"
	"github.com/condominio/backend/internal/middleware"
	"github.com/condominio/backend/internal/services"
	"github.com/condominio/backend/pkg/jwt"
	"github.com/condominio/backend/pkg/oauth"
)

type Services struct {
	Auth         *services.AuthService
	Comunicado   *services.ComunicadoService
	Evento       *services.EventoService
	Tesoreria    *services.TesoreriaService
	Acta         *services.ActaService
	Documento    *services.DocumentoService
	Emergencia   *services.EmergenciaService
	Votacion     *services.VotacionService
	GastoComun   *services.GastoComunService
	Contacto     *services.ContactoService
	Galeria      *services.GaleriaService
	Mapa         *services.MapaService
	Notificacion *services.NotificacionService
}

type OAuthConfig struct {
	GoogleService *oauth.GoogleService
	FrontendURL   string
}

func New(svc *Services, jwtManager *jwt.JWTManager, oauthCfg *OAuthConfig) *chi.Mux {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.CORS)

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Handlers
	var googleService *oauth.GoogleService
	var frontendURL string
	if oauthCfg != nil {
		googleService = oauthCfg.GoogleService
		frontendURL = oauthCfg.FrontendURL
	}
	authHandler := handlers.NewAuthHandler(svc.Auth, googleService, frontendURL)
	comunicadoHandler := handlers.NewComunicadoHandler(svc.Comunicado)
	eventoHandler := handlers.NewEventoHandler(svc.Evento)
	tesoreriaHandler := handlers.NewTesoreriaHandler(svc.Tesoreria)
	actaHandler := handlers.NewActaHandler(svc.Acta)
	documentoHandler := handlers.NewDocumentoHandler(svc.Documento)
	emergenciaHandler := handlers.NewEmergenciaHandler(svc.Emergencia)
	votacionHandler := handlers.NewVotacionHandler(svc.Votacion)
	gastoComunHandler := handlers.NewGastoComunHandler(svc.GastoComun)
	contactoHandler := handlers.NewContactoHandler(svc.Contacto)
	galeriaHandler := handlers.NewGaleriaHandler(svc.Galeria)
	mapaHandler := handlers.NewMapaHandler(svc.Mapa)
	notificacionHandler := handlers.NewNotificacionHandler(svc.Notificacion)

	// Public routes
	r.Get("/health", handlers.Health)

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)

		// Google OAuth
		r.Get("/google", authHandler.GoogleLogin)
		r.Get("/google/callback", authHandler.GoogleCallback)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/me", authHandler.Me)
		})
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Comunicados - public endpoints
		r.Route("/comunicados", func(r chi.Router) {
			r.Get("/", comunicadoHandler.List)
			r.Get("/latest", comunicadoHandler.GetLatest)
			r.Get("/{id}", comunicadoHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", comunicadoHandler.Create)
				r.Put("/{id}", comunicadoHandler.Update)
				r.Delete("/{id}", comunicadoHandler.Delete)
			})
		})

		// Eventos - public endpoints
		r.Route("/eventos", func(r chi.Router) {
			r.Get("/", eventoHandler.List)
			r.Get("/upcoming", eventoHandler.GetUpcoming)
			r.Get("/{id}", eventoHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", eventoHandler.Create)
				r.Put("/{id}", eventoHandler.Update)
				r.Delete("/{id}", eventoHandler.Delete)
			})
		})

		// Tesorería - protected (vecino+)
		r.Route("/tesoreria", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(authMiddleware.RequireRole("vecino", "directiva"))

			r.Get("/", tesoreriaHandler.List)
			r.Get("/resumen", tesoreriaHandler.GetResumen)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", tesoreriaHandler.Create)
			})
		})

		// Actas - protected (vecino+)
		r.Route("/actas", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(authMiddleware.RequireRole("vecino", "directiva"))

			r.Get("/", actaHandler.List)
			r.Get("/{id}", actaHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", actaHandler.Create)
			})
		})

		// Documentos - protected (vecino+)
		r.Route("/documentos", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(authMiddleware.RequireRole("vecino", "directiva"))

			r.Get("/", documentoHandler.List)
			r.Get("/{id}", documentoHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", documentoHandler.Create)
			})
		})

		// Emergencias - public read, protected write
		r.Route("/emergencias", func(r chi.Router) {
			// Public endpoints
			r.Get("/", emergenciaHandler.List)
			r.Get("/active", emergenciaHandler.GetActive)
			r.Get("/{id}", emergenciaHandler.GetByID)

			// Protected endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", emergenciaHandler.Create)
				r.Put("/{id}", emergenciaHandler.Update)
				r.Post("/{id}/resolve", emergenciaHandler.Resolve)
				r.Delete("/{id}", emergenciaHandler.Delete)
			})
		})

		// Votaciones - protected (vecino+)
		r.Route("/votaciones", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Read endpoints: vecino/directiva/familia (admin is always allowed by middleware)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("vecino", "directiva", "familia"))
				r.Get("/", votacionHandler.List)
				r.Get("/active", votacionHandler.GetActive)
				r.Get("/{id}", votacionHandler.GetByID)
				r.Get("/{id}/resultados", votacionHandler.GetResultados)

				// Vote endpoint: allowed roles, but parcela is enforced at service level
				r.Post("/{id}/votar", votacionHandler.EmitirVoto)
			})

			// Admin endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", votacionHandler.Create)
				r.Put("/{id}", votacionHandler.Update)
				r.Post("/{id}/publish", votacionHandler.Publish)
				r.Post("/{id}/close", votacionHandler.Close)
				r.Post("/{id}/cancel", votacionHandler.Cancel)
				r.Delete("/{id}", votacionHandler.Delete)
			})
		})

		// Gastos Comunes - protected (vecino+)
		r.Route("/gastos", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(authMiddleware.RequireRole("vecino", "directiva"))

			// Read endpoints (vecino+)
			r.Get("/periodos", gastoComunHandler.ListPeriodos)
			r.Get("/periodos/actual", gastoComunHandler.GetPeriodoActual)
			r.Get("/periodos/{id}", gastoComunHandler.GetPeriodo)
			r.Get("/periodos/{id}/resumen", gastoComunHandler.GetResumen)
			r.Get("/periodos/{id}/gastos", gastoComunHandler.ListGastos)
			r.Get("/mi-cuenta", gastoComunHandler.GetMiEstadoCuenta)
			r.Get("/{id}", gastoComunHandler.GetGasto)

			// Admin endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/periodos", gastoComunHandler.CreatePeriodo)
				r.Put("/periodos/{id}", gastoComunHandler.UpdatePeriodo)
				r.Post("/{id}/pago", gastoComunHandler.RegistrarPago)
				r.Post("/marcar-vencidos", gastoComunHandler.MarcarVencidos)
			})
		})

		// Contacto - public create, protected management
		r.Route("/contacto", func(r chi.Router) {
			// Public endpoint - anyone can send a message
			r.Post("/", contactoHandler.Create)

			// Protected endpoints (vecino+)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("vecino", "directiva"))
				r.Get("/mis-mensajes", contactoHandler.GetMisMensajes)
			})

			// Admin endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Get("/", contactoHandler.List)
				r.Get("/stats", contactoHandler.GetStats)
				r.Get("/{id}", contactoHandler.GetByID)
				r.Post("/{id}/read", contactoHandler.MarkAsRead)
				r.Post("/{id}/reply", contactoHandler.Reply)
				r.Post("/{id}/archive", contactoHandler.Archive)
				r.Delete("/{id}", contactoHandler.Delete)
			})
		})

		// Galeria - public read for public galleries, protected for private
		r.Route("/galerias", func(r chi.Router) {
			// Public endpoints - list and view public galleries
			r.Get("/", galeriaHandler.List)
			r.Get("/{id}", galeriaHandler.GetByID)

			// Admin endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", galeriaHandler.Create)
				r.Put("/{id}", galeriaHandler.Update)
				r.Delete("/{id}", galeriaHandler.Delete)
				r.Post("/{id}/items", galeriaHandler.AddItem)
				r.Post("/{id}/reorder", galeriaHandler.ReorderItems)
				r.Put("/{id}/items/{itemId}", galeriaHandler.UpdateItem)
				r.Delete("/{id}/items/{itemId}", galeriaHandler.DeleteItem)
			})
		})

		// Mapa - public read, protected write
		r.Route("/mapa", func(r chi.Router) {
			// Public endpoint - get all map data
			r.Get("/", mapaHandler.GetAllMapData)

			// Areas endpoints
			r.Route("/areas", func(r chi.Router) {
				r.Get("/", mapaHandler.ListAreas)
				r.Get("/{id}", mapaHandler.GetAreaByID)

				// Admin endpoints (directiva only)
				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Use(authMiddleware.RequireRole("directiva"))
					r.Post("/", mapaHandler.CreateArea)
					r.Put("/{id}", mapaHandler.UpdateArea)
					r.Delete("/{id}", mapaHandler.DeleteArea)
				})
			})

			// Puntos endpoints
			r.Route("/puntos", func(r chi.Router) {
				r.Get("/", mapaHandler.ListPuntos)
				r.Get("/{id}", mapaHandler.GetPuntoByID)

				// Admin endpoints (directiva only)
				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Use(authMiddleware.RequireRole("directiva"))
					r.Post("/", mapaHandler.CreatePunto)
					r.Put("/{id}", mapaHandler.UpdatePunto)
					r.Delete("/{id}", mapaHandler.DeletePunto)
				})
			})
		})

		// Notificaciones - protected (vecino+)
		r.Route("/notificaciones", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(authMiddleware.RequireRole("vecino", "directiva"))

			// User endpoints
			r.Get("/", notificacionHandler.List)
			r.Get("/stats", notificacionHandler.GetStats)
			r.Get("/{id}", notificacionHandler.GetByID)
			r.Post("/{id}/read", notificacionHandler.MarkAsRead)
			r.Post("/read-all", notificacionHandler.MarkAllAsRead)
			r.Delete("/{id}", notificacionHandler.Delete)
			r.Delete("/", notificacionHandler.DeleteAll)
			r.Delete("/read", notificacionHandler.DeleteRead)

			// Admin endpoints (directiva only)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRole("directiva"))
				r.Post("/", notificacionHandler.Create)
				r.Post("/bulk", notificacionHandler.CreateBulk)
				r.Post("/broadcast", notificacionHandler.CreateBroadcast)
			})
		})
	})

	return r
}
