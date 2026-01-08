# Diseño y Experiencia de Usuario

## Estilo Visual

### Paleta de Colores
Inspirado en el entorno natural del valle:

| Color | Uso | Hex (sugerido) |
|-------|-----|----------------|
| Verde principal | Acciones, headers | #2D5016 |
| Verde claro | Fondos, highlights | #4A7C23 |
| Tierra | Acentos, secundario | #8B7355 |
| Agua/Azul | Links, información | #3B82A0 |
| Blanco | Fondos principales | #FFFFFF |
| Gris oscuro | Textos | #333333 |

### Tipografía
- Títulos: Sans-serif moderna, peso bold
- Cuerpo: Sans-serif legible, peso regular
- Tamaño base: 16px

### Iconografía
- Estilo: Línea simple, consistente
- Tamaño mínimo: 24px para accesibilidad

## Principios de Diseño

### 1. Sobrio y Profesional
- Sin elementos decorativos innecesarios
- Espacios en blanco generosos
- Jerarquía visual clara

### 2. Natural pero Moderno
- Referencias sutiles al entorno rural
- Interfaz limpia y contemporánea
- Balance entre calidez y profesionalismo

### 3. Comunidad Organizada
- Transmitir confianza y orden
- Información estructurada
- Acceso intuitivo a funciones

## Experiencia de Usuario

### Navegación
- Menú principal simple (máximo 6-7 items)
- Breadcrumbs en páginas internas
- Búsqueda global accesible
- Accesos rápidos desde inicio

### Responsive Design
- 100% adaptable a móvil y desktop
- Mobile-first approach
- Breakpoints:
  - Móvil: < 768px
  - Tablet: 768px - 1024px
  - Desktop: > 1024px

### Accesibilidad
- Contraste mínimo WCAG AA
- Navegación por teclado
- Textos alternativos en imágenes
- Formularios con labels claros

## Componentes Clave

### Header
```
┌─────────────────────────────────────────────────────────┐
│  [Logo]  Comunidad Viña Pelvin     [Nav] [User/Login]   │
└─────────────────────────────────────────────────────────┘
```

### Card de Comunicado
```
┌─────────────────────────────────────────────────────────┐
│  [Icono Tipo]  COMUNICADO                               │
│  ───────────────────────────────────────────────────    │
│  Título del comunicado                                  │
│  Fecha: 15 de enero, 2025                              │
│  Extracto del contenido...                             │
│                                            [Leer más]   │
└─────────────────────────────────────────────────────────┘
```

### Alerta de Emergencia
```
┌─────────────────────────────────────────────────────────┐
│  ⚠️  AVISO URGENTE                                      │
│  ───────────────────────────────────────────────────    │
│  Corte de agua programado para mañana 8:00 AM          │
│  Publicado: hace 2 horas                               │
└─────────────────────────────────────────────────────────┘
```

## Estados y Feedback

### Estados de Carga
- Skeleton screens para contenido
- Spinners sutiles para acciones
- Mensajes de progreso en operaciones largas

### Mensajes de Estado
- Éxito: Verde con icono check
- Error: Rojo con icono X
- Información: Azul con icono info
- Advertencia: Amarillo con icono alerta

### Formularios
- Validación en tiempo real
- Mensajes de error claros
- Indicadores de campos requeridos
- Confirmación antes de acciones destructivas
