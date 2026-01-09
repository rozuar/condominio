-- ============================================
-- MIGRACIÓN 004: Mapa de la Comunidad
-- Estructura para mapa interactivo
-- ============================================

CREATE TYPE area_type AS ENUM ('parcela', 'area_comun', 'acceso', 'canal', 'camino');

-- Áreas del mapa (parcelas, áreas comunes, etc.)
CREATE TABLE mapa_areas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parcela_id INTEGER REFERENCES parcelas(id) ON DELETE CASCADE,
    type area_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    -- Coordenadas del polígono (GeoJSON format stored as JSONB)
    coordinates JSONB NOT NULL,
    -- Centro del área para marcadores
    center_lat DECIMAL(10, 8),
    center_lng DECIMAL(11, 8),
    -- Estilo visual
    fill_color VARCHAR(7) DEFAULT '#4A7C23',
    stroke_color VARCHAR(7) DEFAULT '#2D5016',
    is_clickable BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Puntos de interés en el mapa
CREATE TABLE mapa_puntos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    icon VARCHAR(50) DEFAULT 'marker',
    type VARCHAR(50) NOT NULL, -- 'entrada', 'sede', 'plaza', 'bomba_agua', etc.
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mapa_areas_type ON mapa_areas(type);
CREATE INDEX idx_mapa_areas_parcela ON mapa_areas(parcela_id);
CREATE INDEX idx_mapa_puntos_type ON mapa_puntos(type);

CREATE TRIGGER update_mapa_areas_updated_at BEFORE UPDATE ON mapa_areas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_mapa_puntos_updated_at BEFORE UPDATE ON mapa_puntos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
