# Documentacion del Proyecto Condominio Vina Pelvin

## Arquitectura General

```
condominio/
├── source/
│   ├── api/          # Backend Go (API REST)
│   ├── web/          # Frontend Next.js (para vecinos)
│   └── backoffice/   # Admin Panel Next.js (para directiva)
└── context/          # Documentacion de contexto
```

## URLs de Produccion (Railway)

| Servicio   | URL                                                      |
|------------|----------------------------------------------------------|
| API        | https://asistencia-condominio-api-production.up.railway.app |
| Web        | https://web-production-f8c51.up.railway.app              |
| Backoffice | (pendiente configurar Root Directory en Railway)         |

---

## 1. API Backend (Go)

### Ubicacion
`/source/api/`

### Tecnologias
- Go 1.21+
- Chi Router
- PostgreSQL
- JWT para autenticacion

### Puerto Local
`8090`

### Endpoints Principales

#### Autenticacion (sin prefijo /api/v1)
```
POST /auth/login          # Login con email/password
POST /auth/refresh        # Refrescar tokens
GET  /auth/me             # Obtener usuario actual (requiere token)
GET  /auth/google         # Iniciar OAuth con Google
GET  /auth/google/callback
GET  /health              # Health check
```

#### API v1 (prefijo /api/v1)
```
# Comunicados (publico lectura, directiva escritura)
GET    /api/v1/comunicados
GET    /api/v1/comunicados/latest
GET    /api/v1/comunicados/{id}
POST   /api/v1/comunicados          # directiva
PUT    /api/v1/comunicados/{id}     # directiva
DELETE /api/v1/comunicados/{id}     # directiva

# Eventos (publico lectura, directiva escritura)
GET    /api/v1/eventos
GET    /api/v1/eventos/upcoming
GET    /api/v1/eventos/{id}
POST   /api/v1/eventos              # directiva
PUT    /api/v1/eventos/{id}         # directiva
DELETE /api/v1/eventos/{id}         # directiva

# Emergencias (publico lectura, directiva escritura)
GET    /api/v1/emergencias
GET    /api/v1/emergencias/active
GET    /api/v1/emergencias/{id}
POST   /api/v1/emergencias          # directiva
PUT    /api/v1/emergencias/{id}     # directiva
POST   /api/v1/emergencias/{id}/resolve  # directiva
DELETE /api/v1/emergencias/{id}     # directiva

# Votaciones (vecino+ lectura, directiva admin)
GET    /api/v1/votaciones           # vecino+
GET    /api/v1/votaciones/active    # vecino+
GET    /api/v1/votaciones/{id}      # vecino+
GET    /api/v1/votaciones/{id}/resultados  # vecino+
POST   /api/v1/votaciones/{id}/votar       # vecino+
POST   /api/v1/votaciones           # directiva
PUT    /api/v1/votaciones/{id}      # directiva
POST   /api/v1/votaciones/{id}/publish     # directiva
POST   /api/v1/votaciones/{id}/close       # directiva
POST   /api/v1/votaciones/{id}/cancel      # directiva
DELETE /api/v1/votaciones/{id}      # directiva

# Gastos Comunes (vecino+ lectura, directiva admin)
GET    /api/v1/gastos/periodos      # vecino+
GET    /api/v1/gastos/periodos/actual
GET    /api/v1/gastos/periodos/{id}
GET    /api/v1/gastos/periodos/{id}/resumen
GET    /api/v1/gastos/periodos/{id}/gastos
GET    /api/v1/gastos/mi-cuenta     # vecino+
GET    /api/v1/gastos/{id}
POST   /api/v1/gastos/periodos      # directiva
PUT    /api/v1/gastos/periodos/{id} # directiva
POST   /api/v1/gastos/{id}/pago     # directiva

# Contacto (publico crear, directiva gestionar)
POST   /api/v1/contacto             # publico
GET    /api/v1/contacto/mis-mensajes  # vecino+
GET    /api/v1/contacto             # directiva
GET    /api/v1/contacto/{id}        # directiva
POST   /api/v1/contacto/{id}/read   # directiva
POST   /api/v1/contacto/{id}/reply  # directiva
POST   /api/v1/contacto/{id}/archive # directiva
DELETE /api/v1/contacto/{id}        # directiva

# Galerias (publico lectura, directiva admin)
GET    /api/v1/galerias
GET    /api/v1/galerias/{id}
POST   /api/v1/galerias             # directiva
PUT    /api/v1/galerias/{id}        # directiva
DELETE /api/v1/galerias/{id}        # directiva
POST   /api/v1/galerias/{id}/items  # directiva

# Tesoreria (vecino+ lectura, directiva crear)
GET    /api/v1/tesoreria            # vecino+
GET    /api/v1/tesoreria/resumen    # vecino+
POST   /api/v1/tesoreria            # directiva

# Actas (vecino+ lectura, directiva crear)
GET    /api/v1/actas                # vecino+
GET    /api/v1/actas/{id}           # vecino+
POST   /api/v1/actas                # directiva

# Documentos (vecino+ lectura, directiva crear)
GET    /api/v1/documentos           # vecino+
GET    /api/v1/documentos/{id}      # vecino+
POST   /api/v1/documentos           # directiva

# Mapa (publico lectura, directiva admin)
GET    /api/v1/mapa                 # datos completos
GET    /api/v1/mapa/areas
GET    /api/v1/mapa/areas/{id}
POST   /api/v1/mapa/areas           # directiva
PUT    /api/v1/mapa/areas/{id}      # directiva
DELETE /api/v1/mapa/areas/{id}      # directiva
GET    /api/v1/mapa/puntos
GET    /api/v1/mapa/puntos/{id}
POST   /api/v1/mapa/puntos          # directiva
PUT    /api/v1/mapa/puntos/{id}     # directiva
DELETE /api/v1/mapa/puntos/{id}     # directiva

# Notificaciones (vecino+ lectura, directiva crear)
GET    /api/v1/notificaciones       # vecino+
GET    /api/v1/notificaciones/stats # vecino+
POST   /api/v1/notificaciones/{id}/read    # vecino+
POST   /api/v1/notificaciones/read-all     # vecino+
DELETE /api/v1/notificaciones/{id}  # vecino+
POST   /api/v1/notificaciones       # directiva
POST   /api/v1/notificaciones/bulk  # directiva
POST   /api/v1/notificaciones/broadcast    # directiva
```

