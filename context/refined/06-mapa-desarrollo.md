# Mapa de Desarrollo

## Diagrama General de Fases

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FASE 0: FUNDACIÓN                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Setup DB   │  │  Auth Base  │  │  API Core   │  │  CI/CD      │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FASE 1: NÚCLEO PÚBLICO                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                         │
│  │   Inicio    │  │ Comunicados │  │  Calendario │                         │
│  │   (Home)    │  │  Públicos   │  │   Público   │                         │
│  └─────────────┘  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       FASE 2: ÁREA PRIVADA VECINOS                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Comunicados │  │  Tesorería  │  │   Actas y   │  │ Documentos  │        │
│  │  Privados   │  │ Finanzas    │  │  Acuerdos   │  │  Internos   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FASE 3: INTERACCIÓN                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Gastos     │  │ Votaciones  │  │  Contacto   │  │Notificaciones│       │
│  │  Comunes    │  │   Online    │  │  Directiva  │  │   Email     │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        FASE 4: COMPLEMENTARIOS                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                         │
│  │    Mapa     │  │   Galería   │  │ Emergencias │                         │
│  │ Interactivo │  │ Multimedia  │  │  Críticas   │                         │
│  └─────────────┘  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FASE 5: BACKOFFICE                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Gestión    │  │  Gestión    │  │  Reportes   │  │   Admin     │        │
│  │  Usuarios   │  │ Contenido   │  │  Finanzas   │  │  General    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FASE 6: MOBILE                                     │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │              App móvil con funcionalidades críticas              │       │
│  └─────────────────────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Detalle por Fase

### FASE 0: Fundación
> Base técnica sobre la que se construye todo

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Setup Base de Datos | Esquema inicial, migraciones, conexión | - |
| Sistema de Auth | Login Google + Email/Pass, JWT, sesiones | DB |
| API Core | Estructura REST, middlewares, validación | DB, Auth |
| CI/CD | Pipeline de deploy, environments | - |

**Entregables:**
- [ ] Base de datos configurada y desplegada
- [ ] Endpoints de autenticación funcionando
- [ ] Estructura base de API con documentación
- [ ] Deploy automático configurado

---

### FASE 1: Núcleo Público
> Lo que ve cualquier visitante

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Home / Inicio | Landing con accesos rápidos | API Core |
| Comunicados Públicos | CRUD + listado público | API Core |
| Calendario Público | Eventos visibles sin login | API Core |

**Entregables:**
- [ ] Página de inicio responsive
- [ ] Sistema de comunicados (crear, listar, ver)
- [ ] Calendario con vista mensual
- [ ] Diseño base implementado (colores, tipografía)

---

### FASE 2: Área Privada Vecinos
> Contenido exclusivo para vecinos autenticados

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Comunicados Privados | Filtrado por rol, clasificación | Auth, Fase 1 |
| Tesorería | Estados financieros, movimientos | Auth, DB |
| Actas y Acuerdos | Historial, visualización sin descarga | Auth, DB |
| Documentos Internos | Reglamentos, protocolos | Auth, Storage |

**Entregables:**
- [ ] Comunicados con visibilidad por rol
- [ ] Módulo de tesorería con tablas claras
- [ ] Visualizador de actas (sin descarga)
- [ ] Repositorio de documentos internos

---

### FASE 3: Interacción
> Funcionalidades que requieren acción del usuario

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Gastos Comunes | Pago online, historial, estado | Auth, Pasarela de pago |
| Votaciones | Sistema de votos, escrutinio | Auth, DB |
| Contacto Directiva | Formulario privado | Auth, Email |
| Notificaciones | Emails automáticos por eventos | Email service |

**Entregables:**
- [ ] Integración con pasarela de pago
- [ ] Sistema de votación con un voto por usuario
- [ ] Formulario de contacto con anti-spam
- [ ] Sistema de notificaciones por email

---

