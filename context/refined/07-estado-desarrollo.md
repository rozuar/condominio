# Estado del Desarrollo

> Última actualización: 10 Enero 2026

## Resumen General

| Componente | Estado | Progreso |
|------------|--------|----------|
| API (Go) | **Completo** | **100%** |
| Frontend Web (Next.js) | En desarrollo | 90% |
| Backoffice | En desarrollo | 50% |
| Mobile (Android/Kotlin) | En desarrollo | 45% |
| **Base de Datos** | **Schema listo** | **100%** |

---

## Base de Datos (PostgreSQL)

### Migraciones

```
source/api/internal/database/migrations/
├── 001_initial_schema.up.sql       # Schema base
├── 001_initial_schema.down.sql
├── 002_modulos_adicionales.up.sql  # Módulos futuros
├── 002_modulos_adicionales.down.sql
├── 003_seed_parcelas.up.sql        # 73 parcelas + admin
├── 003_seed_parcelas.down.sql
├── 004_mapa_comunidad.up.sql       # Mapa interactivo
└── 004_mapa_comunidad.down.sql
```

### Tablas Definidas

#### Migración 001 - Schema Inicial
| Tabla | Descripción | Estado |
|-------|-------------|--------|
| `parcelas` | 73 parcelas de la comunidad | ✅ |
| `users` | Usuarios con roles | ✅ |
| `comunicados` | Comunicados públicos/privados | ✅ |
| `eventos` | Calendario de eventos | ✅ |
| `movimientos` | Tesorería (ingresos/egresos) | ✅ |
| `actas` | Actas de reuniones | ✅ |
| `documentos` | Documentos internos | ✅ |

#### Migración 002 - Módulos Adicionales
| Tabla | Descripción | Estado |
|-------|-------------|--------|
| `votaciones` | Sistema de votación | ✅ Schema |
| `votacion_opciones` | Opciones de votación | ✅ Schema |
| `votos` | Votos emitidos | ✅ Schema |
| `emergencias` | Avisos urgentes | ✅ Schema |
| `galerias` | Álbumes de fotos | ✅ Schema |
| `galeria_items` | Fotos/videos | ✅ Schema |
| `periodos_gasto` | Períodos de cobro | ✅ Schema |
| `gastos_comunes` | Gastos por parcela | ✅ Schema |
| `pagos` | Pagos realizados | ✅ Schema |
| `mensajes_contacto` | Contacto con directiva | ✅ Schema |
| `notificaciones` | Sistema de notificaciones | ✅ Schema |

#### Migración 004 - Mapa
| Tabla | Descripción | Estado |
|-------|-------------|--------|
| `mapa_areas` | Polígonos del mapa | ✅ Schema |
| `mapa_puntos` | Puntos de interés | ✅ Schema |

### Tipos Enumerados (ENUMs)
- `user_role`: visitor, vecino, directiva
- `comunicado_type`: informativo, seguridad, tesoreria, asamblea
- `evento_type`: reunion, asamblea, trabajo, social
- `movimiento_type`: ingreso, egreso
- `acta_type`: ordinaria, extraordinaria
- `documento_category`: reglamento, protocolo, formulario, otro
- `votacion_status`: draft, active, closed, cancelled
- `emergencia_priority`: low, medium, high, critical
- `emergencia_status`: active, resolved, expired
- `pago_status`: pending, paid, overdue, cancelled
- `contacto_status`: pending, read, replied, archived
- `area_type`: parcela, area_comun, acceso, canal, camino

---

## API (Go)

### Stack Implementado
- **Framework**: chi/v5
- **Base de datos**: PostgreSQL (pgx/v5)
- **Auth**: JWT (golang-jwt/v5)
- **Go version**: 1.22

### Estructura
```
source/api/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── database/
│   │   └── migrations/    # ✅ 4 migraciones
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── router/
│   └── services/
└── pkg/
    ├── jwt/
    ├── email/             # ✅ Sistema de emails
    └── oauth/             # ✅ Google OAuth
```

### Módulos Implementados

| Módulo | Handler | Service | Model | API | DB Schema |
|--------|---------|---------|-------|-----|-----------|
| Auth | ✅ | ✅ | ✅ | ✅ | ✅ |
| Comunicados | ✅ | ✅ | ✅ | ✅ | ✅ |
| Eventos | ✅ | ✅ | ✅ | ✅ | ✅ |
| Tesorería | ✅ | ✅ | ✅ | ✅ | ✅ |
| Actas | ✅ | ✅ | ✅ | ✅ | ✅ |
| Documentos | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Emergencias** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Votaciones** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Gastos Comunes** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Galeria** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Mapa** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Contacto** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Notificaciones** | ✅ | ✅ | ✅ | ✅ | ✅ |