### Roles de Usuario
- `visitor`: Usuario no autenticado
- `vecino`: Residente de la comunidad
- `directiva`: Miembro del comite (admin)
- `admin`: Administrador del sistema

### Usuarios de Prueba (Seed)
| Email                      | Password | Rol       |
|----------------------------|----------|-----------|
| admin@vinapelvin.cl        | password | admin     |
| presidente@vinapelvin.cl   | password | directiva |
| tesorero@vinapelvin.cl     | password | directiva |
| secretaria@vinapelvin.cl   | password | directiva |
| vecino1@vinapelvin.cl      | password | vecino    |

### Variables de Entorno API
```env
DATABASE_URL=postgres://...
JWT_SECRET=your-secret-key
PORT=8090
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
FRONTEND_URL=https://web-production.up.railway.app
```

---

## 2. Web Frontend (Next.js)

### Ubicacion
`/source/web/`

### Tecnologias
- Next.js 14 (App Router)
- React 18
- TypeScript
- Tailwind CSS
- Lucide Icons
- date-fns
- Leaflet (mapas)

### Puerto Local
`3000`

### Estructura de Carpetas
```
src/
├── app/                    # Paginas (App Router)
│   ├── api/health/        # Health check endpoint
│   ├── auth/              # Login, registro, callback
│   ├── comunicados/       # Lista y detalle
│   ├── eventos/           # Calendario y detalle
│   ├── emergencias/       # Lista y detalle
│   ├── votaciones/        # Lista, detalle, votar
│   ├── gastos/            # Estado de cuenta
│   ├── contacto/          # Formulario contacto
│   ├── documentos/        # Lista documentos
│   ├── galeria/           # Galerias de fotos
│   ├── mapa/              # Mapa interactivo
│   ├── tesoreria/         # Movimientos
│   └── actas/             # Lista actas
├── components/
│   ├── layout/            # Header, Footer, Sidebar
│   └── ui/                # Componentes reutilizables
├── contexts/
│   └── AuthContext.tsx    # Contexto de autenticacion
├── lib/
│   ├── api.ts             # Cliente API
│   └── auth.ts            # Helpers de autenticacion
└── types/
    └── index.ts           # Tipos TypeScript
```

