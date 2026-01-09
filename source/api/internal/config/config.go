package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURLDirect string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string

	// JWT
	JWTSecret             string
	JWTExpiryHours        int
	JWTRefreshExpiryHours int

	// Server
	Port string
	Env  string

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string
	SMTPEnabled  bool

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURLDirect:     getEnv("DATABASE_URL", ""),
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		DBUser:                getEnv("DB_USER", "condominio"),
		DBPassword:            getEnv("DB_PASSWORD", "condominio123"),
		DBName:                getEnv("DB_NAME", "condominio"),
		JWTSecret:             getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiryHours:        getEnvInt("JWT_EXPIRY_HOURS", 24),
		JWTRefreshExpiryHours: getEnvInt("JWT_REFRESH_EXPIRY_HOURS", 168),
		Port:                  getEnv("PORT", "8080"),
		Env:                   getEnv("ENV", "development"),
		SMTPHost:              getEnv("SMTP_HOST", ""),
		SMTPPort:              getEnvInt("SMTP_PORT", 587),
		SMTPUser:              getEnv("SMTP_USER", ""),
		SMTPPassword:          getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:              getEnv("SMTP_FROM", "noreply@vinapelvin.cl"),
		SMTPFromName:          getEnv("SMTP_FROM_NAME", "Comunidad Viña Pelvin"),
		SMTPEnabled:           getEnvBool("SMTP_ENABLED", false),
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		FrontendURL:           getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func (c *Config) DatabaseURL() string {
	// Use direct DATABASE_URL if provided (for cloud databases like Railway)
	if c.DatabaseURLDirect != "" {
		return c.DatabaseURLDirect
	}
	// Otherwise build from individual components (for local development)
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