### API Endpoints Actuales

```
GET  /health                    # Health check

POST /auth/login                # Login
POST /auth/refresh              # Refresh token
GET  /auth/me                   # Usuario actual (auth)
GET  /auth/google               # Iniciar login con Google OAuth
GET  /auth/google/callback      # Callback de Google OAuth

GET  /api/v1/comunicados        # Lista pública
GET  /api/v1/comunicados/latest # Últimos comunicados
GET  /api/v1/comunicados/{id}   # Detalle
POST /api/v1/comunicados        # Crear (directiva)
PUT  /api/v1/comunicados/{id}   # Editar (directiva)
DEL  /api/v1/comunicados/{id}   # Eliminar (directiva)

GET  /api/v1/eventos            # Lista
GET  /api/v1/eventos/upcoming   # Próximos
GET  /api/v1/eventos/{id}       # Detalle
POST /api/v1/eventos            # Crear (directiva)
PUT  /api/v1/eventos/{id}       # Editar (directiva)
DEL  /api/v1/eventos/{id}       # Eliminar (directiva)

GET  /api/v1/tesoreria          # Lista (vecino+)
GET  /api/v1/tesoreria/resumen  # Resumen (vecino+)
POST /api/v1/tesoreria          # Crear (directiva)

GET  /api/v1/actas              # Lista (vecino+)
GET  /api/v1/actas/{id}         # Detalle (vecino+)
POST /api/v1/actas              # Crear (directiva)

GET  /api/v1/documentos         # Lista (vecino+)
GET  /api/v1/documentos/{id}    # Detalle (vecino+)
POST /api/v1/documentos         # Crear (directiva)

GET  /api/v1/emergencias        # Lista con filtros
GET  /api/v1/emergencias/active # Solo activas (por prioridad)
GET  /api/v1/emergencias/{id}   # Detalle
POST /api/v1/emergencias        # Crear (directiva)
PUT  /api/v1/emergencias/{id}   # Editar (directiva)
POST /api/v1/emergencias/{id}/resolve  # Resolver (directiva)
DEL  /api/v1/emergencias/{id}   # Eliminar (directiva)

GET  /api/v1/votaciones         # Lista (vecino+)
GET  /api/v1/votaciones/active  # Activas (vecino+)
GET  /api/v1/votaciones/{id}    # Detalle (vecino+)
GET  /api/v1/votaciones/{id}/resultados  # Resultados (vecino+)
POST /api/v1/votaciones/{id}/votar       # Emitir voto (vecino+)
POST /api/v1/votaciones         # Crear (directiva)
PUT  /api/v1/votaciones/{id}    # Editar (directiva)
POST /api/v1/votaciones/{id}/publish  # Publicar (directiva)
POST /api/v1/votaciones/{id}/close    # Cerrar (directiva)
POST /api/v1/votaciones/{id}/cancel   # Cancelar (directiva)
DEL  /api/v1/votaciones/{id}    # Eliminar (directiva)

GET  /api/v1/gastos/periodos         # Lista periodos (vecino+)
GET  /api/v1/gastos/periodos/actual  # Periodo actual (vecino+)
GET  /api/v1/gastos/periodos/{id}    # Detalle periodo (vecino+)
GET  /api/v1/gastos/periodos/{id}/resumen  # Resumen (vecino+)
GET  /api/v1/gastos/periodos/{id}/gastos   # Gastos del periodo (vecino+)
GET  /api/v1/gastos/mi-cuenta        # Mi estado de cuenta (vecino+)
GET  /api/v1/gastos/{id}             # Detalle gasto (vecino+)
POST /api/v1/gastos/periodos         # Crear periodo (directiva)
PUT  /api/v1/gastos/periodos/{id}    # Editar periodo (directiva)
POST /api/v1/gastos/{id}/pago        # Registrar pago (directiva)
POST /api/v1/gastos/marcar-vencidos  # Marcar vencidos (directiva)

POST /api/v1/contacto            # Enviar mensaje (publico)
GET  /api/v1/contacto/mis-mensajes  # Mis mensajes (vecino+)
GET  /api/v1/contacto            # Lista mensajes (directiva)
GET  /api/v1/contacto/stats      # Estadisticas (directiva)
GET  /api/v1/contacto/{id}       # Detalle mensaje (directiva)
POST /api/v1/contacto/{id}/read  # Marcar leido (directiva)
POST /api/v1/contacto/{id}/reply # Responder (directiva)
POST /api/v1/contacto/{id}/archive  # Archivar (directiva)
DEL  /api/v1/contacto/{id}       # Eliminar (directiva)

GET  /api/v1/galerias            # Lista galerias (publico)
GET  /api/v1/galerias/{id}       # Detalle con items (publico)
POST /api/v1/galerias            # Crear galeria (directiva)
PUT  /api/v1/galerias/{id}       # Editar galeria (directiva)
DEL  /api/v1/galerias/{id}       # Eliminar galeria (directiva)
POST /api/v1/galerias/{id}/items # Agregar item (directiva)
POST /api/v1/galerias/{id}/reorder  # Reordenar items (directiva)
PUT  /api/v1/galerias/{id}/items/{itemId}  # Editar item (directiva)
DEL  /api/v1/galerias/{id}/items/{itemId}  # Eliminar item (directiva)

GET  /api/v1/mapa                # Todos los datos del mapa (publico)
GET  /api/v1/mapa/areas          # Lista areas (publico)
GET  /api/v1/mapa/areas/{id}     # Detalle area (publico)
POST /api/v1/mapa/areas          # Crear area (directiva)
PUT  /api/v1/mapa/areas/{id}     # Editar area (directiva)
DEL  /api/v1/mapa/areas/{id}     # Eliminar area (directiva)
GET  /api/v1/mapa/puntos         # Lista puntos (publico)
GET  /api/v1/mapa/puntos/{id}    # Detalle punto (publico)
POST /api/v1/mapa/puntos         # Crear punto (directiva)
PUT  /api/v1/mapa/puntos/{id}    # Editar punto (directiva)
DEL  /api/v1/mapa/puntos/{id}    # Eliminar punto (directiva)

GET  /api/v1/notificaciones       # Lista notificaciones (vecino+)
GET  /api/v1/notificaciones/stats # Estadisticas (vecino+)
GET  /api/v1/notificaciones/{id}  # Detalle (vecino+)
POST /api/v1/notificaciones/{id}/read  # Marcar leida (vecino+)
POST /api/v1/notificaciones/read-all   # Marcar todas leidas (vecino+)
DEL  /api/v1/notificaciones/{id}  # Eliminar una (vecino+)
DEL  /api/v1/notificaciones       # Eliminar todas (vecino+)
DEL  /api/v1/notificaciones/read  # Eliminar leidas (vecino+)
POST /api/v1/notificaciones       # Crear para usuario (directiva)
POST /api/v1/notificaciones/bulk  # Crear para multiples usuarios (directiva)
POST /api/v1/notificaciones/broadcast  # Broadcast por roles (directiva)
```