### Variables de Entorno Web
```env
NEXT_PUBLIC_API_URL=https://asistencia-condominio-api-production.up.railway.app
```

---

## 3. Backoffice (Next.js Admin Panel)

### Ubicacion
`/source/backoffice/`

### Tecnologias
- Next.js 14 (App Router)
- React 18
- TypeScript
- Tailwind CSS (tema oscuro admin)
- Lucide Icons
- date-fns

### Puerto Local
`3001`

### Estructura de Carpetas
```
src/
├── app/
│   ├── api/health/        # Health check
│   ├── login/             # Login (solo directiva/admin)
│   ├── comunicados/       # CRUD
│   ├── eventos/           # CRUD
│   ├── emergencias/       # CRUD + resolver
│   ├── votaciones/        # CRUD + publicar/cerrar/cancelar
│   ├── gastos/            # Periodos y pagos
│   ├── contacto/          # Gestionar mensajes
│   ├── galerias/          # CRUD
│   ├── tesoreria/         # Movimientos
│   ├── actas/             # CRUD
│   ├── documentos/        # CRUD
│   ├── mapa/              # Areas y puntos
│   └── notificaciones/    # Enviar notificaciones
├── components/
│   ├── layout/
│   │   ├── Sidebar.tsx    # Navegacion lateral oscura
│   │   └── AuthGuard.tsx  # Proteccion de rutas
│   └── ui/
│       ├── Button.tsx
│       ├── Input.tsx
│       ├── Select.tsx
│       ├── Textarea.tsx
│       ├── Modal.tsx
│       ├── ConfirmDialog.tsx
│       └── Badge.tsx
├── contexts/
│   └── AuthContext.tsx    # Solo permite directiva/admin
├── lib/
│   └── api.ts             # Cliente API completo
└── types/
    └── index.ts           # Tipos compartidos
```

### Estilo Visual
- Sidebar: `bg-slate-900` (oscuro)
- Sidebar hover: `bg-slate-800`
- Content: `bg-gray-50`
- Cards: `bg-white`
- Accent: `blue-600`

### Variables de Entorno Backoffice
```env
NEXT_PUBLIC_API_URL=https://asistencia-condominio-api-production.up.railway.app
NEXT_PUBLIC_SHOW_DEV_PROFILES=false  # true en desarrollo
```

---

## 4. Despliegue en Railway

### Configuracion por Servicio

#### API
- Root Directory: `source/api`
- Builder: Nixpacks (detecta Go automaticamente)
- Variables:
  - DATABASE_URL
  - JWT_SECRET
  - PORT (automatico)

#### Web
- Root Directory: `source/web`
- Builder: Nixpacks
- Variables:
  - NEXT_PUBLIC_API_URL=https://asistencia-condominio-api-production.up.railway.app

#### Backoffice
- Root Directory: `source/backoffice`
- Builder: Nixpacks
- Variables:
  - NEXT_PUBLIC_API_URL=https://asistencia-condominio-api-production.up.railway.app
  - NEXT_PUBLIC_SHOW_DEV_PROFILES=false

### Archivos de Configuracion
- `railway.toml` - Configuracion de deploy
- `nixpacks.toml` - Configuracion del builder

