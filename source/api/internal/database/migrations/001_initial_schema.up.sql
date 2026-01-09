-- ============================================
-- MIGRACIÓN 001: Schema Inicial
-- Comunidad Viña Pelvin
-- ============================================

-- Extensiones
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- TIPOS ENUMERADOS
-- ============================================

CREATE TYPE user_role AS ENUM ('visitor', 'vecino', 'directiva');
CREATE TYPE comunicado_type AS ENUM ('informativo', 'seguridad', 'tesoreria', 'asamblea');
CREATE TYPE evento_type AS ENUM ('reunion', 'asamblea', 'trabajo', 'social');
CREATE TYPE movimiento_type AS ENUM ('ingreso', 'egreso');
CREATE TYPE acta_type AS ENUM ('ordinaria', 'extraordinaria');
CREATE TYPE documento_category AS ENUM ('reglamento', 'protocolo', 'formulario', 'otro');

-- ============================================
-- TABLAS BASE
-- ============================================

-- Parcelas (73 parcelas en la comunidad)
CREATE TABLE parcelas (
    id SERIAL PRIMARY KEY,
    numero VARCHAR(10) NOT NULL UNIQUE,
    direccion VARCHAR(255),
    superficie_m2 DECIMAL(10, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'vecino',
    parcela_id INTEGER REFERENCES parcelas(id) ON DELETE SET NULL,
    google_id VARCHAR(255) UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_parcela ON users(parcela_id);

-- ============================================
-- COMUNICADOS
-- ============================================

CREATE TABLE comunicados (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type comunicado_type NOT NULL DEFAULT 'informativo',
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comunicados_type ON comunicados(type);
CREATE INDEX idx_comunicados_public ON comunicados(is_public);
CREATE INDEX idx_comunicados_published ON comunicados(published_at DESC);

-- ============================================
-- EVENTOS / CALENDARIO
-- ============================================

CREATE TABLE eventos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMPTZ NOT NULL,
    event_end_date TIMESTAMPTZ,
    location VARCHAR(255),
    type evento_type NOT NULL DEFAULT 'reunion',
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_eventos_date ON eventos(event_date);
CREATE INDEX idx_eventos_type ON eventos(type);
CREATE INDEX idx_eventos_public ON eventos(is_public);

-- ============================================
-- TESORERÍA / MOVIMIENTOS
-- ============================================

CREATE TABLE movimientos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    description VARCHAR(500) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    type movimiento_type NOT NULL,
    category VARCHAR(100),
    date DATE NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movimientos_type ON movimientos(type);
CREATE INDEX idx_movimientos_date ON movimientos(date DESC);
CREATE INDEX idx_movimientos_year_month ON movimientos(EXTRACT(YEAR FROM date), EXTRACT(MONTH FROM date));

-- ============================================
-- ACTAS
-- ============================================

CREATE TABLE actas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    meeting_date DATE NOT NULL,
    type acta_type NOT NULL DEFAULT 'ordinaria',
    attendees_count INTEGER,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_actas_date ON actas(meeting_date DESC);
CREATE INDEX idx_actas_type ON actas(type);

-- ============================================
-- DOCUMENTOS
-- ============================================

CREATE TABLE documentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_url VARCHAR(500) NOT NULL,
    file_size INTEGER,
    mime_type VARCHAR(100),
    category documento_category NOT NULL DEFAULT 'otro',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documentos_category ON documentos(category);
CREATE INDEX idx_documentos_public ON documentos(is_public);

-- ============================================
-- FUNCIONES DE ACTUALIZACIÓN
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers para updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_parcelas_updated_at BEFORE UPDATE ON parcelas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comunicados_updated_at BEFORE UPDATE ON comunicados
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_eventos_updated_at BEFORE UPDATE ON eventos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_movimientos_updated_at BEFORE UPDATE ON movimientos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_actas_updated_at BEFORE UPDATE ON actas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_documentos_updated_at BEFORE UPDATE ON documentos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
