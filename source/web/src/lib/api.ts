const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_URL}${endpoint}`

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`)
  }

  return response.json()
}

// Comunicados
export async function getComunicados(params?: {
  page?: number
  per_page?: number
  type?: string
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)

  const query = searchParams.toString()
  return fetchAPI<{ comunicados: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/comunicados${query ? `?${query}` : ''}`
  )
}

export async function getComunicado(id: string) {
  return fetchAPI<any>(`/api/v1/comunicados/${id}`)
}

export async function getLatestComunicados(limit = 3) {
  return fetchAPI<{ comunicados: any[] }>(`/api/v1/comunicados/latest?limit=${limit}`)
}

// Eventos
export async function getEventos(params?: {
  page?: number
  per_page?: number
  type?: string
  upcoming?: boolean
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  if (params?.upcoming) searchParams.set('upcoming', 'true')

  const query = searchParams.toString()
  return fetchAPI<{ eventos: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/eventos${query ? `?${query}` : ''}`
  )
}

export async function getEvento(id: string) {
  return fetchAPI<any>(`/api/v1/eventos/${id}`)
}

export async function getUpcomingEventos(limit = 3) {
  return fetchAPI<{ eventos: any[] }>(`/api/v1/eventos/upcoming?limit=${limit}`)
}

// Authenticated fetch helper
async function fetchAPIAuth<T>(endpoint: string, token: string, options?: RequestInit): Promise<T> {
  const url = `${API_URL}${endpoint}`

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options?.headers,
    },
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('No autorizado')
    }
    throw new Error(`API Error: ${response.status}`)
  }

  return response.json()
}

// Tesorería
export async function getTesoreriaResumen(token: string) {
  return fetchAPIAuth<{
    total_ingresos: number
    total_egresos: number
    balance: number
    movimientos_count: number
  }>('/api/v1/tesoreria/resumen', token)
}

export async function getMovimientos(token: string, params?: {
  page?: number
  per_page?: number
  type?: string
  year?: number
  month?: number
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  if (params?.year) searchParams.set('year', params.year.toString())
  if (params?.month) searchParams.set('month', params.month.toString())

  const query = searchParams.toString()
  return fetchAPIAuth<{ movimientos: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/tesoreria${query ? `?${query}` : ''}`,
    token
  )
}

// Actas
export async function getActas(token: string, params?: {
  page?: number
  per_page?: number
  type?: string
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)

  const query = searchParams.toString()
  return fetchAPIAuth<{ actas: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/actas${query ? `?${query}` : ''}`,
    token
  )
}

export async function getActa(id: string, token: string) {
  return fetchAPIAuth<any>(`/api/v1/actas/${id}`, token)
}

// Documentos
export async function getDocumentos(token: string, params?: {
  page?: number
  per_page?: number
  category?: string
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.category) searchParams.set('category', params.category)

  const query = searchParams.toString()
  return fetchAPIAuth<{ documentos: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/documentos${query ? `?${query}` : ''}`,
    token
  )
}

export async function getDocumento(id: string, token: string) {
  return fetchAPIAuth<any>(`/api/v1/documentos/${id}`, token)
}

// Emergencias
export async function getEmergencias(params?: {
  page?: number
  per_page?: number
  status?: string
  priority?: string
  active?: boolean
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.status) searchParams.set('status', params.status)
  if (params?.priority) searchParams.set('priority', params.priority)
  if (params?.active) searchParams.set('active', 'true')

  const query = searchParams.toString()
  return fetchAPI<{ emergencias: any[]; total: number; page: number; per_page: number }>(
    `/api/v1/emergencias${query ? `?${query}` : ''}`
  )
}

export async function getActiveEmergencias() {
  return fetchAPI<{ emergencias: any[]; total: number }>('/api/v1/emergencias/active')
}

export async function getEmergencia(id: string) {
  return fetchAPI<any>(`/api/v1/emergencias/${id}`)
}

// Votaciones
import type { Votacion, VotacionListResponse, VotacionResultado } from '@/types'

export async function getVotaciones(token: string, params?: {
  page?: number
  per_page?: number
  status?: string
  active?: boolean
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.status) searchParams.set('status', params.status)
  if (params?.active) searchParams.set('active', 'true')

  const query = searchParams.toString()
  return fetchAPIAuth<VotacionListResponse>(
    `/api/v1/votaciones${query ? `?${query}` : ''}`,
    token
  )
}

export async function getActiveVotaciones(token: string) {
  return fetchAPIAuth<{ votaciones: Votacion[]; total: number }>(
    '/api/v1/votaciones/active',
    token
  )
}

export async function getVotacion(id: string, token: string) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}`, token)
}

export async function getVotacionResultados(id: string, token: string) {
  return fetchAPIAuth<VotacionResultado>(`/api/v1/votaciones/${id}/resultados`, token)
}

export async function emitirVoto(id: string, token: string, data: { opcion_id?: string; is_abstention: boolean }) {
  return fetchAPIAuth<{ message: string }>(
    `/api/v1/votaciones/${id}/votar`,
    token,
    {
      method: 'POST',
      body: JSON.stringify(data),
    }
  )
}

// Gastos Comunes
import type { PeriodoGasto, GastoComun, ResumenGastos, MiEstadoCuenta, PagoStatus } from '@/types'

export interface PeriodoListResponse {
  periodos: PeriodoGasto[]
  total: number
  page: number
  per_page: number
}

export interface GastoComunListResponse {
  gastos: GastoComun[]
  total: number
  page: number
  per_page: number
}

export async function getPeriodos(token: string, params?: {
  year?: number
  page?: number
  per_page?: number
}) {
  const searchParams = new URLSearchParams()
  if (params?.year) searchParams.set('year', params.year.toString())
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())

  const query = searchParams.toString()
  return fetchAPIAuth<PeriodoListResponse>(
    `/api/v1/gastos/periodos${query ? `?${query}` : ''}`,
    token
  )
}