### Pendiente Backend
- [x] ~~Migraciones SQL~~
- [x] ~~Módulo de Emergencias (handler, service, model)~~
- [x] ~~Módulo de Votaciones (handler, service, model)~~
- [x] ~~Módulo de Gastos Comunes (handler, service, model)~~
- [x] ~~Módulo de Contacto (handler, service, model)~~
- [x] ~~Módulo de Galeria (handler, service, model)~~
- [x] ~~Módulo de Mapa (handler, service, model)~~
- [x] ~~Módulo de Notificaciones (handler, service, model)~~
- [x] ~~Sistema de envío de emails~~
- [x] ~~Login con Google OAuth~~

**API 100% Completa**

---

## Frontend Web (Next.js)

### Stack Implementado
- **Framework**: Next.js 14.2
- **UI**: React 18 + TypeScript
- **Estilos**: Tailwind CSS 3.4
- **Iconos**: Lucide React
- **Fechas**: date-fns

### Estructura
```
web/src/
├── app/
│   ├── actas/
│   ├── auth/
│   ├── calendario/
│   ├── comunicados/
│   ├── documentos/
│   ├── emergencias/        # ✅ NUEVO
│   │   ├── [id]/
│   │   └── page.tsx
│   ├── tesoreria/
│   ├── layout.tsx
│   └── page.tsx
├── components/
│   ├── auth/
│   ├── calendario/
│   ├── comunicados/
│   ├── emergencias/        # ✅ NUEVO
│   │   ├── EmergenciaCard.tsx
│   │   └── EmergenciaBanner.tsx
│   ├── layout/
│   └── ui/
├── contexts/
│   └── AuthContext.tsx
├── lib/
│   ├── api.ts              # ✅ +emergencias
│   └── auth.ts
└── types/
    └── index.ts            # ✅ Tipos completos
```

