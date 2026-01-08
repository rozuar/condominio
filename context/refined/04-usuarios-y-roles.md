# Sistema de Usuarios y Roles

## Métodos de Autenticación

### 1. Google (Gmail)
- Login con cuenta de Google
- OAuth 2.0
- Rápido y seguro

### 2. Email y Contraseña
- Registro tradicional
- Verificación por correo
- Recuperación de contraseña

## Roles del Sistema

### 1. Visitante (Público)
**Descripción:** Usuario no autenticado que navega el sitio.

**Permisos:**
- Ver página de inicio (parcial)
- Ver comunicados públicos
- Ver calendario (eventos públicos)
- Ver mapa de la comunidad
- Ver galería (contenido público)

**Restricciones:**
- No puede acceder a información privada
- No puede realizar pagos
- No puede votar
- No puede ver documentos internos

### 2. Vecino Autenticado
**Descripción:** Residente verificado de la comunidad.

**Permisos:**
- Todo lo del visitante
- Ver comunicados privados
- Pagar gastos comunes
- Ver estado de cuenta
- Ver tesorería y finanzas
- Ver actas y acuerdos
- Participar en votaciones
- Ver calendario completo
- Ver galería completa
- Ver documentos internos
- Contactar a la directiva

**Restricciones:**
- No puede crear contenido
- No puede administrar otros usuarios
- No puede editar información

### 3. Directiva / Administrador
**Descripción:** Miembro de la directiva con capacidad de gestión.

**Permisos:**
- Todo lo del vecino
- Crear y editar comunicados
- Crear avisos de emergencia
- Administrar gastos comunes
- Publicar estados financieros
- Crear y publicar actas
- Crear y administrar votaciones
- Gestionar calendario
- Administrar galería
- Gestionar documentos
- Ver mensajes de contacto
- Administrar usuarios

## Flujo de Registro

```
┌─────────────────┐
│  Nuevo Usuario  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Selecciona     │
│  método auth    │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐ ┌──────────┐
│Google │ │Email/Pass│
└───┬───┘ └────┬─────┘
    │          │
    └────┬─────┘
         │
         ▼
┌─────────────────┐
│  Verificación   │
│  de identidad   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Asignación     │
│  de parcela     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Aprobación     │
│  por directiva  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Usuario activo │
│  Rol: Vecino    │
└─────────────────┘
```

## Seguridad

- Control de acceso basado en roles (RBAC)
- Sesiones seguras con expiración
- Protección contra ataques CSRF
- Rate limiting en endpoints críticos
- Logs de acceso y auditoría
