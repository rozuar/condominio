-- ============================================
-- ROLLBACK 003: Seed de Parcelas
-- ============================================

-- Eliminar usuario admin
DELETE FROM users WHERE email = 'admin@vinapelvin.cl';

-- Eliminar todas las parcelas
DELETE FROM parcelas;