export async function getPeriodoActual(token: string) {
  return fetchAPIAuth<PeriodoGasto>('/api/v1/gastos/periodos/actual', token)
}

export async function getPeriodo(id: string, token: string) {
  return fetchAPIAuth<PeriodoGasto>(`/api/v1/gastos/periodos/${id}`, token)
}

export async function getResumenPeriodo(id: string, token: string) {
  return fetchAPIAuth<ResumenGastos>(`/api/v1/gastos/periodos/${id}/resumen`, token)
}

export async function getGastosPeriodo(id: string, token: string, params?: {
  status?: PagoStatus
  page?: number
  per_page?: number
}) {
  const searchParams = new URLSearchParams()
  searchParams.set('periodo_id', id)
  if (params?.status) searchParams.set('status', params.status)
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())

  return fetchAPIAuth<GastoComunListResponse>(
    `/api/v1/gastos/periodos/${id}/gastos?${searchParams.toString()}`,
    token
  )
}

export async function getMiEstadoCuenta(token: string) {
  return fetchAPIAuth<MiEstadoCuenta>('/api/v1/gastos/mi-cuenta', token)
}

export async function getGasto(id: string, token: string) {
  return fetchAPIAuth<GastoComun>(`/api/v1/gastos/${id}`, token)
}

// Contacto
import type { MensajeContacto, ContactoStatus } from '@/types'

export interface MensajeContactoListResponse {
  mensajes: MensajeContacto[]
  total: number
  page: number
  per_page: number
}

export interface ContactoStats {
  total_pending: number
  total_read: number
  total_replied: number
  total_archived: number
  total: number
}

export async function enviarMensajeContacto(data: {
  nombre: string
  email: string
  asunto: string
  mensaje: string
}) {
  return fetchAPI<MensajeContacto>('/api/v1/contacto', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function getMisMensajes(token: string) {
  return fetchAPIAuth<{ mensajes: MensajeContacto[]; total: number }>(
    '/api/v1/contacto/mis-mensajes',
    token
  )
}

// Galeria
import type { Galeria, GaleriaItem } from '@/types'

export interface GaleriaListResponse {
  galerias: Galeria[]
  total: number
  page: number
  per_page: number
}

export interface GaleriaWithItems extends Galeria {
  items: GaleriaItem[]
}

export async function getGalerias(params?: {
  page?: number
  per_page?: number
  is_public?: boolean
}) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.is_public !== undefined) searchParams.set('is_public', params.is_public.toString())

  const query = searchParams.toString()
  return fetchAPI<GaleriaListResponse>(
    `/api/v1/galerias${query ? `?${query}` : ''}`
  )
}

export async function getGaleria(id: string) {
  return fetchAPI<GaleriaWithItems>(`/api/v1/galerias/${id}`)
}

// Mapa
import type { MapaArea, MapaPunto } from '@/types'

export interface MapaData {
  areas: MapaArea[]
  puntos: MapaPunto[]
}

export async function getMapaData() {
  return fetchAPI<MapaData>('/api/v1/mapa')
}

export async function getMapaAreas(params?: {
  type?: string
  page?: number
  per_page?: number
}) {
  const searchParams = new URLSearchParams()
  if (params?.type) searchParams.set('type', params.type)
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())

  const query = searchParams.toString()
  return fetchAPI<{ areas: MapaArea[]; total: number; page: number; per_page: number }>(
    `/api/v1/mapa/areas${query ? `?${query}` : ''}`
  )
}

export async function getMapaPuntos(params?: {
  type?: string
  is_public?: boolean
  page?: number
  per_page?: number
}) {
  const searchParams = new URLSearchParams()
  if (params?.type) searchParams.set('type', params.type)
  if (params?.is_public !== undefined) searchParams.set('is_public', params.is_public.toString())
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())

  const query = searchParams.toString()
  return fetchAPI<{ puntos: MapaPunto[]; total: number; page: number; per_page: number }>(
    `/api/v1/mapa/puntos${query ? `?${query}` : ''}`
  )
}
