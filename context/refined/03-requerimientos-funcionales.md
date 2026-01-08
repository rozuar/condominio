# Requerimientos Funcionales

> El contenido debe adaptarse automáticamente según el rol del usuario.

## Secciones del Sitio (12 módulos)

### 1. Inicio
- Mensajes relevantes de la comunidad
- Accesos rápidos a:
  - Comunicados
  - Emergencias
  - Pagar gastos comunes
  - Próximas fechas importantes

### 2. Comunicados
- Comunicados públicos (sin login)
- Comunicados privados (solo vecinos autenticados)
- Clasificación por tipo:
  - Informativo
  - Seguridad
  - Tesorería
  - Asamblea
- Historial ordenado por fecha

### 3. Avisos de Emergencia
- Sección prioritaria y visible
- Avisos críticos (seguridad, cortes, incidentes)
- Posibilidad de marcar avisos como urgentes
- Envío automático de notificación por correo a los vecinos

### 4. Gastos Comunes
- Pago online de gastos comunes
- Visualización del estado de pago del vecino
- Historial de pagos

### 5. Tesorería y Transparencia
- Publicación mensual de:
  - Movimientos de tesorería
  - Estado financiero actual
- Tablas claras y fáciles de leer
- Acceso solo con sesión iniciada

### 6. Actas y Acuerdos
- Registro histórico de:
  - Actas de reuniones
  - Acuerdos comunitarios
- Acceso solo para vecinos autenticados
- Lectura en línea, sin opción de descarga

### 7. Votaciones
- Sistema de votación online:
  - Votaciones simples
  - Escrutinios formales
- Un voto por usuario
- Resultados visibles según permisos
- Registro histórico de votaciones

### 8. Calendario Comunitario
- Calendario de:
  - Reuniones
  - Asambleas
  - Trabajos comunitarios
  - Actividades sociales
- Vista mensual y por evento
- Eventos visibles según rol

### 9. Mapa de la Comunidad
- Mapa interactivo con:
  - Parcelas
  - Áreas comunes
  - Accesos
  - Canal y referencias naturales

### 10. Galería Multimedia
- Fotos y videos de actividades comunitarias
- Organización por eventos o fechas
- Acceso privado o público según configuración

### 11. Documentos Internos
- Reglamentos, protocolos, información relevante
- Acceso solo con sesión iniciada
- Visualización sin descarga

### 12. Contacto con la Directiva
- Formulario privado de contacto
- Mensajes dirigidos a la directiva
- Protección anti-spam

## Matriz de Acceso por Sección

| Sección | Visitante | Vecino | Directiva |
|---------|-----------|--------|-----------|
| Inicio | Parcial | Completo | Completo |
| Comunicados públicos | Si | Si | Si + Editar |
| Comunicados privados | No | Si | Si + Editar |
| Emergencias | Parcial | Completo | Completo + Crear |
| Gastos comunes | No | Si | Si + Admin |
| Tesorería | No | Si | Si + Editar |
| Actas | No | Si | Si + Crear |
| Votaciones | No | Si | Si + Crear |
| Calendario | Parcial | Completo | Completo + Editar |
| Mapa | Si | Si | Si |
| Galería | Parcial | Completo | Completo + Editar |
| Documentos | No | Si | Si + Editar |
| Contacto | No | Si | Si + Ver mensajes |

## Sistema de Notificaciones

Envío automático de notificaciones por:

| Evento | Email | Push (Android) |
|--------|-------|----------------|
| Nuevo comunicado | ✓ | ✓ |
| Aviso de emergencia | ✓ (prioritario) | ✓ (inmediato) |
| Nueva votación | ✓ | ✓ |
| Publicación de acta | ✓ | - |
| Recordatorio de pago | ✓ | ✓ |
| Confirmación de pago | ✓ | ✓ |
