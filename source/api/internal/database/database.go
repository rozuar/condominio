package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) RunMigrations(ctx context.Context) error {
	migrations := []string{
		migrationUsers,
		migrationComunicados,
		migrationEventos,
		migrationTesoreria,
		migrationActas,
		migrationDocumentos,
		migrationEmergencias,
		migrationVotaciones,
		migrationGalerias,
		migrationMapaPuntos,
		migrationMensajesContacto,
		migrationNotificaciones,
	}

	for i, migration := range migrations {
		if _, err := db.Pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	return nil
}

const migrationUsers = `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'vecino',
    parcela_id INTEGER,
    google_id VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
`

const migrationComunicados = `
CREATE TABLE IF NOT EXISTS comunicados (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'informativo',
    is_public BOOLEAN NOT NULL DEFAULT true,
    author_id UUID REFERENCES users(id),
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comunicados_type ON comunicados(type);
CREATE INDEX IF NOT EXISTS idx_comunicados_is_public ON comunicados(is_public);
CREATE INDEX IF NOT EXISTS idx_comunicados_published_at ON comunicados(published_at DESC);
`

const migrationEventos = `
CREATE TABLE IF NOT EXISTS eventos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_end_date TIMESTAMP WITH TIME ZONE,
    location VARCHAR(255),
    type VARCHAR(50) NOT NULL DEFAULT 'reunion',
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_eventos_event_date ON eventos(event_date);
CREATE INDEX IF NOT EXISTS idx_eventos_type ON eventos(type);
CREATE INDEX IF NOT EXISTS idx_eventos_is_public ON eventos(is_public);
`

const migrationTesoreria = `
CREATE TABLE IF NOT EXISTS movimientos_tesoreria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('ingreso', 'egreso')),
    category VARCHAR(100),
    date DATE NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tesoreria_date ON movimientos_tesoreria(date DESC);
CREATE INDEX IF NOT EXISTS idx_tesoreria_type ON movimientos_tesoreria(type);
`

const migrationActas = `
CREATE TABLE IF NOT EXISTS actas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    meeting_date DATE NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'ordinaria' CHECK (type IN ('ordinaria', 'extraordinaria')),
    attendees_count INTEGER,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_actas_meeting_date ON actas(meeting_date DESC);
CREATE INDEX IF NOT EXISTS idx_actas_type ON actas(type);
`

const migrationDocumentos = `
CREATE TABLE IF NOT EXISTS documentos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_url VARCHAR(500),
    category VARCHAR(50) NOT NULL DEFAULT 'otro' CHECK (category IN ('reglamento', 'protocolo', 'formulario', 'otro')),
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_documentos_category ON documentos(category);
CREATE INDEX IF NOT EXISTS idx_documentos_is_public ON documentos(is_public);
`

const migrationEmergencias = `
CREATE TABLE IF NOT EXISTS emergencias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'resolved', 'expired')),
    expires_at TIMESTAMP WITH TIME ZONE,
    notify_email BOOLEAN NOT NULL DEFAULT TRUE,
    notify_push BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_emergencias_status ON emergencias(status);
CREATE INDEX IF NOT EXISTS idx_emergencias_priority ON emergencias(priority);
`

const migrationVotaciones = `
CREATE TABLE IF NOT EXISTS votaciones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'closed', 'cancelled')),
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    requires_quorum BOOLEAN NOT NULL DEFAULT FALSE,
    quorum_percentage INTEGER DEFAULT 50,
    allow_abstention BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS votacion_opciones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    votacion_id UUID NOT NULL REFERENCES votaciones(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS votos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    votacion_id UUID NOT NULL REFERENCES votaciones(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opcion_id UUID REFERENCES votacion_opciones(id) ON DELETE CASCADE,
    is_abstention BOOLEAN NOT NULL DEFAULT FALSE,
    voted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(votacion_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_votaciones_status ON votaciones(status);
CREATE INDEX IF NOT EXISTS idx_votos_votacion ON votos(votacion_id);
`

const migrationGalerias = `
CREATE TABLE IF NOT EXISTS galerias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date DATE,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    cover_image_url VARCHAR(500),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS galeria_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    galeria_id UUID NOT NULL REFERENCES galerias(id) ON DELETE CASCADE,
    file_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    file_type VARCHAR(50) NOT NULL,
    caption VARCHAR(500),
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_galerias_public ON galerias(is_public);
CREATE INDEX IF NOT EXISTS idx_galeria_items_galeria ON galeria_items(galeria_id);
`

const migrationMapaPuntos = `
CREATE TABLE IF NOT EXISTS mapa_puntos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    icon VARCHAR(50) DEFAULT 'marker',
    type VARCHAR(50) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mapa_puntos_type ON mapa_puntos(type);
`

const migrationMensajesContacto = `
CREATE TABLE IF NOT EXISTS mensajes_contacto (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    asunto VARCHAR(255) NOT NULL,
    mensaje TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'read', 'replied', 'archived')),
    read_at TIMESTAMP WITH TIME ZONE,
    read_by UUID REFERENCES users(id),
    replied_at TIMESTAMP WITH TIME ZONE,
    replied_by UUID REFERENCES users(id),
    respuesta TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mensajes_status ON mensajes_contacto(status);
CREATE INDEX IF NOT EXISTS idx_mensajes_user ON mensajes_contacto(user_id);
`

const migrationNotificaciones = `
CREATE TABLE IF NOT EXISTS notificaciones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    reference_id UUID,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notificaciones_user ON notificaciones(user_id, is_read);
`
