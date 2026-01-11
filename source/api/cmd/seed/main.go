package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/condominio/backend/internal/database"
)

func main() {
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	db, err := database.New(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Ensure schema exists (seed can be run on a fresh DB)
	log.Println("Running migrations...")
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	pool := db.Pool

	log.Println("Connected to database, seeding data...")

	// Clear existing data
	log.Println("Clearing existing data...")
	clearTables := []string{
		"pagos",
		"gastos_comunes",
		"periodos_gasto",
		"notificaciones",
		"mensajes_contacto",
		"galeria_items",
		"galerias",
		"votos",
		"votacion_opciones",
		"votaciones",
		"emergencias",
		"mapa_puntos",
		"mapa_areas",
		"documentos",
		"actas",
		"movimientos_tesoreria",
		"eventos",
		"comunicados",
		"users",
		"parcelas",
	}
	for _, table := range clearTables {
		_, err = pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			log.Printf("Warning clearing %s: %v", table, err)
		}
	}

	// ============================================
	// PARCELAS (requerido para Gastos Comunes)
	// ============================================
	log.Println("Inserting parcelas...")
	_, err = pool.Exec(ctx, `
		INSERT INTO parcelas (id, numero, direccion) VALUES
		(1, '1', 'Calle Los Álamos 101'),
		(2, '2', 'Calle Los Álamos 102'),
		(3, '3', 'Calle Los Álamos 103'),
		(4, '4', 'Calle Los Álamos 104'),
		(5, '5', 'Calle Los Álamos 105'),
		(6, '6', 'Calle Los Álamos 106'),
		(7, '7', 'Calle Los Álamos 107'),
		(8, '8', 'Calle Los Álamos 108'),
		(9, '9', 'Calle Los Álamos 109'),
		(10, '10', 'Calle Los Álamos 110'),
		(11, '11', 'Calle Los Álamos 111'),
		(12, '12', 'Calle Los Álamos 112')
	`)
	if err != nil {
		log.Printf("Warning inserting parcelas: %v", err)
	}

	// ============================================
	// USERS
	// ============================================
	log.Println("Inserting users...")
	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, name, role, parcela_id, email_verified) VALUES
		('a0000000-0000-0000-0000-000000000001', 'admin@vinapelvin.cl', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Carlos Mendoza', 'admin', 1, true),
		('a0000000-0000-0000-0000-000000000002', 'presidente@vinapelvin.cl', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'María González', 'directiva', 1, true),
		('a0000000-0000-0000-0000-000000000003', 'tesorero@vinapelvin.cl', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Roberto Silva', 'directiva', 2, true),
		('a0000000-0000-0000-0000-000000000004', 'secretaria@vinapelvin.cl', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Ana Martínez', 'directiva', 3, true),
		('a0000000-0000-0000-0000-000000000005', 'juan.perez@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Juan Pérez', 'vecino', 4, true),
		('a0000000-0000-0000-0000-000000000011', 'familia.guest@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Invitado Familia', 'familia', NULL, true),
		('a0000000-0000-0000-0000-000000000012', 'familia.parcela4@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Familia Parcela 4', 'familia', 4, true),
		('a0000000-0000-0000-0000-000000000006', 'laura.rojas@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Laura Rojas', 'vecino', 5, true),
		('a0000000-0000-0000-0000-000000000007', 'pedro.castro@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Pedro Castro', 'vecino', 6, true),
		('a0000000-0000-0000-0000-000000000008', 'carmen.lopez@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Carmen López', 'vecino', 7, true),
		('a0000000-0000-0000-0000-000000000009', 'miguel.fuentes@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Miguel Fuentes', 'vecino', 8, true),
		('a0000000-0000-0000-0000-000000000010', 'sofia.herrera@email.com', '$2a$10$QQ/g64VCks/027EZYEh.XeGfr.PVTUHhG9aoEdDeg/bD/i5p3idgm', 'Sofía Herrera', 'vecino', 9, true)
	`)
	if err != nil {
		log.Fatalf("Failed to insert users: %v", err)
	}

	// ============================================
	// COMUNICADOS
	// ============================================
	log.Println("Inserting comunicados...")
	comunicados := []struct {
		title, content, tipo string
		daysAgo              int
	}{
		{"Bienvenidos al nuevo portal de la comunidad",
			"Estimados vecinos,\n\nNos complace anunciar el lanzamiento de nuestro nuevo portal web para la comunidad Viña Pelvin. A través de esta plataforma podrán revisar comunicados, consultar el calendario de eventos, ver el estado de tesorería y acceder a documentos importantes.\n\n¡Esperamos que esta herramienta facilite la comunicación entre todos!\n\nSaludos cordiales,\nLa Directiva",
			"informativo", 30},
		{"Mantención de áreas verdes - Enero 2026",
			"Informamos a la comunidad que durante la próxima semana se realizará mantención de las áreas verdes comunes. Los trabajos incluyen poda de árboles, riego profundo y reparación del sistema de riego automático.\n\nSe solicita mantener mascotas alejadas de las áreas de trabajo durante el horario de 8:00 a 17:00 hrs.",
			"informativo", 15},
		{"Corte de agua programado",
			"AVISO IMPORTANTE\n\nEl día miércoles 15 de enero se realizará un corte de agua programado para reparaciones en la matriz principal. El corte será entre las 09:00 y las 14:00 hrs aproximadamente.\n\nSe recomienda almacenar agua con anticipación.",
			"urgente", 10},
		{"Resultados Asamblea Ordinaria 2025",
			"Estimados vecinos,\n\nCompartimos un resumen de los acuerdos alcanzados en la Asamblea Ordinaria:\n\n1. Se aprobó el presupuesto 2026 por unanimidad\n2. Se acordó aumentar el fondo de reserva en un 10%\n3. Se aprobó el proyecto de mejoramiento de luminarias\n4. Se eligió nueva directiva para el período 2026-2028\n\nEl acta completa estará disponible en la sección de documentos.",
			"informativo", 5},
		{"Recordatorio: Pago de gastos comunes",
			"Estimados vecinos,\n\nLes recordamos que el plazo para el pago de gastos comunes del mes de enero vence el día 10.\n\nPueden realizar el pago mediante transferencia bancaria, depósito en cuenta corriente o pago presencial en tesorería.\n\nEvite recargos pagando a tiempo.\n\nTesorería",
			"recordatorio", 2},
		{"Normas de convivencia - Ruidos molestos",
			"Recordamos a todos los residentes que según el reglamento interno:\n\n- El horario de silencio es de 22:00 a 08:00 hrs en días de semana\n- Fines de semana y festivos: 00:00 a 09:00 hrs\n- Los trabajos de construcción solo están permitidos de lunes a viernes de 09:00 a 18:00 hrs\n\nAgradecemos respetar estas normas para una mejor convivencia.",
			"informativo", 1},
	}

	for _, c := range comunicados {
		_, err = pool.Exec(ctx, `
			INSERT INTO comunicados (title, content, type, is_public, author_id, published_at)
			VALUES ($1, $2, $3, true, 'a0000000-0000-0000-0000-000000000002', NOW() - $4::interval)
		`, c.title, c.content, c.tipo, fmt.Sprintf("%d days", c.daysAgo))
		if err != nil {
			log.Printf("Warning inserting comunicado: %v", err)
		}
	}

	// ============================================
	// EVENTOS
	// ============================================
	log.Println("Inserting eventos...")
	eventos := []struct {
		title, desc, location, tipo string
		daysOffset                  int
		hours                       int
	}{
		{"Asamblea Ordinaria 2025", "Asamblea ordinaria anual para revisar gestión del año.", "Sede Social", "reunion", -20, 3},
		{"Taller de Jardinería", "Taller práctico sobre cuidado de jardines y plantas.", "Área verde central", "social", -10, 2},
		{"Reunión de Directiva", "Reunión mensual de la directiva.", "Sede Social", "reunion", 3, 2},
		{"Celebración Día del Niño", "Actividades recreativas: juegos inflables, pintacaritas, show de magia.", "Plaza Central", "social", 15, 4},
		{"Mantención Piscina Comunitaria", "Cierre temporal por mantención anual.", "Piscina", "mantencion", 20, 48},
		{"Asamblea Extraordinaria - Proyecto Seguridad", "Votación del proyecto de mejoramiento del sistema de seguridad.", "Sede Social", "reunion", 30, 2},
		{"Feria Navideña", "Feria de manualidades y productos locales.", "Explanada Principal", "social", 45, 6},
	}

	for _, e := range eventos {
		_, err = pool.Exec(ctx, `
			INSERT INTO eventos (title, description, event_date, event_end_date, location, type, is_public, created_by)
			VALUES ($1, $2, NOW() + $3::interval, NOW() + $3::interval + $4::interval, $5, $6, true, 'a0000000-0000-0000-0000-000000000002')
		`, e.title, e.desc, fmt.Sprintf("%d days", e.daysOffset), fmt.Sprintf("%d hours", e.hours), e.location, e.tipo)
		if err != nil {
			log.Printf("Warning inserting evento: %v", err)
		}
	}

	// ============================================
	// MOVIMIENTOS TESORERIA
	// ============================================
	log.Println("Inserting movimientos tesoreria...")
	movimientos := []struct {
		desc     string
		amount   float64
		tipo     string
		category string
		date     string
	}{
		{"Gastos comunes Octubre 2025", 4500000, "ingreso", "gastos_comunes", "2025-10-15"},
		{"Gastos comunes Noviembre 2025", 4650000, "ingreso", "gastos_comunes", "2025-11-15"},
		{"Gastos comunes Diciembre 2025", 4580000, "ingreso", "gastos_comunes", "2025-12-15"},
		{"Arriendo sede social - Evento privado", 150000, "ingreso", "arriendos", "2025-11-20"},
		{"Multas por atraso en pagos", 85000, "ingreso", "multas", "2025-12-10"},
		{"Gastos comunes Enero 2026 (parcial)", 2100000, "ingreso", "gastos_comunes", "2026-01-05"},
		{"Servicio de aseo áreas comunes - Octubre", 450000, "egreso", "servicios", "2025-10-05"},
		{"Electricidad áreas comunes - Octubre", 380000, "egreso", "servicios_basicos", "2025-10-10"},
		{"Mantención jardines - Octubre", 280000, "egreso", "mantencion", "2025-10-12"},
		{"Servicio de aseo áreas comunes - Noviembre", 450000, "egreso", "servicios", "2025-11-05"},
		{"Electricidad áreas comunes - Noviembre", 420000, "egreso", "servicios_basicos", "2025-11-10"},
		{"Reparación portón principal", 650000, "egreso", "reparaciones", "2025-11-18"},
		{"Servicio de aseo áreas comunes - Diciembre", 450000, "egreso", "servicios", "2025-12-05"},
		{"Electricidad áreas comunes - Diciembre", 520000, "egreso", "servicios_basicos", "2025-12-10"},
		{"Evento navideño comunidad", 380000, "egreso", "eventos", "2025-12-20"},
		{"Honorarios administración - 4to trim", 900000, "egreso", "administracion", "2025-12-28"},
		{"Seguro áreas comunes 2026", 1200000, "egreso", "seguros", "2026-01-02"},
		{"Servicio de aseo áreas comunes - Enero", 450000, "egreso", "servicios", "2026-01-05"},
	}

	for _, m := range movimientos {
		_, err = pool.Exec(ctx, `
			INSERT INTO movimientos_tesoreria (description, amount, type, category, date, created_by)
			VALUES ($1, $2, $3, $4, $5, 'a0000000-0000-0000-0000-000000000003')
		`, m.desc, m.amount, m.tipo, m.category, m.date)
		if err != nil {
			log.Printf("Warning inserting movimiento: %v", err)
		}
	}

	// ============================================
	// GASTOS COMUNES - PERIODOS + GASTOS + PAGOS
	// ============================================
	log.Println("Inserting gastos comunes (periodos/gastos/pagos)...")
	// Periodos (incluye meses anteriores y el actual)
	_, err = pool.Exec(ctx, `
		INSERT INTO periodos_gasto (id, year, month, monto_base, fecha_vencimiento, descripcion) VALUES
		('90000000-0000-0000-0000-000000000010', 2025, 11, 65000, '2025-11-10', 'Gastos comunes Noviembre 2025'),
		('90000000-0000-0000-0000-000000000011', 2025, 12, 68000, '2025-12-10', 'Gastos comunes Diciembre 2025'),
		('90000000-0000-0000-0000-000000000012', 2026, 1, 70000, '2026-01-10', 'Gastos comunes Enero 2026'),
		('90000000-0000-0000-0000-000000000013', 2026, 2, 70000, '2026-02-10', 'Gastos comunes Febrero 2026')
	`)
	if err != nil {
		log.Printf("Warning inserting periodos_gasto: %v", err)
	}

	// Gastos por parcela (1..9) y por periodo (4)
	// Mezclamos status para tener pagados / pendientes / vencidos
	_, err = pool.Exec(ctx, `
		INSERT INTO gastos_comunes (
			id, periodo_id, parcela_id, user_id, monto, monto_pagado, status, fecha_pago, metodo_pago, referencia_pago
		) VALUES
		-- Noviembre 2025: varios pagados
		('91000000-0000-0000-0000-000000000101', '90000000-0000-0000-0000-000000000010', 1, 'a0000000-0000-0000-0000-000000000002', 65000, 65000, 'paid', NOW() - INTERVAL '55 days', 'transferencia', 'TRX-2025-11-0001'),
		('91000000-0000-0000-0000-000000000102', '90000000-0000-0000-0000-000000000010', 2, 'a0000000-0000-0000-0000-000000000003', 65000, 65000, 'paid', NOW() - INTERVAL '54 days', 'transferencia', 'TRX-2025-11-0002'),
		('91000000-0000-0000-0000-000000000103', '90000000-0000-0000-0000-000000000010', 3, 'a0000000-0000-0000-0000-000000000004', 65000, 65000, 'paid', NOW() - INTERVAL '53 days', 'transferencia', 'TRX-2025-11-0003'),
		('91000000-0000-0000-0000-000000000104', '90000000-0000-0000-0000-000000000010', 4, 'a0000000-0000-0000-0000-000000000005', 65000, 65000, 'paid', NOW() - INTERVAL '52 days', 'transferencia', 'TRX-2025-11-0004'),
		('91000000-0000-0000-0000-000000000105', '90000000-0000-0000-0000-000000000010', 5, 'a0000000-0000-0000-0000-000000000006', 65000, 0, 'overdue', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000106', '90000000-0000-0000-0000-000000000010', 6, 'a0000000-0000-0000-0000-000000000007', 65000, 65000, 'paid', NOW() - INTERVAL '50 days', 'transferencia', 'TRX-2025-11-0006'),
		('91000000-0000-0000-0000-000000000107', '90000000-0000-0000-0000-000000000010', 7, 'a0000000-0000-0000-0000-000000000008', 65000, 0, 'overdue', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000108', '90000000-0000-0000-0000-000000000010', 8, 'a0000000-0000-0000-0000-000000000009', 65000, 65000, 'paid', NOW() - INTERVAL '49 days', 'transferencia', 'TRX-2025-11-0008'),
		('91000000-0000-0000-0000-000000000109', '90000000-0000-0000-0000-000000000010', 9, 'a0000000-0000-0000-0000-000000000010', 65000, 0, 'pending', NULL, NULL, NULL),

		-- Diciembre 2025: algunos pagos parciales + vencidos
		('91000000-0000-0000-0000-000000000201', '90000000-0000-0000-0000-000000000011', 1, 'a0000000-0000-0000-0000-000000000002', 68000, 68000, 'paid', NOW() - INTERVAL '25 days', 'transferencia', 'TRX-2025-12-0001'),
		('91000000-0000-0000-0000-000000000202', '90000000-0000-0000-0000-000000000011', 2, 'a0000000-0000-0000-0000-000000000003', 68000, 68000, 'paid', NOW() - INTERVAL '24 days', 'transferencia', 'TRX-2025-12-0002'),
		('91000000-0000-0000-0000-000000000203', '90000000-0000-0000-0000-000000000011', 3, 'a0000000-0000-0000-0000-000000000004', 68000, 34000, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000204', '90000000-0000-0000-0000-000000000011', 4, 'a0000000-0000-0000-0000-000000000005', 68000, 0, 'overdue', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000205', '90000000-0000-0000-0000-000000000011', 5, 'a0000000-0000-0000-0000-000000000006', 68000, 68000, 'paid', NOW() - INTERVAL '20 days', 'transferencia', 'TRX-2025-12-0005'),
		('91000000-0000-0000-0000-000000000206', '90000000-0000-0000-0000-000000000011', 6, 'a0000000-0000-0000-0000-000000000007', 68000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000207', '90000000-0000-0000-0000-000000000011', 7, 'a0000000-0000-0000-0000-000000000008', 68000, 68000, 'paid', NOW() - INTERVAL '18 days', 'transferencia', 'TRX-2025-12-0007'),
		('91000000-0000-0000-0000-000000000208', '90000000-0000-0000-0000-000000000011', 8, 'a0000000-0000-0000-0000-000000000009', 68000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000209', '90000000-0000-0000-0000-000000000011', 9, 'a0000000-0000-0000-0000-000000000010', 68000, 68000, 'paid', NOW() - INTERVAL '16 days', 'transferencia', 'TRX-2025-12-0009'),

		-- Enero 2026: periodo actual (mixto)
		('91000000-0000-0000-0000-000000000301', '90000000-0000-0000-0000-000000000012', 1, 'a0000000-0000-0000-0000-000000000002', 70000, 70000, 'paid', NOW() - INTERVAL '5 days', 'transferencia', 'TRX-2026-01-0001'),
		('91000000-0000-0000-0000-000000000302', '90000000-0000-0000-0000-000000000012', 2, 'a0000000-0000-0000-0000-000000000003', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000303', '90000000-0000-0000-0000-000000000012', 3, 'a0000000-0000-0000-0000-000000000004', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000304', '90000000-0000-0000-0000-000000000012', 4, 'a0000000-0000-0000-0000-000000000005', 70000, 70000, 'paid', NOW() - INTERVAL '3 days', 'transferencia', 'TRX-2026-01-0004'),
		('91000000-0000-0000-0000-000000000305', '90000000-0000-0000-0000-000000000012', 5, 'a0000000-0000-0000-0000-000000000006', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000306', '90000000-0000-0000-0000-000000000012', 6, 'a0000000-0000-0000-0000-000000000007', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000307', '90000000-0000-0000-0000-000000000012', 7, 'a0000000-0000-0000-0000-000000000008', 70000, 35000, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000308', '90000000-0000-0000-0000-000000000012', 8, 'a0000000-0000-0000-0000-000000000009', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000309', '90000000-0000-0000-0000-000000000012', 9, 'a0000000-0000-0000-0000-000000000010', 70000, 70000, 'paid', NOW() - INTERVAL '2 days', 'transferencia', 'TRX-2026-01-0009'),

		-- Febrero 2026: futuro (todo pendiente)
		('91000000-0000-0000-0000-000000000401', '90000000-0000-0000-0000-000000000013', 1, 'a0000000-0000-0000-0000-000000000002', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000402', '90000000-0000-0000-0000-000000000013', 2, 'a0000000-0000-0000-0000-000000000003', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000403', '90000000-0000-0000-0000-000000000013', 3, 'a0000000-0000-0000-0000-000000000004', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000404', '90000000-0000-0000-0000-000000000013', 4, 'a0000000-0000-0000-0000-000000000005', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000405', '90000000-0000-0000-0000-000000000013', 5, 'a0000000-0000-0000-0000-000000000006', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000406', '90000000-0000-0000-0000-000000000013', 6, 'a0000000-0000-0000-0000-000000000007', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000407', '90000000-0000-0000-0000-000000000013', 7, 'a0000000-0000-0000-0000-000000000008', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000408', '90000000-0000-0000-0000-000000000013', 8, 'a0000000-0000-0000-0000-000000000009', 70000, 0, 'pending', NULL, NULL, NULL),
		('91000000-0000-0000-0000-000000000409', '90000000-0000-0000-0000-000000000013', 9, 'a0000000-0000-0000-0000-000000000010', 70000, 0, 'pending', NULL, NULL, NULL)
	`)
	if err != nil {
		log.Printf("Warning inserting gastos_comunes: %v", err)
	}

	// Pagos asociados a los "paid"
	_, err = pool.Exec(ctx, `
		INSERT INTO pagos (gasto_comun_id, monto, metodo, referencia_externa, estado, detalles) VALUES
		('91000000-0000-0000-0000-000000000101', 65000, 'transferencia', 'TRX-2025-11-0001', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000102', 65000, 'transferencia', 'TRX-2025-11-0002', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000103', 65000, 'transferencia', 'TRX-2025-11-0003', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000104', 65000, 'transferencia', 'TRX-2025-11-0004', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000106', 65000, 'transferencia', 'TRX-2025-11-0006', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000108', 65000, 'transferencia', 'TRX-2025-11-0008', 'approved', '{"banco":"Banco Estado"}'::jsonb),

		('91000000-0000-0000-0000-000000000201', 68000, 'transferencia', 'TRX-2025-12-0001', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000202', 68000, 'transferencia', 'TRX-2025-12-0002', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000205', 68000, 'transferencia', 'TRX-2025-12-0005', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000207', 68000, 'transferencia', 'TRX-2025-12-0007', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000209', 68000, 'transferencia', 'TRX-2025-12-0009', 'approved', '{"banco":"Banco Estado"}'::jsonb),

		('91000000-0000-0000-0000-000000000301', 70000, 'transferencia', 'TRX-2026-01-0001', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000304', 70000, 'transferencia', 'TRX-2026-01-0004', 'approved', '{"banco":"Banco Estado"}'::jsonb),
		('91000000-0000-0000-0000-000000000309', 70000, 'transferencia', 'TRX-2026-01-0009', 'approved', '{"banco":"Banco Estado"}'::jsonb)
	`)
	if err != nil {
		log.Printf("Warning inserting pagos (gastos comunes): %v", err)
	}

	// ============================================
	// ACTAS
	// ============================================
	log.Println("Inserting actas...")
	actas := []struct {
		title, content string
		date           string
		tipo           string
		attendees      int
	}{
		{"Acta Asamblea Ordinaria - Diciembre 2025",
			"ACTA DE ASAMBLEA ORDINARIA\nComunidad Viña Pelvin\n\nFecha: 15 de Diciembre de 2025\nAsistentes: 45 propietarios (75%)\n\nACUERDOS:\n1. Se aprobó presupuesto 2026 por unanimidad\n2. Se aprobó aumento fondo de reserva 10%\n3. Se aprobó proyecto luminarias LED\n4. Se eligió nueva directiva 2026-2028\n\nFirma: Ana Martínez, Secretaria",
			"2025-12-15", "ordinaria", 45},
		{"Acta Reunión Directiva - Noviembre 2025",
			"ACTA REUNIÓN DE DIRECTIVA\n\nFecha: 10 de Noviembre de 2025\nAsistentes: Directiva completa\n\nTEMAS:\n1. Revisión morosidad - 8 parcelas pendientes\n2. Preparación Asamblea Ordinaria\n3. Reparación portón - aprobado $650.000\n4. Evento navideño - aprobado $380.000\n\nFirma: Ana Martínez, Secretaria",
			"2025-11-10", "ordinaria", 4},
		{"Acta Asamblea Extraordinaria - Seguridad",
			"ACTA DE ASAMBLEA EXTRAORDINARIA\n\nFecha: 20 de Septiembre de 2025\nAsistentes: 38 propietarios (63%)\n\nTEMA: Proyecto de mejoramiento sistema de seguridad\n\nPropuesta: 12 cámaras HD, grabación 30 días, control acceso automatizado\nInversión: $8.500.000\n\nVotación: 32 a favor, 4 en contra, 2 abstención\nSE APRUEBA EL PROYECTO.\n\nFirma: María González, Presidenta",
			"2025-09-20", "extraordinaria", 38},
	}

	for _, a := range actas {
		_, err = pool.Exec(ctx, `
			INSERT INTO actas (title, content, meeting_date, type, attendees_count, created_by)
			VALUES ($1, $2, $3, $4, $5, 'a0000000-0000-0000-0000-000000000004')
		`, a.title, a.content, a.date, a.tipo, a.attendees)
		if err != nil {
			log.Printf("Warning inserting acta: %v", err)
		}
	}

	// ============================================
	// DOCUMENTOS
	// ============================================
	log.Println("Inserting documentos...")
	documentos := []struct {
		title, desc, url, category string
		isPublic                   bool
	}{
		{"Reglamento Interno de Copropiedad", "Reglamento vigente de normas de convivencia.", "/documentos/reglamento-interno-2024.pdf", "reglamento", true},
		{"Reglamento de Piscina", "Normas de uso de la piscina comunitaria.", "/documentos/reglamento-piscina.pdf", "reglamento", true},
		{"Protocolo de Emergencias", "Procedimientos ante emergencias.", "/documentos/protocolo-emergencias.pdf", "protocolo", true},
		{"Protocolo de Mudanzas", "Procedimiento y horarios para mudanzas.", "/documentos/protocolo-mudanzas.pdf", "protocolo", true},
		{"Formulario Solicitud de Obras", "Formulario para autorización de modificaciones.", "/documentos/formulario-obras.pdf", "formulario", true},
		{"Formulario Arriendo Sede Social", "Solicitud para arriendo de sede social.", "/documentos/formulario-arriendo-sede.pdf", "formulario", true},
		{"Plano General del Condominio", "Plano con ubicación de parcelas y áreas comunes.", "/documentos/plano-condominio.pdf", "otro", true},
		{"Presupuesto 2026", "Presupuesto anual aprobado en asamblea.", "/documentos/presupuesto-2026.pdf", "otro", false},
	}

	for _, d := range documentos {
		_, err = pool.Exec(ctx, `
			INSERT INTO documentos (title, description, file_url, category, is_public, created_by)
			VALUES ($1, $2, $3, $4, $5, 'a0000000-0000-0000-0000-000000000002')
		`, d.title, d.desc, d.url, d.category, d.isPublic)
		if err != nil {
			log.Printf("Warning inserting documento: %v", err)
		}
	}

	// ============================================
	// EMERGENCIAS
	// ============================================
	log.Println("Inserting emergencias...")
	_, err = pool.Exec(ctx, `
		INSERT INTO emergencias (id, title, content, priority, status, expires_at, created_by, created_at) VALUES
		('e0000000-0000-0000-0000-000000000001',
		 'Corte de suministro eléctrico',
		 'Se informa que debido a trabajos de la empresa eléctrica, habrá corte de suministro el día de mañana entre las 09:00 y 13:00 hrs. Se recomienda tomar las precauciones necesarias.',
		 'medium', 'active', NOW() + INTERVAL '2 days',
		 'a0000000-0000-0000-0000-000000000002', NOW() - INTERVAL '1 hour'),
		('e0000000-0000-0000-0000-000000000002',
		 'Alerta de seguridad - Vehículo sospechoso',
		 'Se ha detectado un vehículo sospechoso rondando la comunidad. Favor reportar cualquier actividad inusual al guardia o directiva. Vehículo: Sedan gris, patente no identificada.',
		 'high', 'active', NOW() + INTERVAL '1 day',
		 'a0000000-0000-0000-0000-000000000002', NOW() - INTERVAL '3 hours'),
		('e0000000-0000-0000-0000-000000000003',
		 'Fuga de agua en calle principal',
		 'Se detectó una fuga de agua en la calle principal cerca de la entrada. El equipo de mantención está trabajando en la reparación. Se solicita precaución al transitar por la zona.',
		 'medium', 'resolved', NULL,
		 'a0000000-0000-0000-0000-000000000003', NOW() - INTERVAL '5 days')
	`)
	if err != nil {
		log.Printf("Warning inserting emergencias: %v", err)
	}

	// ============================================
	// VOTACIONES
	// ============================================
	log.Println("Inserting votaciones...")
	// Votación activa
	_, err = pool.Exec(ctx, `
		INSERT INTO votaciones (id, title, description, status, start_date, end_date, requires_quorum, quorum_percentage, allow_abstention, created_by) VALUES
		('b0000000-0000-0000-0000-000000000001',
		 'Proyecto de iluminación LED',
		 'Votación para aprobar el cambio de luminarias a tecnología LED en todas las áreas comunes. Inversión estimada: $3.500.000. Ahorro proyectado: 40% en consumo eléctrico.',
		 'active', NOW() - INTERVAL '2 days', NOW() + INTERVAL '5 days',
		 true, 50, true, 'a0000000-0000-0000-0000-000000000002'),
		('b0000000-0000-0000-0000-000000000002',
		 'Horario de piscina temporada verano',
		 'Seleccione el horario preferido para la piscina durante la temporada de verano 2026.',
		 'active', NOW() - INTERVAL '1 day', NOW() + INTERVAL '7 days',
		 false, 0, true, 'a0000000-0000-0000-0000-000000000002'),
		('b0000000-0000-0000-0000-000000000003',
		 'Contratación servicio de jardinería',
		 'Se sometió a votación la renovación del contrato con la empresa de jardinería actual o cambio a nuevo proveedor.',
		 'closed', NOW() - INTERVAL '30 days', NOW() - INTERVAL '20 days',
		 true, 50, true, 'a0000000-0000-0000-0000-000000000002')
	`)
	if err != nil {
		log.Printf("Warning inserting votaciones: %v", err)
	}

	// Opciones de votación
	_, err = pool.Exec(ctx, `
		INSERT INTO votacion_opciones (id, votacion_id, label, description, order_index) VALUES
		-- Opciones para iluminación LED
		('c0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'Aprobar proyecto', 'Aprobar la inversión y ejecución del proyecto de iluminación LED', 1),
		('c0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', 'Rechazar proyecto', 'Mantener el sistema de iluminación actual', 2),
		('c0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001', 'Postergar decisión', 'Solicitar más información antes de decidir', 3),
		-- Opciones para horario piscina
		('c0000000-0000-0000-0000-000000000004', 'b0000000-0000-0000-0000-000000000002', '09:00 - 20:00', 'Horario extendido de mañana a noche', 1),
		('c0000000-0000-0000-0000-000000000005', 'b0000000-0000-0000-0000-000000000002', '10:00 - 19:00', 'Horario estándar', 2),
		('c0000000-0000-0000-0000-000000000006', 'b0000000-0000-0000-0000-000000000002', '11:00 - 21:00', 'Horario tarde-noche', 3),
		-- Opciones para jardinería (cerrada)
		('c0000000-0000-0000-0000-000000000007', 'b0000000-0000-0000-0000-000000000003', 'Renovar contrato actual', 'Mantener empresa JardinPro por 1 año más', 1),
		('c0000000-0000-0000-0000-000000000008', 'b0000000-0000-0000-0000-000000000003', 'Cambiar proveedor', 'Contratar nueva empresa GreenService', 2)
	`)
	if err != nil {
		log.Printf("Warning inserting votacion_opciones: %v", err)
	}

	// Votos (para la votación cerrada)
	_, err = pool.Exec(ctx, `
		INSERT INTO votos (votacion_id, user_id, opcion_id, is_abstention, voted_at) VALUES
		('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000007', false, NOW() - INTERVAL '25 days'),
		('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000007', false, NOW() - INTERVAL '24 days'),
		('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000008', false, NOW() - INTERVAL '23 days'),
		('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000007', false, NOW() - INTERVAL '22 days'),
		('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000009', NULL, true, NOW() - INTERVAL '21 days')
	`)
	if err != nil {
		log.Printf("Warning inserting votos: %v", err)
	}

	// ============================================
	// GALERIAS
	// ============================================
	log.Println("Inserting galerias...")
	_, err = pool.Exec(ctx, `
		INSERT INTO galerias (id, title, description, event_date, is_public, cover_image_url, created_by) VALUES
		('d0000000-0000-0000-0000-000000000001',
		 'Celebración Fiestas Patrias 2025',
		 'Fotos de la celebración de fiestas patrias con asado comunitario, juegos y cueca.',
		 '2025-09-18', true, 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=800',
		 'a0000000-0000-0000-0000-000000000004'),
		('d0000000-0000-0000-0000-000000000002',
		 'Navidad en la Comunidad 2025',
		 'Celebración navideña con visita de Viejito Pascuero y regalos para los niños.',
		 '2025-12-20', true, 'https://images.unsplash.com/photo-1482517967863-00e15c9b44be?w=800',
		 'a0000000-0000-0000-0000-000000000004'),
		('d0000000-0000-0000-0000-000000000003',
		 'Mejoras áreas verdes 2025',
		 'Registro fotográfico de los trabajos de mejoramiento de áreas verdes.',
		 '2025-11-15', true, 'https://images.unsplash.com/photo-1558904541-efa843a96f01?w=800',
		 'a0000000-0000-0000-0000-000000000002')
	`)
	if err != nil {
		log.Printf("Warning inserting galerias: %v", err)
	}

	// Items de galería
	_, err = pool.Exec(ctx, `
		INSERT INTO galeria_items (galeria_id, file_url, thumbnail_url, file_type, caption, order_index) VALUES
		('d0000000-0000-0000-0000-000000000001', 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=1200', 'https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=400', 'image', 'Asado comunitario', 1),
		('d0000000-0000-0000-0000-000000000001', 'https://images.unsplash.com/photo-1528495612343-9ca9f4a4de28?w=1200', 'https://images.unsplash.com/photo-1528495612343-9ca9f4a4de28?w=400', 'image', 'Presentación de cueca', 2),
		('d0000000-0000-0000-0000-000000000001', 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=1200', 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=400', 'image', 'Mesa de comida típica', 3),
		('d0000000-0000-0000-0000-000000000002', 'https://images.unsplash.com/photo-1482517967863-00e15c9b44be?w=1200', 'https://images.unsplash.com/photo-1482517967863-00e15c9b44be?w=400', 'image', 'Árbol de navidad comunidad', 1),
		('d0000000-0000-0000-0000-000000000002', 'https://images.unsplash.com/photo-1512909006721-3d6018887383?w=1200', 'https://images.unsplash.com/photo-1512909006721-3d6018887383?w=400', 'image', 'Viejito Pascuero con niños', 2),
		('d0000000-0000-0000-0000-000000000002', 'https://images.unsplash.com/photo-1513297887119-d46091b24bfa?w=1200', 'https://images.unsplash.com/photo-1513297887119-d46091b24bfa?w=400', 'image', 'Entrega de regalos', 3),
		('d0000000-0000-0000-0000-000000000003', 'https://images.unsplash.com/photo-1558904541-efa843a96f01?w=1200', 'https://images.unsplash.com/photo-1558904541-efa843a96f01?w=400', 'image', 'Nuevas plantas', 1),
		('d0000000-0000-0000-0000-000000000003', 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=1200', 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=400', 'image', 'Sistema de riego', 2)
	`)
	if err != nil {
		log.Printf("Warning inserting galeria_items: %v", err)
	}

	// ============================================
	// MAPA - PUNTOS DE INTERÉS
	// ============================================
	log.Println("Inserting mapa puntos...")
	_, err = pool.Exec(ctx, `
		INSERT INTO mapa_puntos (id, name, description, lat, lng, icon, type, is_public) VALUES
		('f0000000-0000-0000-0000-000000000001', 'Entrada Principal', 'Acceso principal con control de guardia 24/7', -33.4500, -70.6500, 'gate', 'entrada', true),
		('f0000000-0000-0000-0000-000000000002', 'Sede Social', 'Sala de reuniones y eventos comunitarios', -33.4505, -70.6495, 'building', 'sede', true),
		('f0000000-0000-0000-0000-000000000003', 'Piscina Comunitaria', 'Piscina temperada disponible en verano', -33.4510, -70.6490, 'pool', 'recreacion', true),
		('f0000000-0000-0000-0000-000000000004', 'Plaza Central', 'Área de juegos infantiles y bancas', -33.4508, -70.6502, 'park', 'plaza', true),
		('f0000000-0000-0000-0000-000000000005', 'Cancha Multiuso', 'Cancha de fútbol y básquetbol', -33.4515, -70.6485, 'sports', 'recreacion', true),
		('f0000000-0000-0000-0000-000000000006', 'Estacionamiento Visitas', 'Estacionamiento para visitas', -33.4498, -70.6505, 'parking', 'estacionamiento', true),
		('f0000000-0000-0000-0000-000000000007', 'Bomba de Agua', 'Sistema de bombeo de agua potable', -33.4520, -70.6480, 'water', 'infraestructura', false),
		('f0000000-0000-0000-0000-000000000008', 'Sala de Basura', 'Punto de acopio y reciclaje', -33.4512, -70.6508, 'trash', 'servicios', true)
	`)
	if err != nil {
		log.Printf("Warning inserting mapa_puntos: %v", err)
	}

	// ============================================
	// MENSAJES DE CONTACTO
	// ============================================
	log.Println("Inserting mensajes contacto...")
	_, err = pool.Exec(ctx, `
		INSERT INTO mensajes_contacto (id, user_id, nombre, email, asunto, mensaje, status, created_at) VALUES
		('e0000000-0000-0000-0000-000000000001',
		 'a0000000-0000-0000-0000-000000000005',
		 'Juan Pérez', 'juan.perez@email.com',
		 'Consulta sobre estacionamiento de visitas',
		 'Estimada directiva, quisiera consultar sobre el procedimiento para solicitar el estacionamiento de visitas para un evento familiar que realizaré el próximo fin de semana. ¿Cuántos vehículos puedo registrar?',
		 'replied', NOW() - INTERVAL '5 days'),
		('e0000000-0000-0000-0000-000000000002',
		 'a0000000-0000-0000-0000-000000000006',
		 'Laura Rojas', 'laura.rojas@email.com',
		 'Ruidos molestos parcela vecina',
		 'Quiero reportar ruidos molestos provenientes de la parcela 12 durante las noches. Esto ha ocurrido los últimos 3 fines de semana después de las 23:00 hrs.',
		 'read', NOW() - INTERVAL '2 days'),
		('e0000000-0000-0000-0000-000000000003',
		 'a0000000-0000-0000-0000-000000000007',
		 'Pedro Castro', 'pedro.castro@email.com',
		 'Solicitud de poda de árbol',
		 'Solicito evaluación para podar el árbol ubicado frente a mi parcela, ya que las ramas están interfiriendo con el tendido eléctrico.',
		 'pending', NOW() - INTERVAL '1 day'),
		('e0000000-0000-0000-0000-000000000004',
		 NULL,
		 'María Fernández', 'maria.fernandez@email.com',
		 'Consulta compra de parcela',
		 'Buenas tardes, estoy interesada en conocer más sobre la comunidad. ¿Hay parcelas disponibles actualmente? Me gustaría agendar una visita.',
		 'pending', NOW() - INTERVAL '6 hours')
	`)
	if err != nil {
		log.Printf("Warning inserting mensajes_contacto: %v", err)
	}

	// ============================================
	// NOTIFICACIONES
	// ============================================
	log.Println("Inserting notificaciones...")
	// Notificaciones para Juan Pérez (usuario de prueba principal)
	_, err = pool.Exec(ctx, `
		INSERT INTO notificaciones (user_id, title, body, type, reference_id, is_read, created_at) VALUES
		('a0000000-0000-0000-0000-000000000005', 'Nuevo comunicado publicado', 'Se ha publicado: Normas de convivencia - Ruidos molestos', 'comunicado', NULL, false, NOW() - INTERVAL '1 day'),
		('a0000000-0000-0000-0000-000000000005', 'Recordatorio de pago', 'El plazo para el pago de gastos comunes vence en 3 días', 'pago', NULL, false, NOW() - INTERVAL '2 days'),
		('a0000000-0000-0000-0000-000000000005', 'Nueva votación disponible', 'Participa en: Proyecto de iluminación LED', 'votacion', 'b0000000-0000-0000-0000-000000000001', false, NOW() - INTERVAL '2 days'),
		('a0000000-0000-0000-0000-000000000005', 'Evento próximo', 'Recuerda: Celebración Día del Niño en 15 días', 'evento', NULL, true, NOW() - INTERVAL '3 days'),
		('a0000000-0000-0000-0000-000000000005', 'Respuesta a tu consulta', 'La directiva ha respondido tu mensaje sobre estacionamiento', 'contacto', 'e0000000-0000-0000-0000-000000000001', true, NOW() - INTERVAL '4 days'),
		('a0000000-0000-0000-0000-000000000005', 'Emergencia activa', 'Alerta de seguridad - Vehículo sospechoso', 'emergencia', 'e0000000-0000-0000-0000-000000000002', false, NOW() - INTERVAL '3 hours'),
		-- Notificaciones para otros usuarios
		('a0000000-0000-0000-0000-000000000006', 'Nuevo comunicado publicado', 'Se ha publicado: Normas de convivencia - Ruidos molestos', 'comunicado', NULL, true, NOW() - INTERVAL '1 day'),
		('a0000000-0000-0000-0000-000000000006', 'Nueva votación disponible', 'Participa en: Proyecto de iluminación LED', 'votacion', 'b0000000-0000-0000-0000-000000000001', false, NOW() - INTERVAL '2 days'),
		('a0000000-0000-0000-0000-000000000007', 'Nuevo comunicado publicado', 'Se ha publicado: Normas de convivencia - Ruidos molestos', 'comunicado', NULL, false, NOW() - INTERVAL '1 day'),
		('a0000000-0000-0000-0000-000000000008', 'Emergencia activa', 'Alerta de seguridad - Vehículo sospechoso', 'emergencia', 'e0000000-0000-0000-0000-000000000002', false, NOW() - INTERVAL '3 hours')
	`)
	if err != nil {
		log.Printf("Warning inserting notificaciones: %v", err)
	}

	log.Println("Seed data completed successfully!")

	// Show summary
	fmt.Println("\n=== RESUMEN DE DATOS INSERTADOS ===")
	tables := []string{
		"parcelas", "users", "periodos_gasto", "gastos_comunes", "pagos",
		"comunicados", "eventos", "movimientos_tesoreria",
		"actas", "documentos", "emergencias", "votaciones",
		"votacion_opciones", "votos", "galerias", "galeria_items",
		"mapa_puntos", "mensajes_contacto", "notificaciones",
	}
	for _, table := range tables {
		var count int
		err := pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("- %s: (tabla no existe)\n", table)
		} else {
			fmt.Printf("- %s: %d\n", table, count)
		}
	}
}
