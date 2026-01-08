# Arquitectura Técnica

## Visión General
Arquitectura modular, orientada a eventos, con separación clara entre componentes.

## Stack Tecnológico

| Componente | Tecnología |
|------------|------------|
| Backend | **Go** |
| Frontend Web | **Next.js + TypeScript** |
| Backoffice | **Next.js + TypeScript** |
| Mobile | **Android nativo (Kotlin)** |
| Base de Datos | PostgreSQL |
| Infraestructura | Railway → GCP |

## Componentes Principales

### 1. Backend de Negocio (Go)
- API REST central del sistema
- Lógica de negocio
- Gestión de datos
- Autenticación JWT

### 2. Frontend Web (Next.js + TypeScript)
- Sitio público y privado para vecinos
- Interfaz responsive (desktop y móvil)
- SSR para SEO y rendimiento

### 3. Frontend Backoffice (Next.js + TypeScript)
- Panel de administración para la directiva
- Gestión de contenido y usuarios
- Dashboards y reportes

### 4. Mobile (Android nativo - Kotlin)
- App Android nativa
- Notificaciones push
- Funcionalidades críticas
- UI moderna con Material Design 3

### 5. Canales de Notificación
- Envío de correos automáticos
- Notificaciones por eventos

### 6. Capa de Datos
- PostgreSQL como base principal
- Migraciones versionadas

### 7. Infraestructura
- Portable: Railway → GCP
- Escalable según necesidades

## Diagrama de Componentes
```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Web App    │  │  Backoffice │  │   Mobile    │     │
│  │  Next.js    │  │  Next.js    │  │Android/Kotlin│    │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘     │
└─────────┼────────────────┼────────────────┼─────────────┘
          │                │                │
          └────────────────┼────────────────┘
                           │ REST API
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    BACKEND (Go)                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │                  API REST                        │   │
│  └─────────────────────────────────────────────────┘   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Auth      │  │   Business  │  │   Events    │     │
│  │   Service   │  │   Logic     │  │   Handler   │     │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘     │
└─────────┼────────────────┼────────────────┼─────────────┘
          │                │                │
          └────────────────┼────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   SERVICIOS                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ PostgreSQL  │  │    Email    │  │   Storage   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

## Notificaciones Automáticas

El sistema enviará notificaciones automáticas a través de múltiples canales:

### Email
- Nuevos comunicados (públicos y privados)
- Avisos de emergencia (prioridad alta)
- Nuevas votaciones abiertas
- Publicación de actas
- Recordatorios de pago de gastos comunes

### Push (Mobile Android)
- Alertas de emergencia en tiempo real
- Nuevos comunicados importantes
- Recordatorios de votaciones activas
- Confirmaciones de pago

## Servicios de Terceros

| Servicio | Proveedor | Uso |
|----------|-----------|-----|
| Email transaccional | SendGrid / Resend | Notificaciones automáticas |
| Push notifications | Firebase Cloud Messaging | Alertas mobile |
| Pagos | Transbank / MercadoPago | Gastos comunes |
| Mapas | Mapbox / Google Maps | Mapa interactivo |
| Storage | S3 / Google Cloud Storage | Archivos y multimedia |
