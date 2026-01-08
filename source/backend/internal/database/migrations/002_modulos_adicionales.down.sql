-- ============================================
-- ROLLBACK 002: Módulos Adicionales
-- ============================================

-- Eliminar triggers
DROP TRIGGER IF EXISTS update_mensajes_contacto_updated_at ON mensajes_contacto;
DROP TRIGGER IF EXISTS update_gastos_comunes_updated_at ON gastos_comunes;
DROP TRIGGER IF EXISTS update_periodos_gasto_updated_at ON periodos_gasto;
DROP TRIGGER IF EXISTS update_galerias_updated_at ON galerias;
DROP TRIGGER IF EXISTS update_emergencias_updated_at ON emergencias;
DROP TRIGGER IF EXISTS update_votaciones_updated_at ON votaciones;

-- Eliminar tablas (orden inverso por dependencias)
DROP TABLE IF EXISTS notificaciones;
DROP TABLE IF EXISTS mensajes_contacto;
DROP TABLE IF EXISTS pagos;
DROP TABLE IF EXISTS gastos_comunes;
DROP TABLE IF EXISTS periodos_gasto;
DROP TABLE IF EXISTS galeria_items;
DROP TABLE IF EXISTS galerias;
DROP TABLE IF EXISTS emergencias;
DROP TABLE IF EXISTS votos;
DROP TABLE IF EXISTS votacion_opciones;
DROP TABLE IF EXISTS votaciones;

-- Eliminar tipos enumerados
DROP TYPE IF EXISTS contacto_status;
DROP TYPE IF EXISTS pago_status;
DROP TYPE IF EXISTS emergencia_status;
DROP TYPE IF EXISTS emergencia_priority;
DROP TYPE IF EXISTS votacion_status;
