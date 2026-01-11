# Sistema de Usuarios y Roles

## Métodos de Autenticación

### 1. Google (Gmail)
- Login con cuenta de Google
- OAuth 2.0
- Rápido y seguro

### 2. Email y Contraseña
- **Sin auto-registro**: solo usuarios previamente creados por la directiva
- Login tradicional con email/contraseña
- Recuperación de contraseña (pendiente)

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

### 4. Familia / Invitado (Autenticado)
**Descripción:** Usuario autenticado perteneciente a una familia asociada al condominio, o invitado interno. Puede existir **sin parcela asignada** (p. ej. familiares que usan la app, cuidadores, etc.).

**Permisos (app móvil / web):**
- Ver comunicados y eventos (según visibilidad del contenido)
- Ver avisos/alertas de emergencias
- Contactar a la directiva
- Acceder a menús generales de información

**Restricciones:**
- Si **no tiene `parcela_id`**, **no tiene Gastos Comunes** (no es error; se muestra estado informativo)
- No administra tesorería/finanzas
- No administra usuarios ni contenido

## Condición: usuario sin parcela (Gastos Comunes)

- **Contexto**: la parcela representa una unidad que “genera gastos”. Es normal que existan cuentas internas **sin parcela**.
- **Comportamiento esperado**:
  - El módulo “Gastos” debe mostrar un estado informativo (no error): **“Ud no posee asociada una parcela que genere gastos”**.
  - El usuario mantiene acceso a módulos generales.

## Alta de usuarios (sin auto-registro)

Las cuentas **no se crean desde la web pública**. La directiva debe:

1. Crear la cuenta del vecino (asignando rol y, si aplica, parcela).
2. Entregar credenciales (o habilitar login por Google vinculando el email).
3. El vecino inicia sesión desde `/auth/login`.

**Nota:** El login por Google **no crea usuarios automáticamente**. Solo permite entrar a cuentas ya existentes (y puede vincular Google por email).

## Seguridad

- Control de acceso basado en roles (RBAC)
- Sesiones seguras con expiración
- Protección contra ataques CSRF
- Rate limiting en endpoints críticos
- Logs de acceso y auditoría

## Matriz resumida de acceso (móvil)

| Módulo | Público | Familia/Invitado | Vecino | Directiva | Admin |
|---|---:|---:|---:|---:|---:|
| Comunicados | ✅ | ✅ | ✅ | ✅ | ✅ |
| Eventos | ✅ | ✅ | ✅ | ✅ | ✅ |
| Emergencias | ✅ | ✅ | ✅ | ✅ | ✅ |
| Contacto Directiva | ❌ | ✅ | ✅ | ✅ | ✅ |
| Actas | ❌ | ✅ | ✅ | ✅ | ✅ |
| Documentos | ❌ | ✅ | ✅ | ✅ | ✅ |
| Notificaciones | ❌ | ✅ | ✅ | ✅ | ✅ |
| Votaciones | ❌ | ✅ (solo ver)* | ✅ | ✅ | ✅ |
| Tesorería | ❌ | ❌ | ✅ | ✅ | ✅ |
| Gastos | ❌ | ✅* | ✅ | ✅ | ✅ |

\* **Si no hay `parcela_id`**:
- **Gastos**: no hay gasto asociado y se muestra estado informativo.
- **Votaciones**: solo lectura (ver estado/resultados), **no puede votar**.
