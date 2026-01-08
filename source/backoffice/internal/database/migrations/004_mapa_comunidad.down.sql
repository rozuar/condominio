-- ============================================
-- ROLLBACK 004: Mapa de la Comunidad
-- ============================================

DROP TRIGGER IF EXISTS update_mapa_puntos_updated_at ON mapa_puntos;
DROP TRIGGER IF EXISTS update_mapa_areas_updated_at ON mapa_areas;

DROP TABLE IF EXISTS mapa_puntos;
DROP TABLE IF EXISTS mapa_areas;

DROP TYPE IF EXISTS area_type;