### Tipos TypeScript

| Módulo | Tipos | Labels | Colors |
|--------|-------|--------|--------|
| User/Auth | ✅ | - | - |
| Parcela | ✅ | - | - |
| Comunicado | ✅ | ✅ | ✅ |
| Evento | ✅ | ✅ | ✅ |
| Movimiento | ✅ | ✅ | ✅ |
| Acta | ✅ | ✅ | ✅ |
| Documento | ✅ | ✅ | ✅ |
| Votacion | ✅ | ✅ | ✅ |
| Emergencia | ✅ | ✅ | ✅ |
| Galeria | ✅ | - | - |
| GastoComun | ✅ | ✅ | ✅ |
| Contacto | ✅ | ✅ | - |
| Notificacion | ✅ | - | - |
| Mapa | ✅ | - | - |

### Páginas Implementadas

| Página | Ruta | UI | API |
|--------|------|-----|-----|
| Home | `/` | ✅ | ✅ |
| Comunicados | `/comunicados` | ✅ | ✅ |
| Calendario | `/calendario` | ✅ | ✅ |
| Tesorería | `/tesoreria` | ✅ | ✅ |
| Actas | `/actas` | ✅ | ✅ |
| Documentos | `/documentos` | ✅ | ✅ |
| Auth | `/auth` | ✅ | ✅ |
| **Emergencias** | `/emergencias` | ✅ | ✅ |
| **Votaciones** | `/votaciones` | ✅ | ✅ |
| **Gastos Comunes** | `/gastos` | ✅ | ✅ |
| **Galeria** | `/galeria` | ✅ | ✅ |
| **Mapa** | `/mapa` | ✅ | ✅ |
| **Contacto** | `/contacto` | ✅ | ✅ |

### Diseño
- Colores del brand implementados en Tailwind:
  - `primary`: #2D5016 (verde principal)
  - `primary-light`: #4A7C23
  - `tierra`: #8B7355
  - `agua`: #3B82A0

### Pendiente Frontend
- [x] ~~Tipos TypeScript para todos los módulos~~
- [x] ~~Página de Emergencias~~
- [x] ~~Página de Votaciones~~
- [x] ~~Página de Gastos Comunes~~ (pasarela de pago pendiente)
- [x] ~~Página de Contacto~~
- [x] ~~Página de Galeria~~ (con lightbox para ver fotos/videos)
- [x] ~~Página de Mapa~~ (interfaz con sidebar, pendiente integrar libreria de mapas)

---

## Backoffice

**Estado**: En desarrollo

Proyecto Next.js creado en `source/backoffice/` para panel de administración (directiva).

---

## Mobile (Android/Kotlin)

**Estado**: En desarrollo

Proyecto Android nativo en `source/mobile/` (Jetpack Compose + Hilt + Retrofit + FCM).

### Implementado ✅
- Auth (email/pass) con JWT (almacenamiento en DataStore)
- Navegación base (Login → Home)
- Módulos: Comunicados, Eventos, Emergencias, Votaciones, Gastos (mi cuenta), Tesorería, Actas, Documentos, Notificaciones, Contacto
- Push: Firebase Cloud Messaging (topics + canales de notificación)

### Pendiente 📌
- Módulos: Galería, Mapa
- Ajustes de UX según roles (modo visitante vs autenticado)
- Registro de token FCM en backend (si se requiere push por usuario)

---

## Infraestructura

### Archivos de Configuración
- `source/api/Dockerfile` ✅
- `source/api/docker-compose.yml` ✅
- `source/api/Makefile` ✅

### Pendiente
- [ ] CI/CD pipeline
- [ ] Configuración Railway/GCP
- [ ] Variables de entorno producción

---

## Progreso por Fase

| Fase | Descripción | Estado | Detalle |
|------|-------------|--------|---------|
| Fase 0 | Fundación | 🟢 85% | DB schema completo, Auth funcionando |
| Fase 1 | Núcleo Público | 🟡 70% | Home, Comunicados, Calendario listos |
| Fase 2 | Área Privada | 🟢 80% | Tesorería, Actas, Documentos, Gastos listos |
| Fase 3 | Interacción | 🟢 100% | **Votaciones y Contacto completos** |
| Fase 4 | Complementarios | 🟢 100% | **Emergencias, Galeria y Mapa completos** |
| Fase 5 | Backoffice | 🟡 50% | En desarrollo |
| Fase 6 | Mobile | ❌ 0% | No iniciado |

