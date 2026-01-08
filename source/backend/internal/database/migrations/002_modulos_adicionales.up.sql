-- ============================================
-- MIGRACIÓN 002: Módulos Adicionales
-- Votaciones, Emergencias, Galería, Gastos Comunes, Contacto
-- ============================================

-- ============================================
-- TIPOS ENUMERADOS ADICIONALES
-- ============================================

CREATE TYPE votacion_status AS ENUM ('draft', 'active', 'closed', 'cancelled');
CREATE TYPE emergencia_priority AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE emergencia_status AS ENUM ('active', 'resolved', 'expired');
CREATE TYPE pago_status AS ENUM ('pending', 'paid', 'overdue', 'cancelled');
CREATE TYPE contacto_status AS ENUM ('pending', 'read', 'replied', 'archived');

-- ============================================
-- VOTACIONES
-- ============================================

CREATE TABLE votaciones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status votacion_status NOT NULL DEFAULT 'draft',
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    requires_quorum BOOLEAN NOT NULL DEFAULT FALSE,
    quorum_percentage INTEGER DEFAULT 50,
    allow_abstention BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE votacion_opciones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    votacion_id UUID NOT NULL REFERENCES votaciones(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE votos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    votacion_id UUID NOT NULL REFERENCES votaciones(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opcion_id UUID REFERENCES votacion_opciones(id) ON DELETE CASCADE,
    is_abstention BOOLEAN NOT NULL DEFAULT FALSE,
    voted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(votacion_id, user_id)
);

CREATE INDEX idx_votaciones_status ON votaciones(status);
CREATE INDEX idx_votaciones_dates ON votaciones(start_date, end_date);
CREATE INDEX idx_votos_votacion ON votos(votacion_id);

-- ============================================
-- EMERGENCIAS
-- ============================================

CREATE TABLE emergencias (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    priority emergencia_priority NOT NULL DEFAULT 'medium',
    status emergencia_status NOT NULL DEFAULT 'active',
    expires_at TIMESTAMPTZ,
    notify_email BOOLEAN NOT NULL DEFAULT TRUE,
    notify_push BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emergencias_status ON emergencias(status);
CREATE INDEX idx_emergencias_priority ON emergencias(priority);
CREATE INDEX idx_emergencias_active ON emergencias(status, created_at DESC) WHERE status = 'active';

-- ============================================
-- GALERÍA MULTIMEDIA
-- ============================================

CREATE TABLE galerias (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date DATE,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    cover_image_url VARCHAR(500),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE galeria_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    galeria_id UUID NOT NULL REFERENCES galerias(id) ON DELETE CASCADE,
    file_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    file_type VARCHAR(50) NOT NULL, -- 'image' or 'video'
    caption VARCHAR(500),
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_galerias_public ON galerias(is_public);
CREATE INDEX idx_galerias_date ON galerias(event_date DESC);
CREATE INDEX idx_galeria_items_galeria ON galeria_items(galeria_id);

-- ============================================
-- GASTOS COMUNES
-- ============================================

CREATE TABLE periodos_gasto (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    monto_base DECIMAL(12, 2) NOT NULL,
    fecha_vencimiento DATE NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(year, month)
);

CREATE TABLE gastos_comunes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    periodo_id UUID NOT NULL REFERENCES periodos_gasto(id) ON DELETE CASCADE,
    parcela_id INTEGER NOT NULL REFERENCES parcelas(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    monto DECIMAL(12, 2) NOT NULL,
    monto_pagado DECIMAL(12, 2) NOT NULL DEFAULT 0,
    status pago_status NOT NULL DEFAULT 'pending',
    fecha_pago TIMESTAMPTZ,
    metodo_pago VARCHAR(50),
    referencia_pago VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(periodo_id, parcela_id)
);

CREATE TABLE pagos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    gasto_comun_id UUID NOT NULL REFERENCES gastos_comunes(id) ON DELETE CASCADE,
    monto DECIMAL(12, 2) NOT NULL,
    metodo VARCHAR(50) NOT NULL, -- 'transbank', 'mercadopago', 'transferencia', 'efectivo'
    referencia_externa VARCHAR(255),
    estado VARCHAR(50) NOT NULL, -- 'pending', 'approved', 'rejected'
    detalles JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gastos_periodo ON gastos_comunes(periodo_id);
CREATE INDEX idx_gastos_parcela ON gastos_comunes(parcela_id);
CREATE INDEX idx_gastos_status ON gastos_comunes(status);
CREATE INDEX idx_pagos_gasto ON pagos(gasto_comun_id);

-- ============================================
-- CONTACTO CON DIRECTIVA
-- ============================================

CREATE TABLE mensajes_contacto (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    asunto VARCHAR(255) NOT NULL,
    mensaje TEXT NOT NULL,
    status contacto_status NOT NULL DEFAULT 'pending',
    read_at TIMESTAMPTZ,
    read_by UUID REFERENCES users(id) ON DELETE SET NULL,
    replied_at TIMESTAMPTZ,
    replied_by UUID REFERENCES users(id) ON DELETE SET NULL,
    respuesta TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mensajes_status ON mensajes_contacto(status);
CREATE INDEX idx_mensajes_user ON mensajes_contacto(user_id);

-- ============================================
-- NOTIFICACIONES
-- ============================================

CREATE TABLE notificaciones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'comunicado', 'emergencia', 'votacion', 'pago', etc.
    reference_id UUID, -- ID del objeto relacionado
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notificaciones_user ON notificaciones(user_id, is_read);
CREATE INDEX idx_notificaciones_unread ON notificaciones(user_id, created_at DESC) WHERE is_read = FALSE;

-- ============================================
-- TRIGGERS ADICIONALES
-- ============================================

CREATE TRIGGER update_votaciones_updated_at BEFORE UPDATE ON votaciones
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_emergencias_updated_at BEFORE UPDATE ON emergencias
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_galerias_updated_at BEFORE UPDATE ON galerias
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_periodos_gasto_updated_at BEFORE UPDATE ON periodos_gasto
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_gastos_comunes_updated_at BEFORE UPDATE ON gastos_comunes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_mensajes_contacto_updated_at BEFORE UPDATE ON mensajes_contacto
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