### FASE 4: Complementarios
> Funcionalidades adicionales de valor

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Mapa Interactivo | Parcelas, áreas comunes, accesos | Frontend, Mapas API |
| Galería Multimedia | Fotos/videos por eventos | Storage, Auth |
| Emergencias | Avisos urgentes, notificación inmediata | Notificaciones |

**Entregables:**
- [ ] Mapa con 73 parcelas y áreas comunes
- [ ] Galería con organización por evento
- [ ] Sistema de alertas urgentes con email inmediato

---

### FASE 5: Backoffice
> Panel de administración para la directiva

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| Gestión Usuarios | CRUD usuarios, asignación roles | Auth, DB |
| Gestión Contenido | Admin de comunicados, actas, docs | Fases 1-4 |
| Reportes Finanzas | Dashboard financiero | Tesorería, Pagos |
| Admin General | Configuración del sistema | Todo |

**Entregables:**
- [ ] Panel de usuarios con roles
- [ ] CMS para todo el contenido
- [ ] Dashboard con métricas financieras
- [ ] Configuración general del sitio

---

### FASE 6: Mobile
> Aplicación Android nativa con Kotlin

| Componente | Descripción | Dependencias |
|------------|-------------|--------------|
| App Core | Autenticación, navegación, Jetpack Compose | API completa |
| Notificaciones Push | Firebase Cloud Messaging | Backend notificaciones |
| Funciones Críticas | Comunicados, emergencias, pagos | Fases 1-4 |

**Entregables:**
- [ ] App Android funcional (Play Store)
- [ ] Push notifications con FCM
- [ ] Funcionalidades principales portadas
- [ ] Material Design 3 implementado

---

## Grafo de Dependencias

```
                    ┌──────────┐
                    │    DB    │
                    └────┬─────┘
                         │
                         ▼
                    ┌──────────┐
                    │   Auth   │
                    └────┬─────┘
                         │
                         ▼
                    ┌──────────┐
                    │ API Core │
                    └────┬─────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌─────────┐    ┌──────────┐    ┌──────────┐
    │  Home   │    │Comunicados│   │Calendario│
    └────┬────┘    └─────┬────┘    └────┬─────┘
         │               │               │
         └───────────────┼───────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌─────────┐    ┌──────────┐    ┌──────────┐
    │Tesorería│    │  Actas   │    │Documentos│
    └────┬────┘    └──────────┘    └──────────┘
         │
         ▼
    ┌──────────┐
    │  Pagos   │◄────────────────┐
    └────┬─────┘                 │
         │                       │
         ▼                       │
    ┌──────────┐           ┌──────────┐
    │Votaciones│           │  Email   │
    └──────────┘           └────┬─────┘
                                │
                                ▼
                          ┌──────────┐
                          │Emergencias│
                          └──────────┘
```

---

## Stack Tecnológico Definido

| Capa | Tecnología | Justificación |
|------|------------|---------------|
| **Backend** | Go | Rendimiento, tipado, bajo consumo de recursos |
| **Base de Datos** | PostgreSQL | Relacional, robusto |
| **Frontend Web** | Next.js + TypeScript | SSR, SEO, ecosistema React |
| **Backoffice** | Next.js + TypeScript | Código compartido con web |
| **Mobile** | Android nativo (Kotlin) | Rendimiento nativo, Material Design 3 |
| **Email** | SendGrid / Resend | Deliverability |
| **Pagos** | Transbank / MercadoPago | Chile-friendly |
| **Storage** | S3 / GCS | Escalable |
| **Hosting** | Railway → GCP | Portabilidad |
| **Mapas** | Mapbox / Google Maps | Interactividad |

---

## Métricas de Avance

| Fase | Módulos | Peso |
|------|---------|------|
| Fase 0 | 4 | 15% |
| Fase 1 | 3 | 15% |
| Fase 2 | 4 | 20% |
| Fase 3 | 4 | 20% |
| Fase 4 | 3 | 10% |
| Fase 5 | 4 | 15% |
| Fase 6 | 3 | 5% |

**Total: 25 módulos**
