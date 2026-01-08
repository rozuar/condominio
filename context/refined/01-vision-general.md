# Visión General del Proyecto

## Nombre del Proyecto
Sitio Web Comunidad de Viña Pelvin

## Contexto
Comunidad privada compuesta por **73 parcelas**, ubicada en un valle entre cerros y un canal, con entorno natural, vida comunitaria activa y necesidad real de orden, seguridad, transparencia y participación.

El sitio será el **canal digital oficial y único** de la comunidad.

## Objetivo Principal
Centralizar en un solo lugar:
- Comunicación oficial
- Gestión comunitaria
- Finanzas y transparencia
- Participación vecinal
- Seguridad y emergencias

**Evitar WhatsApp como canal crítico. Este sitio manda.**

## Criterios Clave
- Seguridad y control de acceso
- Transparencia total
- Fácil administración para la directiva
- Escalable a futuro

## Tono General
- Cercano, claro y serio
- Comunidad organizada, no improvisada
- Infraestructura digital comunitaria, no un blog

## Estructura del Proyecto
```
├── context
│   ├── raw          # Información sin procesar
│   └── refined      # Información organizada
└── source
    ├── arquitectura # Diagramas y documentación técnica
    ├── backend      # Backend de negocio (Go)
    │   ├── cmd
    │   ├── internal
    │   └── pkg
    ├── backoffice   # Panel de administración (Next.js)
    │   └── src
    ├── web          # Frontend público (Next.js)
    │   └── src
    └── mobile       # App Android (Kotlin)
```