---

## 5. Tipos TypeScript Principales

```typescript
// Usuario
interface User {
  id: string
  email: string
  name: string
  role: 'visitor' | 'vecino' | 'directiva' | 'admin'
  parcela_id?: number
  email_verified: boolean
  created_at: string
  updated_at: string
}

// Autenticacion
interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

// Comunicado
interface Comunicado {
  id: string
  title: string
  content: string
  type: 'general' | 'urgente' | 'mantenimiento' | 'seguridad' | 'evento'
  priority: 'low' | 'medium' | 'high'
  author_id: string
  author_name?: string
  created_at: string
  updated_at: string
}

// Evento
interface Evento {
  id: string
  title: string
  description: string
  location: string
  start_date: string
  end_date: string
  type: 'reunion' | 'social' | 'mantenimiento' | 'deportivo' | 'otro'
  is_mandatory: boolean
  created_at: string
}

// Emergencia
interface Emergencia {
  id: string
  title: string
  description: string
  priority: 'low' | 'medium' | 'high' | 'critical'
  status: 'active' | 'resolved' | 'cancelled'
  reported_by: string
  resolved_at?: string
  resolved_by?: string
  created_at: string
}

// Votacion
interface Votacion {
  id: string
  title: string
  description: string
  options: VotacionOpcion[]
  start_date: string
  end_date: string
  status: 'draft' | 'active' | 'closed' | 'cancelled'
  requires_quorum: boolean
  quorum_percentage: number
  created_at: string
}

// Gasto Comun
interface GastoComun {
  id: string
  periodo_id: string
  user_id: string
  parcela_id: number
  monto_base: number
  monto_extra: number
  monto_total: number
  status: 'pending' | 'paid' | 'overdue' | 'partial'
  fecha_pago?: string
  metodo_pago?: string
}
```

---

## 6. Comandos Utiles

### Desarrollo Local
```bash
# API
cd source/api
go run cmd/api/main.go

# Web
cd source/web
npm run dev

# Backoffice
cd source/backoffice
npm run dev
```

### Build
```bash
# Web/Backoffice
npm run build
npm start
```

### Seed Database
```bash
cd source/api
go run cmd/seed/main.go
```

---

## 7. Pendientes para App Movil (Kotlin)

### Endpoints a Consumir
La app movil debe consumir los mismos endpoints que la web:
- Autenticacion: `/auth/login`, `/auth/refresh`, `/auth/me`
- Comunicados: `/api/v1/comunicados`
- Eventos: `/api/v1/eventos`
- Emergencias: `/api/v1/emergencias`
- Votaciones: `/api/v1/votaciones` + votar
- Gastos: `/api/v1/gastos/mi-cuenta`
- Notificaciones: `/api/v1/notificaciones`

### Headers Requeridos
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

### Flujo de Autenticacion
1. Login con email/password o Google OAuth
2. Guardar access_token y refresh_token (SecureStorage)
3. Usar access_token en cada request
4. Cuando expire (401), usar refresh_token para obtener nuevos tokens
5. Si refresh falla, redirigir a login

### Push Notifications
- Implementar Firebase Cloud Messaging
- Registrar device token con el backend
- Endpoint pendiente: `POST /api/v1/devices/register`

---

## 8. Notas Importantes

1. **NEXT_PUBLIC_* variables**: Se "bake" durante el build. Deben estar configuradas en Railway ANTES del build.

2. **Root Directory en Railway**: Cada servicio debe tener configurado su directorio raiz (source/api, source/web, source/backoffice).

3. **Health Checks**:
   - API: `/health`
   - Web/Backoffice: `/api/health`

4. **CORS**: El API tiene CORS configurado para permitir requests desde cualquier origen.

5. **Roles**:
   - `directiva` incluye presidente, tesorero, secretaria
   - `admin` es superusuario del sistema
   - Backoffice solo permite `directiva` y `admin`