---

## Checklist General

### Completado ✅
- [x] Schema de base de datos completo (4 migraciones)
- [x] Tipos TypeScript para todos los módulos
- [x] Backend: Auth, Comunicados, Eventos, Tesorería, Actas, Documentos
- [x] Frontend: Home, Comunicados, Calendario, Tesorería, Actas, Documentos
- [x] Seed de 73 parcelas
- [x] Usuario admin por defecto
- [x] Corrección typo carpeta `arquitectura`
- [x] **Módulo de Emergencias completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints con filtros y resolución
  - Página con lista, detalle y filtros
  - Componentes: EmergenciaCard, EmergenciaBanner
  - Prioridades visuales (critical, high, medium, low)
- [x] **Módulo de Votaciones completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints: CRUD, publicar, cerrar, cancelar, votar, resultados
  - Página con lista de votaciones y filtros por estado
  - Página de detalle con formulario de votación
  - Componentes: VotacionCard
  - Sistema de quorum y abstención
  - Visualización de resultados con porcentajes
- [x] **Módulo de Gastos Comunes completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints: periodos CRUD, gastos, mi-cuenta, pagos
  - Página Mi Estado de Cuenta con gastos pendientes y pagados
  - Sistema de periodos mensuales
  - Control de pagos parciales y vencidos
  - Generación automática de gastos para todas las parcelas
- [x] **Módulo de Contacto completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints: enviar, mis-mensajes, list, read, reply, archive
  - Página de contacto con formulario público
  - Vista de mensajes enviados para usuarios autenticados
  - Sistema de estados: pendiente, leído, respondido, archivado
- [x] **Módulo de Galeria completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints: CRUD galerias, CRUD items, reordenar items
  - Página de galeria con grid de albumes
  - Página de detalle con lightbox para ver fotos/videos
  - Soporte para imagenes y videos
  - Navegación con teclado en lightbox (flechas y Escape)
- [x] **Módulo de Mapa completo (Backend + Frontend)**
  - Model, Service, Handler en Go
  - API endpoints: GET mapa completo, CRUD areas, CRUD puntos
  - Página de mapa con sidebar interactivo
  - Visualización de areas por tipo (parcelas, areas comunes, accesos, etc.)
  - Lista de puntos de interes
  - Panel de información para items seleccionados
  - Estadísticas de la comunidad (parcelas, areas, puntos)
- [x] **Módulo de Notificaciones completo (Backend)**
  - Model, Service, Handler en Go
  - 10 tipos de notificación (comunicado, emergencia, votacion, pago, evento, documento, acta, contacto, gasto_comun, sistema)
  - API endpoints: lista, stats, detalle, marcar leida, marcar todas leidas
  - API endpoints para eliminar (una, todas, solo leidas)
  - API admin: crear para usuario, crear bulk, broadcast por roles
- [x] **Sistema de Emails completo (Backend)**
  - Servicio SMTP con soporte TLS/STARTTLS
  - Configuración via variables de entorno (SMTP_HOST, SMTP_PORT, etc.)
  - Plantillas HTML para emails:
    - Confirmación de mensaje de contacto
    - Respuesta de directiva a mensaje de contacto
    - Notificaciones generales
    - Alertas de emergencia
    - Notificaciones de gastos comunes
    - Email de bienvenida
  - Integrado con módulo de Contacto (envía confirmación y respuesta)
  - Integrado con módulo de Notificaciones (broadcast con opción de email)
  - Branding consistente con colores de Viña Pelvin
- [x] **Login con Google OAuth completo (Backend)**
  - Servicio OAuth en `pkg/oauth/google.go`
  - Configuración via variables de entorno (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, etc.)
  - Flujo completo: /auth/google -> Google -> /auth/google/callback -> Frontend
  - Protección CSRF con state aleatorio en cookie
  - Login automático si usuario existe (por Google ID o email)
  - **Sin registro automático**: si el usuario no existe, el login por Google se rechaza
  - Vinculación de cuenta Google a cuenta existente por email
  - Redirección a frontend con tokens en URL (/auth/callback)

### Próximos Pasos
1. Integrar libreria de mapas (Leaflet) en frontend
2. Integrar pasarela de pago (Transbank/MercadoPago)
3. Implementar página /auth/callback en frontend para recibir tokens OAuth
4. Continuar backoffice (incluye gestión de usuarios)
5. Iniciar app mobile Android
