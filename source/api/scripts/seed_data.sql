-- ============================================
-- DATOS DE PRUEBA - Condominio Viña Pelvin
-- ============================================

-- Limpiar datos existentes (opcional)
TRUNCATE TABLE documentos, actas, movimientos_tesoreria, eventos, comunicados, users CASCADE;

-- ============================================
-- USUARIOS
-- ============================================
INSERT INTO users (id, email, password_hash, name, role, parcela_id, email_verified) VALUES
-- Administrador
('a0000000-0000-0000-0000-000000000001', 'admin@vinapelvin.cl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Carlos Mendoza', 'admin', NULL, true),

-- Directiva
('a0000000-0000-0000-0000-000000000002', 'presidente@vinapelvin.cl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'María González', 'directiva', 1, true),
('a0000000-0000-0000-0000-000000000003', 'tesorero@vinapelvin.cl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Roberto Silva', 'directiva', 2, true),
('a0000000-0000-0000-0000-000000000004', 'secretaria@vinapelvin.cl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Ana Martínez', 'directiva', 3, true),

-- Vecinos
('a0000000-0000-0000-0000-000000000005', 'juan.perez@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Juan Pérez', 'vecino', 4, true),
('a0000000-0000-0000-0000-000000000006', 'laura.rojas@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Laura Rojas', 'vecino', 5, true),
('a0000000-0000-0000-0000-000000000007', 'pedro.castro@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Pedro Castro', 'vecino', 6, true),
('a0000000-0000-0000-0000-000000000008', 'carmen.lopez@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Carmen López', 'vecino', 7, true),
('a0000000-0000-0000-0000-000000000009', 'miguel.fuentes@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Miguel Fuentes', 'vecino', 8, true),
('a0000000-0000-0000-0000-000000000010', 'sofia.herrera@email.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.SfXlJVHlF2OQqmgVvO', 'Sofía Herrera', 'vecino', 9, true);

-- ============================================
-- COMUNICADOS
-- ============================================
INSERT INTO comunicados (title, content, type, is_public, author_id, published_at) VALUES
-- Comunicados públicos
('Bienvenidos al nuevo portal de la comunidad',
'Estimados vecinos,

Nos complace anunciar el lanzamiento de nuestro nuevo portal web para la comunidad Viña Pelvin. A través de esta plataforma podrán:

- Revisar comunicados y noticias
- Consultar el calendario de eventos
- Ver el estado de tesorería
- Acceder a documentos importantes
- Participar en votaciones

¡Esperamos que esta herramienta facilite la comunicación entre todos!

Saludos cordiales,
La Directiva',
'informativo', true, 'a0000000-0000-0000-0000-000000000002', NOW() - INTERVAL '30 days'),

('Mantención de áreas verdes - Enero 2026',
'Informamos a la comunidad que durante la próxima semana se realizará mantención de las áreas verdes comunes. Los trabajos incluyen:

- Poda de árboles y arbustos
- Riego profundo
- Aplicación de fertilizantes
- Reparación del sistema de riego automático

Se solicita mantener mascotas alejadas de las áreas de trabajo durante el horario de 8:00 a 17:00 hrs.

Agradecemos su comprensión.',
'informativo', true, 'a0000000-0000-0000-0000-000000000002', NOW() - INTERVAL '15 days'),

('Corte de agua programado',
'AVISO IMPORTANTE

El día miércoles 15 de enero se realizará un corte de agua programado para reparaciones en la matriz principal. El corte será entre las 09:00 y las 14:00 hrs aproximadamente.

Se recomienda almacenar agua con anticipación para uso durante ese período.

Disculpen las molestias.',
'urgente', true, 'a0000000-0000-0000-0000-000000000003', NOW() - INTERVAL '10 days'),

('Resultados Asamblea Ordinaria 2025',
'Estimados vecinos,

Compartimos un resumen de los acuerdos alcanzados en la Asamblea Ordinaria realizada el pasado sábado:

1. Se aprobó el presupuesto 2026 por unanimidad
2. Se acordó aumentar el fondo de reserva en un 10%
3. Se aprobó el proyecto de mejoramiento de luminarias
4. Se eligió nueva directiva para el período 2026-2028

El acta completa estará disponible en la sección de documentos.

Gracias por su participación.',
'informativo', true, 'a0000000-0000-0000-0000-000000000004', NOW() - INTERVAL '5 days'),

('Recordatorio: Pago de gastos comunes',
'Estimados vecinos,

Les recordamos que el plazo para el pago de gastos comunes del mes de enero vence el día 10.

Pueden realizar el pago mediante:
- Transferencia bancaria
- Depósito en cuenta corriente
- Pago presencial en tesorería (martes y jueves de 18:00 a 20:00)

Evite recargos pagando a tiempo.

Tesorería',
'recordatorio', true, 'a0000000-0000-0000-0000-000000000003', NOW() - INTERVAL '2 days'),

('Normas de convivencia - Ruidos molestos',
'Recordamos a todos los residentes que según el reglamento interno:

- El horario de silencio es de 22:00 a 08:00 hrs en días de semana
- Fines de semana y festivos: 00:00 a 09:00 hrs
- Los trabajos de construcción solo están permitidos de lunes a viernes de 09:00 a 18:00 hrs
- El uso de herramientas ruidosas los sábados está permitido solo de 10:00 a 14:00 hrs

Agradecemos respetar estas normas para una mejor convivencia.',
'informativo', true, 'a0000000-0000-0000-0000-000000000002', NOW() - INTERVAL '1 day');

-- ============================================
-- EVENTOS
-- ============================================
INSERT INTO eventos (title, description, event_date, event_end_date, location, type, is_public, created_by) VALUES
-- Eventos pasados
('Asamblea Ordinaria 2025',
'Asamblea ordinaria anual para revisar gestión del año, aprobar presupuesto y elegir nueva directiva.',
NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days' + INTERVAL '3 hours',
'Sede Social', 'reunion', true, 'a0000000-0000-0000-0000-000000000002'),

('Taller de Jardinería',
'Taller práctico sobre cuidado de jardines y plantas de temporada. Abierto a todos los vecinos.',
NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days' + INTERVAL '2 hours',
'Área verde central', 'social', true, 'a0000000-0000-0000-0000-000000000004'),

-- Eventos futuros
('Reunión de Directiva',
'Reunión mensual de la directiva para revisión de temas pendientes y planificación.',
NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days' + INTERVAL '2 hours',
'Sede Social', 'reunion', false, 'a0000000-0000-0000-0000-000000000002'),

('Celebración Día del Niño',
'Actividades recreativas para los niños de la comunidad: juegos inflables, pintacaritas, show de magia y sorpresas.',
NOW() + INTERVAL '15 days', NOW() + INTERVAL '15 days' + INTERVAL '4 hours',
'Plaza Central', 'social', true, 'a0000000-0000-0000-0000-000000000004'),

('Mantención Piscina Comunitaria',
'Cierre temporal de la piscina por mantención anual. Limpieza profunda y revisión de equipos.',
NOW() + INTERVAL '20 days', NOW() + INTERVAL '22 days',
'Piscina', 'mantencion', true, 'a0000000-0000-0000-0000-000000000003'),

('Asamblea Extraordinaria - Proyecto Seguridad',
'Votación del proyecto de mejoramiento del sistema de seguridad perimetral.',
NOW() + INTERVAL '30 days', NOW() + INTERVAL '30 days' + INTERVAL '2 hours',
'Sede Social', 'reunion', true, 'a0000000-0000-0000-0000-000000000002'),

('Feria Navideña',
'Feria de manualidades y productos locales. Vecinos interesados en participar deben inscribirse con la directiva.',
NOW() + INTERVAL '45 days', NOW() + INTERVAL '45 days' + INTERVAL '6 hours',
'Explanada Principal', 'social', true, 'a0000000-0000-0000-0000-000000000004');

-- ============================================
-- MOVIMIENTOS TESORERÍA
-- ============================================
INSERT INTO movimientos_tesoreria (description, amount, type, category, date, created_by) VALUES
-- Ingresos
('Gastos comunes Octubre 2025', 4500000.00, 'ingreso', 'gastos_comunes', '2025-10-15', 'a0000000-0000-0000-0000-000000000003'),
('Gastos comunes Noviembre 2025', 4650000.00, 'ingreso', 'gastos_comunes', '2025-11-15', 'a0000000-0000-0000-0000-000000000003'),
('Gastos comunes Diciembre 2025', 4580000.00, 'ingreso', 'gastos_comunes', '2025-12-15', 'a0000000-0000-0000-0000-000000000003'),
('Arriendo sede social - Evento privado', 150000.00, 'ingreso', 'arriendos', '2025-11-20', 'a0000000-0000-0000-0000-000000000003'),
('Multas por atraso en pagos', 85000.00, 'ingreso', 'multas', '2025-12-10', 'a0000000-0000-0000-0000-000000000003'),
('Gastos comunes Enero 2026 (parcial)', 2100000.00, 'ingreso', 'gastos_comunes', '2026-01-05', 'a0000000-0000-0000-0000-000000000003'),

-- Egresos
('Servicio de aseo áreas comunes - Octubre', 450000.00, 'egreso', 'servicios', '2025-10-05', 'a0000000-0000-0000-0000-000000000003'),
('Electricidad áreas comunes - Octubre', 380000.00, 'egreso', 'servicios_basicos', '2025-10-10', 'a0000000-0000-0000-0000-000000000003'),
('Mantención jardines - Octubre', 280000.00, 'egreso', 'mantencion', '2025-10-12', 'a0000000-0000-0000-0000-000000000003'),
('Servicio de aseo áreas comunes - Noviembre', 450000.00, 'egreso', 'servicios', '2025-11-05', 'a0000000-0000-0000-0000-000000000003'),
('Electricidad áreas comunes - Noviembre', 420000.00, 'egreso', 'servicios_basicos', '2025-11-10', 'a0000000-0000-0000-0000-000000000003'),
('Reparación portón principal', 650000.00, 'egreso', 'reparaciones', '2025-11-18', 'a0000000-0000-0000-0000-000000000003'),
('Servicio de aseo áreas comunes - Diciembre', 450000.00, 'egreso', 'servicios', '2025-12-05', 'a0000000-0000-0000-0000-000000000003'),
('Electricidad áreas comunes - Diciembre', 520000.00, 'egreso', 'servicios_basicos', '2025-12-10', 'a0000000-0000-0000-0000-000000000003'),
('Evento navideño comunidad', 380000.00, 'egreso', 'eventos', '2025-12-20', 'a0000000-0000-0000-0000-000000000003'),
('Honorarios administración - 4to trim', 900000.00, 'egreso', 'administracion', '2025-12-28', 'a0000000-0000-0000-0000-000000000003'),
('Seguro áreas comunes 2026', 1200000.00, 'egreso', 'seguros', '2026-01-02', 'a0000000-0000-0000-0000-000000000003'),
('Servicio de aseo áreas comunes - Enero', 450000.00, 'egreso', 'servicios', '2026-01-05', 'a0000000-0000-0000-0000-000000000003');

-- ============================================
-- ACTAS
-- ============================================
INSERT INTO actas (title, content, meeting_date, type, attendees_count, created_by) VALUES
('Acta Asamblea Ordinaria - Diciembre 2025',
'ACTA DE ASAMBLEA ORDINARIA
Comunidad Viña Pelvin

Fecha: 15 de Diciembre de 2025
Hora inicio: 10:00 hrs
Lugar: Sede Social
Asistentes: 45 propietarios (75% del total)

TABLA:
1. Cuenta del Presidente
2. Balance Financiero 2025
3. Presupuesto 2026
4. Elección de Directiva
5. Varios

DESARROLLO:

1. CUENTA DEL PRESIDENTE
La presidenta María González presenta resumen de gestión 2025:
- Mejoras en áreas verdes completadas
- Sistema de seguridad actualizado
- Reducción de morosidad en 15%

2. BALANCE FINANCIERO
El tesorero Roberto Silva presenta:
- Ingresos totales: $54.200.000
- Egresos totales: $48.500.000
- Saldo positivo: $5.700.000
- Fondo de reserva: $12.000.000

Se aprueba balance por unanimidad.

3. PRESUPUESTO 2026
Se presenta presupuesto con aumento de 5% en gastos comunes.
Votación: 40 a favor, 3 en contra, 2 abstenciones.
SE APRUEBA.

4. ELECCIÓN DIRECTIVA
Se presentan candidatos y se procede a votación:
- Presidente: María González (reelecta)
- Tesorero: Roberto Silva (reelecto)
- Secretaria: Ana Martínez (nueva)

5. VARIOS
- Se acuerda estudiar proyecto de luminarias LED
- Se designa comité para fiestas patrias 2026

Sin más que tratar, se cierra la sesión a las 12:30 hrs.

Firma Secretaria: Ana Martínez
Firma Presidente: María González',
'2025-12-15', 'ordinaria', 45, 'a0000000-0000-0000-0000-000000000004'),

('Acta Reunión Directiva - Noviembre 2025',
'ACTA REUNIÓN DE DIRECTIVA
Comunidad Viña Pelvin

Fecha: 10 de Noviembre de 2025
Asistentes: Directiva completa

TEMAS TRATADOS:

1. Revisión estado de morosidad
   - 8 parcelas con pagos pendientes
   - Se acordó enviar cartas de cobranza

2. Preparación Asamblea Ordinaria
   - Se fija fecha: 15 de diciembre
   - Se prepara documentación

3. Mantención portón principal
   - Se aprueba reparación urgente
   - Presupuesto: $650.000

4. Evento navideño
   - Se aprueba presupuesto de $380.000
   - Fecha: 20 de diciembre

Próxima reunión: 8 de diciembre

Firma: Ana Martínez, Secretaria',
'2025-11-10', 'ordinaria', 4, 'a0000000-0000-0000-0000-000000000004'),

('Acta Asamblea Extraordinaria - Seguridad',
'ACTA DE ASAMBLEA EXTRAORDINARIA
Comunidad Viña Pelvin

Fecha: 20 de Septiembre de 2025
Hora: 19:00 hrs
Lugar: Sede Social
Asistentes: 38 propietarios (63% del total)

TEMA ÚNICO: Proyecto de mejoramiento sistema de seguridad

DESARROLLO:
Se presenta propuesta de actualización del sistema de seguridad que incluye:
- 12 cámaras nuevas de alta definición
- Sistema de grabación 30 días
- Control de acceso vehicular automatizado
- App para residentes

Inversión total: $8.500.000
Financiamiento propuesto: Fondo de reserva + cuota extraordinaria

Votación:
- A favor: 32
- En contra: 4
- Abstención: 2

SE APRUEBA EL PROYECTO.

Se acuerda implementación en enero 2026.

Firma: Secretaria anterior
Firma Presidente: María González',
'2025-09-20', 'extraordinaria', 38, 'a0000000-0000-0000-0000-000000000004');

-- ============================================
-- DOCUMENTOS
-- ============================================
INSERT INTO documentos (title, description, file_url, category, is_public, created_by) VALUES
('Reglamento Interno de Copropiedad',
'Reglamento vigente que establece las normas de convivencia y administración de la comunidad.',
'/documentos/reglamento-interno-2024.pdf', 'reglamento', true, 'a0000000-0000-0000-0000-000000000002'),

('Reglamento de Piscina',
'Normas de uso de la piscina comunitaria, horarios y medidas de seguridad.',
'/documentos/reglamento-piscina.pdf', 'reglamento', true, 'a0000000-0000-0000-0000-000000000002'),

('Protocolo de Emergencias',
'Procedimientos ante emergencias: incendio, sismo, cortes de suministros.',
'/documentos/protocolo-emergencias.pdf', 'protocolo', true, 'a0000000-0000-0000-0000-000000000002'),

('Protocolo de Mudanzas',
'Procedimiento y horarios permitidos para realizar mudanzas.',
'/documentos/protocolo-mudanzas.pdf', 'protocolo', true, 'a0000000-0000-0000-0000-000000000004'),

('Formulario Solicitud de Obras',
'Formulario para solicitar autorización de modificaciones en la propiedad.',
'/documentos/formulario-obras.pdf', 'formulario', true, 'a0000000-0000-0000-0000-000000000004'),

('Formulario Arriendo Sede Social',
'Solicitud y condiciones para el arriendo de la sede social.',
'/documentos/formulario-arriendo-sede.pdf', 'formulario', true, 'a0000000-0000-0000-0000-000000000004'),

('Plano General del Condominio',
'Plano con ubicación de parcelas, áreas comunes y servicios.',
'/documentos/plano-condominio.pdf', 'otro', true, 'a0000000-0000-0000-0000-000000000002'),

('Presupuesto 2026',
'Presupuesto anual aprobado en asamblea ordinaria.',
'/documentos/presupuesto-2026.pdf', 'otro', false, 'a0000000-0000-0000-0000-000000000003');

-- ============================================
-- FIN DE DATOS DE PRUEBA
-- ============================================
