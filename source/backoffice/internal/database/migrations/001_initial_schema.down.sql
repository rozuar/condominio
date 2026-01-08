-- ============================================
-- ROLLBACK 001: Schema Inicial
-- ============================================

-- Eliminar triggers
DROP TRIGGER IF EXISTS update_documentos_updated_at ON documentos;
DROP TRIGGER IF EXISTS update_actas_updated_at ON actas;
DROP TRIGGER IF EXISTS update_movimientos_updated_at ON movimientos;
DROP TRIGGER IF EXISTS update_eventos_updated_at ON eventos;
DROP TRIGGER IF EXISTS update_comunicados_updated_at ON comunicados;
DROP TRIGGER IF EXISTS update_parcelas_updated_at ON parcelas;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Eliminar función
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Eliminar tablas (orden inverso por dependencias)
DROP TABLE IF EXISTS documentos;
DROP TABLE IF EXISTS actas;
DROP TABLE IF EXISTS movimientos;
DROP TABLE IF EXISTS eventos;
DROP TABLE IF EXISTS comunicados;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS parcelas;

-- Eliminar tipos enumerados
DROP TYPE IF EXISTS documento_category;
DROP TYPE IF EXISTS acta_type;
DROP TYPE IF EXISTS movimiento_type;
DROP TYPE IF EXISTS evento_type;
DROP TYPE IF EXISTS comunicado_type;
DROP TYPE IF EXISTS user_role;
