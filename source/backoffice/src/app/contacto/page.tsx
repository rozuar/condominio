'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getMensajes, getMensaje, markMensajeAsRead, replyMensaje, archiveMensaje, deleteMensaje } from '@/lib/api'
import type { MensajeContacto } from '@/types'
import { CONTACTO_STATUS_LABELS, CONTACTO_STATUS_COLORS } from '@/types'
import { Loader2, Mail, MailOpen, Reply, Archive, Trash2 } from 'lucide-react'
import { format, formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import Textarea from '@/components/ui/Textarea'

export default function ContactoPage() {
  const { getToken } = useAuth()
  const [mensajes, setMensajes] = useState<MensajeContacto[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('')

  const [selectedMensaje, setSelectedMensaje] = useState<MensajeContacto | null>(null)
  const [showDetailModal, setShowDetailModal] = useState(false)
  const [showReplyModal, setShowReplyModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [replyText, setReplyText] = useState('')
  const [isSaving, setIsSaving] = useState(false)

  const fetchMensajes = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getMensajes(token, { page, per_page: 10, status: statusFilter || undefined })
      setMensajes(data.mensajes)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchMensajes()
  }, [page, statusFilter])

  const openDetail = async (mensaje: MensajeContacto) => {
    const token = getToken()
    if (!token) return

    try {
      const detail = await getMensaje(mensaje.id, token)
      setSelectedMensaje(detail)
      setShowDetailModal(true)
      if (detail.status === 'pending') {
        await markMensajeAsRead(mensaje.id, token)
        fetchMensajes()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  const handleReply = async () => {
    if (!selectedMensaje || !replyText.trim()) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await replyMensaje(selectedMensaje.id, token, replyText)
      setShowReplyModal(false)
      setReplyText('')
      setShowDetailModal(false)
      fetchMensajes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al responder')
    } finally {
      setIsSaving(false)
    }
  }

  const handleArchive = async (mensaje: MensajeContacto) => {
    const token = getToken()
    if (!token) return

    try {
      await archiveMensaje(mensaje.id, token)
      fetchMensajes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  const handleDelete = async () => {
    if (!selectedMensaje) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteMensaje(selectedMensaje.id, token)
      setShowDeleteDialog(false)
      setShowDetailModal(false)
      fetchMensajes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar')
    } finally {
      setIsSaving(false)
    }
  }

  const totalPages = Math.ceil(total / 10)
  const filters = [
    { value: '', label: 'Todos' },
    { value: 'pending', label: 'Pendientes' },
    { value: 'read', label: 'Leidos' },
    { value: 'replied', label: 'Respondidos' },
    { value: 'archived', label: 'Archivados' },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Mensajes de Contacto</h1>
          <p className="text-gray-500 mt-1">{total} mensajes</p>
        </div>
      </div>

      {/* Filters */}
      <div className="flex gap-2 mb-6">
        {filters.map((f) => (
          <button
            key={f.value}
            onClick={() => { setStatusFilter(f.value); setPage(1); }}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${statusFilter === f.value ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
        ) : mensajes.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay mensajes</div>
        ) : (
          <div className="divide-y divide-gray-100">
            {mensajes.map((mensaje) => (
              <div
                key={mensaje.id}
                onClick={() => openDetail(mensaje)}
                className={`p-4 hover:bg-gray-50 cursor-pointer transition-colors ${mensaje.status === 'pending' ? 'bg-blue-50' : ''}`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-3">
                    <div className={`p-2 rounded-full ${mensaje.status === 'pending' ? 'bg-blue-100' : 'bg-gray-100'}`}>
                      {mensaje.status === 'pending' ? <Mail className="h-5 w-5 text-blue-600" /> : <MailOpen className="h-5 w-5 text-gray-400" />}
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className={`font-medium ${mensaje.status === 'pending' ? 'text-gray-900' : 'text-gray-700'}`}>{mensaje.nombre}</p>
                        <Badge color={CONTACTO_STATUS_COLORS[mensaje.status]}>{CONTACTO_STATUS_LABELS[mensaje.status]}</Badge>
                      </div>
                      <p className="text-sm text-gray-500">{mensaje.email}</p>
                      <p className="text-sm font-medium text-gray-800 mt-1">{mensaje.asunto}</p>
                      <p className="text-sm text-gray-500 mt-1 line-clamp-2">{mensaje.mensaje}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-xs text-gray-400">{formatDistanceToNow(new Date(mensaje.created_at), { addSuffix: true, locale: es })}</p>
                    <div className="mt-2 flex gap-1 justify-end" onClick={(e) => e.stopPropagation()}>
                      {mensaje.status !== 'archived' && (
                        <button onClick={() => handleArchive(mensaje)} className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded" title="Archivar">
                          <Archive size={16} />
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
            <p className="text-sm text-gray-500">Pagina {page} de {totalPages}</p>
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" onClick={() => setPage(page - 1)} disabled={page === 1}>Anterior</Button>
              <Button variant="secondary" size="sm" onClick={() => setPage(page + 1)} disabled={page === totalPages}>Siguiente</Button>
            </div>
          </div>
        )}
      </div>

      {/* Detail Modal */}
      <Modal open={showDetailModal} onClose={() => setShowDetailModal(false)} title="Detalle del Mensaje" size="lg">
        {selectedMensaje && (
          <div className="space-y-6">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div><span className="text-gray-500">De:</span> <span className="font-medium">{selectedMensaje.nombre}</span></div>
              <div><span className="text-gray-500">Email:</span> <span className="font-medium">{selectedMensaje.email}</span></div>
              <div><span className="text-gray-500">Fecha:</span> <span>{format(new Date(selectedMensaje.created_at), "d MMMM yyyy, HH:mm", { locale: es })}</span></div>
              <div><span className="text-gray-500">Estado:</span> <Badge color={CONTACTO_STATUS_COLORS[selectedMensaje.status]}>{CONTACTO_STATUS_LABELS[selectedMensaje.status]}</Badge></div>
            </div>

            <div>
              <h3 className="font-medium text-gray-900 mb-2">{selectedMensaje.asunto}</h3>
              <div className="p-4 bg-gray-50 rounded-lg text-gray-700 whitespace-pre-wrap">{selectedMensaje.mensaje}</div>
            </div>

            {selectedMensaje.respuesta && (
              <div>
                <h4 className="font-medium text-gray-900 mb-2">Respuesta enviada</h4>
                <div className="p-4 bg-green-50 rounded-lg text-gray-700 whitespace-pre-wrap">{selectedMensaje.respuesta}</div>
              </div>
            )}

            <div className="flex gap-3 pt-4 border-t">
              {selectedMensaje.status !== 'replied' && (
                <Button onClick={() => setShowReplyModal(true)} icon={<Reply size={18} />}>Responder</Button>
              )}
              {selectedMensaje.status !== 'archived' && (
                <Button variant="secondary" onClick={() => handleArchive(selectedMensaje)} icon={<Archive size={18} />}>Archivar</Button>
              )}
              <Button variant="danger" onClick={() => setShowDeleteDialog(true)} icon={<Trash2 size={18} />}>Eliminar</Button>
            </div>
          </div>
        )}
      </Modal>

      {/* Reply Modal */}
      <Modal open={showReplyModal} onClose={() => setShowReplyModal(false)} title="Responder Mensaje" size="lg">
        <div className="space-y-6">
          <p className="text-sm text-gray-600">Respondiendo a: <strong>{selectedMensaje?.nombre}</strong> ({selectedMensaje?.email})</p>
          <Textarea label="Tu respuesta" value={replyText} onChange={(e) => setReplyText(e.target.value)} rows={6} placeholder="Escribe tu respuesta..." />
          <div className="flex gap-3">
            <Button variant="secondary" onClick={() => setShowReplyModal(false)} className="flex-1">Cancelar</Button>
            <Button onClick={handleReply} loading={isSaving} className="flex-1">Enviar Respuesta</Button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} onConfirm={handleDelete} title="Eliminar Mensaje" message="¿Eliminar este mensaje?" confirmText="Eliminar" isLoading={isSaving} />
    </div>
  )
}
