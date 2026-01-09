export type ComunicadoType = 'informativo' | 'seguridad' | 'tesoreria' | 'asamblea'

export interface Comunicado {
  id: string
  title: string
  content: string
  type: ComunicadoType
  is_public: boolean
  author_id?: string
  author_name?: string
  published_at?: string
  created_at: string
  updated_at: string
}

export interface ComunicadoListResponse {
  comunicados: Comunicado[]
  total: number
  page: number
  per_page: number
}

export type EventoType = 'reunion' | 'asamblea' | 'trabajo' | 'social'

export interface Evento {
  id: string
  title: string
  description?: string
  event_date: string
  event_end_date?: string
  location?: string
  type: EventoType
  is_public: boolean
  created_by?: string
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface EventoListResponse {
  eventos: Evento[]
  total: number
  page: number
  per_page: number
}

export const COMUNICADO_TYPE_LABELS: Record<ComunicadoType, string> = {
  informativo: 'Informativo',
  seguridad: 'Seguridad',
  tesoreria: 'Tesoreria',
  asamblea: 'Asamblea',
}

export const EVENTO_TYPE_LABELS: Record<EventoType, string> = {
  reunion: 'Reunion',
  asamblea: 'Asamblea',
  trabajo: 'Trabajo Comunitario',
  social: 'Actividad Social',
}

export const COMUNICADO_TYPE_COLORS: Record<ComunicadoType, string> = {
  informativo: 'bg-blue-100 text-blue-800',
  seguridad: 'bg-red-100 text-red-800',
  tesoreria: 'bg-yellow-100 text-yellow-800',
  asamblea: 'bg-purple-100 text-purple-800',
}

export const EVENTO_TYPE_COLORS: Record<EventoType, string> = {
  reunion: 'bg-blue-100 text-blue-800',
  asamblea: 'bg-purple-100 text-purple-800',
  trabajo: 'bg-orange-100 text-orange-800',
  social: 'bg-green-100 text-green-800',
}

export type MovimientoType = 'ingreso' | 'egreso'

export interface Movimiento {
  id: string
  description: string
  amount: number
  type: MovimientoType
  category: string
  date: string
  created_by?: string
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface MovimientoListResponse {
  movimientos: Movimiento[]
  total: number
  page: number
  per_page: number
}

export interface TesoreriaResumen {
  total_ingresos: number
  total_egresos: number
  balance: number
  movimientos_count: number
}

export const MOVIMIENTO_TYPE_LABELS: Record<MovimientoType, string> = {
  ingreso: 'Ingreso',
  egreso: 'Egreso',
}

export const MOVIMIENTO_TYPE_COLORS: Record<MovimientoType, string> = {
  ingreso: 'bg-green-100 text-green-800',
  egreso: 'bg-red-100 text-red-800',
}

export type ActaType = 'ordinaria' | 'extraordinaria'

export interface Acta {
  id: string
  title: string
  content: string
  meeting_date: string
  type: ActaType
  attendees_count: number
  created_by?: string
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface ActaListResponse {
  actas: Acta[]
  total: number
  page: number
  per_page: number
}

export const ACTA_TYPE_LABELS: Record<ActaType, string> = {
  ordinaria: 'Ordinaria',
  extraordinaria: 'Extraordinaria',
}

export const ACTA_TYPE_COLORS: Record<ActaType, string> = {
  ordinaria: 'bg-blue-100 text-blue-800',
  extraordinaria: 'bg-purple-100 text-purple-800',
}

export type DocumentoCategory = 'reglamento' | 'protocolo' | 'formulario' | 'otro'

export interface Documento {
  id: string
  title: string
  description?: string
  file_url: string
  category: DocumentoCategory
  is_public: boolean
  created_by?: string
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface DocumentoListResponse {
  documentos: Documento[]
  total: number
  page: number
  per_page: number
}

export const DOCUMENTO_CATEGORY_LABELS: Record<DocumentoCategory, string> = {
  reglamento: 'Reglamento',
  protocolo: 'Protocolo',
  formulario: 'Formulario',
  otro: 'Otro',
}

export const DOCUMENTO_CATEGORY_COLORS: Record<DocumentoCategory, string> = {
  reglamento: 'bg-blue-100 text-blue-800',
  protocolo: 'bg-yellow-100 text-yellow-800',
  formulario: 'bg-green-100 text-green-800',
  otro: 'bg-gray-100 text-gray-800',
}

export type UserRole = 'visitor' | 'vecino' | 'directiva' | 'admin'

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  parcela_id?: number
  email_verified: boolean
  phone?: string
  avatar_url?: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

export interface Parcela {
  id: number
  numero: string
  direccion?: string
  superficie_m2?: number
}

export type VotacionStatus = 'draft' | 'active' | 'closed' | 'cancelled'

// Payload helpers for create/update
export interface VotacionOpcionInput {
  label: string
  description?: string
  order_index: number
}

export interface VotacionUpsertPayload {
  title: string
  description?: string
  requires_quorum: boolean
  quorum_percentage: number
  allow_abstention: boolean
  opciones: VotacionOpcionInput[]
}

export interface Votacion {
  id: string
  title: string
  description?: string
  status: VotacionStatus
  start_date?: string
  end_date?: string
  requires_quorum: boolean
  quorum_percentage: number
  allow_abstention: boolean
  opciones: VotacionOpcion[]
  created_by?: string
  creator_name?: string
  created_at: string
  updated_at: string
  total_votos?: number
  has_voted?: boolean
}

export interface VotacionOpcion {
  id: string
  votacion_id: string
  label: string
  description?: string
  order_index: number
  votos_count?: number
}

export interface VotacionResultado {
  votacion: Votacion
  total_votos: number
  total_abstenciones: number
  resultados: { opcion_id: string; label: string; count: number; percentage: number }[]
  quorum_alcanzado: boolean
  total_vecinos: number
  participacion: number
}

export interface VotacionListResponse {
  votaciones: Votacion[]
  total: number
  page: number
  per_page: number
}

export const VOTACION_STATUS_LABELS: Record<VotacionStatus, string> = {
  draft: 'Borrador',
  active: 'Activa',
  closed: 'Cerrada',
  cancelled: 'Cancelada',
}

export const VOTACION_STATUS_COLORS: Record<VotacionStatus, string> = {
  draft: 'bg-gray-100 text-gray-800',
  active: 'bg-green-100 text-green-800',
  closed: 'bg-blue-100 text-blue-800',
  cancelled: 'bg-red-100 text-red-800',
}

export type EmergenciaPriority = 'low' | 'medium' | 'high' | 'critical'
export type EmergenciaStatus = 'active' | 'resolved' | 'expired'

export interface Emergencia {
  id: string
  title: string
  content: string
  priority: EmergenciaPriority
  status: EmergenciaStatus
  expires_at?: string
  notify_email: boolean
  notify_push: boolean
  created_by?: string
  creator_name?: string
  resolved_at?: string
  resolved_by?: string
  created_at: string
  updated_at: string
}

export interface EmergenciaListResponse {
  emergencias: Emergencia[]
  total: number
  page: number
  per_page: number
}

export const EMERGENCIA_PRIORITY_LABELS: Record<EmergenciaPriority, string> = {
  low: 'Baja',
  medium: 'Media',
  high: 'Alta',
  critical: 'Critica',
}

export const EMERGENCIA_PRIORITY_COLORS: Record<EmergenciaPriority, string> = {
  low: 'bg-blue-100 text-blue-800',
  medium: 'bg-yellow-100 text-yellow-800',
  high: 'bg-orange-100 text-orange-800',
  critical: 'bg-red-100 text-red-800',
}

export const EMERGENCIA_STATUS_LABELS: Record<EmergenciaStatus, string> = {
  active: 'Activa',
  resolved: 'Resuelta',
  expired: 'Expirada',
}

export const EMERGENCIA_STATUS_COLORS: Record<EmergenciaStatus, string> = {
  active: 'bg-green-100 text-green-800',
  resolved: 'bg-blue-100 text-blue-800',
  expired: 'bg-gray-100 text-gray-800',
}

export interface Galeria {
  id: string
  title: string
  description?: string
  event_date?: string
  is_public: boolean
  cover_image_url?: string
  items_count?: number
  created_by?: string
  created_at: string
  updated_at: string
}

export interface GaleriaItem {
  id: string
  galeria_id: string
  file_url: string
  thumbnail_url?: string
  file_type: 'image' | 'video'
  caption?: string
  order_index: number
}

export interface GaleriaListResponse {
  galerias: Galeria[]
  total: number
  page: number
  per_page: number
}

export type PagoStatus = 'pending' | 'paid' | 'overdue' | 'cancelled'

export interface PeriodoGasto {
  id: string
  year: number
  month: number
  monto_base: number
  fecha_vencimiento: string
  descripcion?: string
  created_at?: string
  updated_at?: string
  total_parcelas?: number
  total_pagados?: number
  total_pendientes?: number
  monto_recaudado?: number
  monto_pendiente?: number
}

export interface GastoComun {
  id: string
  periodo_id: string
  parcela_id: number
  parcela_numero?: string
  user_id?: string
  user_name?: string
  monto: number
  monto_pagado: number
  status: PagoStatus
  fecha_pago?: string
  metodo_pago?: string
  referencia_pago?: string
  created_at?: string
  updated_at?: string
  periodo?: PeriodoGasto
}

export interface ResumenGastos {
  periodo: PeriodoGasto
  total_parcelas: number
  total_pagados: number
  total_pendientes: number
  total_vencidos: number
  monto_total: number
  monto_recaudado: number
  monto_pendiente: number
  porcentaje_recaudo: number
}

export interface MiEstadoCuenta {
  parcela_id: number
  parcela_numero: string
  gastos_pendientes: GastoComun[]
  gastos_pagados: GastoComun[]
  total_pendiente: number
  total_pagado: number
}

export const PAGO_STATUS_LABELS: Record<PagoStatus, string> = {
  pending: 'Pendiente',
  paid: 'Pagado',
  overdue: 'Vencido',
  cancelled: 'Cancelado',
}

export const PAGO_STATUS_COLORS: Record<PagoStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  paid: 'bg-green-100 text-green-800',
  overdue: 'bg-red-100 text-red-800',
  cancelled: 'bg-gray-100 text-gray-800',
}

export type ContactoStatus = 'pending' | 'read' | 'replied' | 'archived'

export interface MensajeContacto {
  id: string
  user_id?: string
  nombre: string
  email: string
  asunto: string
  mensaje: string
  status: ContactoStatus
  read_at?: string
  read_by?: string
  replied_at?: string
  replied_by?: string
  respuesta?: string
  created_at: string
  updated_at: string
}

export interface MensajeContactoListResponse {
  mensajes: MensajeContacto[]
  total: number
  page: number
  per_page: number
}

export const CONTACTO_STATUS_LABELS: Record<ContactoStatus, string> = {
  pending: 'Pendiente',
  read: 'Leido',
  replied: 'Respondido',
  archived: 'Archivado',
}

export const CONTACTO_STATUS_COLORS: Record<ContactoStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  read: 'bg-blue-100 text-blue-800',
  replied: 'bg-green-100 text-green-800',
  archived: 'bg-gray-100 text-gray-800',
}

export interface Notificacion {
  id: string
  user_id: string
  title: string
  body: string
  type: string
  reference_id?: string
  is_read: boolean
  read_at?: string
  created_at: string
}

export type AreaType = 'parcela' | 'area_comun' | 'acceso' | 'canal' | 'camino'

export interface MapaArea {
  id: string
  parcela_id?: number
  type: AreaType
  name: string
  description?: string
  coordinates: number[][][]
  center_lat?: number
  center_lng?: number
  fill_color: string
  stroke_color: string
  is_clickable: boolean
}

export interface MapaPunto {
  id: string
  name: string
  description?: string
  lat: number
  lng: number
  icon: string
  type: string
  is_public: boolean
}
