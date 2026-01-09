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
    const error = await response.json().catch(() => ({ message: `Error ${response.status}` }))
    throw new Error(error.message || error.error || `API Error: ${response.status}`)
  }

  return response.json()
}

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
    const error = await response.json().catch(() => ({ message: `Error ${response.status}` }))
    throw new Error(error.message || error.error || `API Error: ${response.status}`)
  }

  return response.json()
}

// Auth
import type { AuthResponse, User } from '@/types'

export async function login(email: string, password: string) {
  return fetchAPI<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function getMe(token: string) {
  // API returns the user object directly
  return fetchAPIAuth<User>('/auth/me', token)
}

export async function refreshTokens(refreshToken: string) {
  // API returns the same shape as /auth/login
  return fetchAPI<AuthResponse>('/auth/refresh', {
    method: 'POST',
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
}

// Comunicados
import type { Comunicado, ComunicadoListResponse } from '@/types'

export async function getComunicados(token: string, params?: { page?: number; per_page?: number; type?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  const query = searchParams.toString()
  return fetchAPIAuth<ComunicadoListResponse>(`/api/v1/comunicados${query ? `?${query}` : ''}`, token)
}

export async function getComunicado(id: string, token: string) {
  return fetchAPIAuth<Comunicado>(`/api/v1/comunicados/${id}`, token)
}

export async function createComunicado(token: string, data: Partial<Comunicado>) {
  return fetchAPIAuth<Comunicado>('/api/v1/comunicados', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateComunicado(id: string, token: string, data: Partial<Comunicado>) {
  return fetchAPIAuth<Comunicado>(`/api/v1/comunicados/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteComunicado(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/comunicados/${id}`, token, {
    method: 'DELETE',
  })
}

// Eventos
import type { Evento, EventoListResponse } from '@/types'

export async function getEventos(token: string, params?: { page?: number; per_page?: number; type?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  const query = searchParams.toString()
  return fetchAPIAuth<EventoListResponse>(`/api/v1/eventos${query ? `?${query}` : ''}`, token)
}

export async function getEvento(id: string, token: string) {
  return fetchAPIAuth<Evento>(`/api/v1/eventos/${id}`, token)
}

export async function createEvento(token: string, data: Partial<Evento>) {
  return fetchAPIAuth<Evento>('/api/v1/eventos', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateEvento(id: string, token: string, data: Partial<Evento>) {
  return fetchAPIAuth<Evento>(`/api/v1/eventos/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteEvento(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/eventos/${id}`, token, {
    method: 'DELETE',
  })
}

// Emergencias
import type { Emergencia, EmergenciaListResponse } from '@/types'

export async function getEmergencias(token: string, params?: { page?: number; per_page?: number; status?: string; priority?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.status) searchParams.set('status', params.status)
  if (params?.priority) searchParams.set('priority', params.priority)
  const query = searchParams.toString()
  return fetchAPIAuth<EmergenciaListResponse>(`/api/v1/emergencias${query ? `?${query}` : ''}`, token)
}

export async function getEmergencia(id: string, token: string) {
  return fetchAPIAuth<Emergencia>(`/api/v1/emergencias/${id}`, token)
}

export async function createEmergencia(token: string, data: Partial<Emergencia>) {
  return fetchAPIAuth<Emergencia>('/api/v1/emergencias', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateEmergencia(id: string, token: string, data: Partial<Emergencia>) {
  return fetchAPIAuth<Emergencia>(`/api/v1/emergencias/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function resolveEmergencia(id: string, token: string) {
  return fetchAPIAuth<Emergencia>(`/api/v1/emergencias/${id}/resolve`, token, {
    method: 'POST',
  })
}

export async function deleteEmergencia(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/emergencias/${id}`, token, {
    method: 'DELETE',
  })
}

// Votaciones
import type { Votacion, VotacionListResponse, VotacionResultado, VotacionUpsertPayload } from '@/types'

export async function getVotaciones(token: string, params?: { page?: number; per_page?: number; status?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.status) searchParams.set('status', params.status)
  const query = searchParams.toString()
  return fetchAPIAuth<VotacionListResponse>(`/api/v1/votaciones${query ? `?${query}` : ''}`, token)
}

export async function getVotacion(id: string, token: string) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}`, token)
}

export async function getVotacionResultados(id: string, token: string) {
  return fetchAPIAuth<VotacionResultado>(`/api/v1/votaciones/${id}/resultados`, token)
}

export async function createVotacion(token: string, data: VotacionUpsertPayload) {
  return fetchAPIAuth<Votacion>('/api/v1/votaciones', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateVotacion(id: string, token: string, data: VotacionUpsertPayload) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function publishVotacion(id: string, token: string) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}/publish`, token, { method: 'POST' })
}

export async function closeVotacion(id: string, token: string) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}/close`, token, { method: 'POST' })
}

export async function cancelVotacion(id: string, token: string) {
  return fetchAPIAuth<Votacion>(`/api/v1/votaciones/${id}/cancel`, token, { method: 'POST' })
}

export async function deleteVotacion(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/votaciones/${id}`, token, { method: 'DELETE' })
}

// Gastos Comunes
import type { PeriodoGasto, GastoComun, ResumenGastos } from '@/types'

export async function getPeriodos(token: string, params?: { year?: number; page?: number; per_page?: number }) {
  const searchParams = new URLSearchParams()
  if (params?.year) searchParams.set('year', params.year.toString())
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  const query = searchParams.toString()
  return fetchAPIAuth<{ periodos: PeriodoGasto[]; total: number }>(`/api/v1/gastos/periodos${query ? `?${query}` : ''}`, token)
}

export async function getPeriodo(id: string, token: string) {
  return fetchAPIAuth<PeriodoGasto>(`/api/v1/gastos/periodos/${id}`, token)
}

export async function getResumenPeriodo(id: string, token: string) {
  return fetchAPIAuth<ResumenGastos>(`/api/v1/gastos/periodos/${id}/resumen`, token)
}

export async function getGastosPeriodo(id: string, token: string, params?: { status?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.status) searchParams.set('status', params.status)
  const query = searchParams.toString()
  return fetchAPIAuth<{ gastos: GastoComun[]; total: number }>(`/api/v1/gastos/periodos/${id}/gastos${query ? `?${query}` : ''}`, token)
}

export async function createPeriodo(token: string, data: Partial<PeriodoGasto>) {
  return fetchAPIAuth<PeriodoGasto>('/api/v1/gastos/periodos', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updatePeriodo(id: string, token: string, data: Partial<PeriodoGasto>) {
  return fetchAPIAuth<PeriodoGasto>(`/api/v1/gastos/periodos/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function registrarPago(gastoId: string, token: string, data: { monto_pagado: number; metodo_pago?: string; referencia_pago?: string }) {
  return fetchAPIAuth<GastoComun>(`/api/v1/gastos/${gastoId}/pago`, token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

// Contacto
import type { MensajeContacto, MensajeContactoListResponse } from '@/types'

export async function getMensajes(token: string, params?: { page?: number; per_page?: number; status?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.status) searchParams.set('status', params.status)
  const query = searchParams.toString()
  return fetchAPIAuth<MensajeContactoListResponse>(`/api/v1/contacto${query ? `?${query}` : ''}`, token)
}

export async function getMensaje(id: string, token: string) {
  return fetchAPIAuth<MensajeContacto>(`/api/v1/contacto/${id}`, token)
}

export async function markMensajeAsRead(id: string, token: string) {
  return fetchAPIAuth<MensajeContacto>(`/api/v1/contacto/${id}/read`, token, { method: 'POST' })
}

export async function replyMensaje(id: string, token: string, respuesta: string) {
  return fetchAPIAuth<MensajeContacto>(`/api/v1/contacto/${id}/reply`, token, {
    method: 'POST',
    body: JSON.stringify({ respuesta }),
  })
}

export async function archiveMensaje(id: string, token: string) {
  return fetchAPIAuth<MensajeContacto>(`/api/v1/contacto/${id}/archive`, token, { method: 'POST' })
}

export async function deleteMensaje(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/contacto/${id}`, token, { method: 'DELETE' })
}

// Galerias
import type { Galeria, GaleriaItem, GaleriaListResponse } from '@/types'

export async function getGalerias(token: string, params?: { page?: number; per_page?: number }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  const query = searchParams.toString()
  return fetchAPIAuth<GaleriaListResponse>(`/api/v1/galerias${query ? `?${query}` : ''}`, token)
}

export async function getGaleria(id: string, token: string) {
  return fetchAPIAuth<Galeria & { items: GaleriaItem[] }>(`/api/v1/galerias/${id}`, token)
}

export async function createGaleria(token: string, data: Partial<Galeria>) {
  return fetchAPIAuth<Galeria>('/api/v1/galerias', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateGaleria(id: string, token: string, data: Partial<Galeria>) {
  return fetchAPIAuth<Galeria>(`/api/v1/galerias/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteGaleria(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/galerias/${id}`, token, { method: 'DELETE' })
}

// Tesoreria
import type { Movimiento, MovimientoListResponse, TesoreriaResumen } from '@/types'

export async function getTesoreriaResumen(token: string) {
  return fetchAPIAuth<TesoreriaResumen>('/api/v1/tesoreria/resumen', token)
}

export async function getMovimientos(token: string, params?: { page?: number; per_page?: number; type?: string; year?: number; month?: number }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  if (params?.year) searchParams.set('year', params.year.toString())
  if (params?.month) searchParams.set('month', params.month.toString())
  const query = searchParams.toString()
  return fetchAPIAuth<MovimientoListResponse>(`/api/v1/tesoreria${query ? `?${query}` : ''}`, token)
}

export async function createMovimiento(token: string, data: Partial<Movimiento>) {
  return fetchAPIAuth<Movimiento>('/api/v1/tesoreria', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

// Actas
import type { Acta, ActaListResponse } from '@/types'

export async function getActas(token: string, params?: { page?: number; per_page?: number; type?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.type) searchParams.set('type', params.type)
  const query = searchParams.toString()
  return fetchAPIAuth<ActaListResponse>(`/api/v1/actas${query ? `?${query}` : ''}`, token)
}

export async function getActa(id: string, token: string) {
  return fetchAPIAuth<Acta>(`/api/v1/actas/${id}`, token)
}

export async function createActa(token: string, data: Partial<Acta>) {
  return fetchAPIAuth<Acta>('/api/v1/actas', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

// Documentos
import type { Documento, DocumentoListResponse } from '@/types'

export async function getDocumentos(token: string, params?: { page?: number; per_page?: number; category?: string }) {
  const searchParams = new URLSearchParams()
  if (params?.page) searchParams.set('page', params.page.toString())
  if (params?.per_page) searchParams.set('per_page', params.per_page.toString())
  if (params?.category) searchParams.set('category', params.category)
  const query = searchParams.toString()
  return fetchAPIAuth<DocumentoListResponse>(`/api/v1/documentos${query ? `?${query}` : ''}`, token)
}

export async function getDocumento(id: string, token: string) {
  return fetchAPIAuth<Documento>(`/api/v1/documentos/${id}`, token)
}

export async function createDocumento(token: string, data: Partial<Documento>) {
  return fetchAPIAuth<Documento>('/api/v1/documentos', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

// Notificaciones
import type { Notificacion } from '@/types'

export async function createNotificacion(token: string, data: { user_id: string; title: string; body: string; type: string }) {
  return fetchAPIAuth<Notificacion>('/api/v1/notificaciones', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function createBulkNotificaciones(token: string, data: { user_ids: string[]; title: string; body: string; type: string }) {
  return fetchAPIAuth<{ count: number }>('/api/v1/notificaciones/bulk', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function broadcastNotificacion(token: string, data: { title: string; body: string; type: string; roles?: string[] }) {
  return fetchAPIAuth<{ count: number }>('/api/v1/notificaciones/broadcast', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

// Mapa
import type { MapaArea, MapaPunto } from '@/types'

export async function getMapaAreas(token: string) {
  return fetchAPIAuth<{ areas: MapaArea[] }>('/api/v1/mapa/areas', token)
}

export async function createMapaArea(token: string, data: Partial<MapaArea>) {
  return fetchAPIAuth<MapaArea>('/api/v1/mapa/areas', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateMapaArea(id: string, token: string, data: Partial<MapaArea>) {
  return fetchAPIAuth<MapaArea>(`/api/v1/mapa/areas/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteMapaArea(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/mapa/areas/${id}`, token, { method: 'DELETE' })
}

export async function getMapaPuntos(token: string) {
  return fetchAPIAuth<{ puntos: MapaPunto[] }>('/api/v1/mapa/puntos', token)
}

export async function createMapaPunto(token: string, data: Partial<MapaPunto>) {
  return fetchAPIAuth<MapaPunto>('/api/v1/mapa/puntos', token, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateMapaPunto(id: string, token: string, data: Partial<MapaPunto>) {
  return fetchAPIAuth<MapaPunto>(`/api/v1/mapa/puntos/${id}`, token, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteMapaPunto(id: string, token: string) {
  return fetchAPIAuth<{ message: string }>(`/api/v1/mapa/puntos/${id}`, token, { method: 'DELETE' })
}

// Dashboard stats
export async function getDashboardStats(token: string) {
  const [comunicados, eventos, emergencias, votaciones, mensajes] = await Promise.all([
    getComunicados(token, { per_page: 1 }),
    getEventos(token, { per_page: 1 }),
    getEmergencias(token, { per_page: 1, status: 'active' }),
    getVotaciones(token, { per_page: 1, status: 'active' }),
    getMensajes(token, { per_page: 1, status: 'pending' }),
  ])

  return {
    totalComunicados: comunicados.total,
    totalEventos: eventos.total,
    emergenciasActivas: emergencias.total,
    votacionesActivas: votaciones.total,
    mensajesPendientes: mensajes.total,
  }
}
