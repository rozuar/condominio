'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getEventos, createEvento, updateEvento, deleteEvento } from '@/lib/api'
import type { Evento, EventoType } from '@/types'
import { EVENTO_TYPE_LABELS, EVENTO_TYPE_COLORS } from '@/types'
import { Plus, Pencil, Trash2, Loader2, MapPin, Clock } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'

const typeOptions = [
  { value: 'reunion', label: 'Reunion' },
  { value: 'asamblea', label: 'Asamblea' },
  { value: 'trabajo', label: 'Trabajo Comunitario' },
  { value: 'social', label: 'Actividad Social' },
]

export default function EventosPage() {
  const { getToken } = useAuth()
  const [eventos, setEventos] = useState<Evento[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [selectedEvento, setSelectedEvento] = useState<Evento | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    event_date: '',
    location: '',
    type: 'reunion' as EventoType,
    is_public: true,
  })

  const fetchEventos = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getEventos(token, { page, per_page: 10 })
      setEventos(data.eventos)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar eventos')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchEventos()
  }, [page])

  const openCreateModal = () => {
    setSelectedEvento(null)
    setFormData({ title: '', description: '', event_date: '', location: '', type: 'reunion', is_public: true })
    setShowModal(true)
  }

  const openEditModal = (evento: Evento) => {
    setSelectedEvento(evento)
    setFormData({
      title: evento.title,
      description: evento.description || '',
      event_date: evento.event_date.slice(0, 16),
      location: evento.location || '',
      type: evento.type,
      is_public: evento.is_public,
    })
    setShowModal(true)
  }

  const openDeleteDialog = (evento: Evento) => {
    setSelectedEvento(evento)
    setShowDeleteDialog(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      if (selectedEvento) {
        await updateEvento(selectedEvento.id, token, formData)
      } else {
        await createEvento(token, formData)
      }
      setShowModal(false)
      fetchEventos()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedEvento) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteEvento(selectedEvento.id, token)
      setShowDeleteDialog(false)
      setSelectedEvento(null)
      fetchEventos()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar')
    } finally {
      setIsSaving(false)
    }
  }

  const totalPages = Math.ceil(total / 10)

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Eventos</h1>
          <p className="text-gray-500 mt-1">{total} eventos en total</p>
        </div>
        <Button onClick={openCreateModal} icon={<Plus size={20} />}>
          Nuevo Evento
        </Button>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          {error}
        </div>
      )}

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        ) : eventos.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay eventos</div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Evento</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Ubicacion</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {eventos.map((evento) => (
                <tr key={evento.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="text-sm font-medium text-gray-900">{evento.title}</div>
                    {evento.description && (
                      <div className="text-sm text-gray-500 truncate max-w-md">{evento.description.substring(0, 80)}...</div>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={EVENTO_TYPE_COLORS[evento.type]}>{EVENTO_TYPE_LABELS[evento.type]}</Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center gap-1 text-sm text-gray-500">
                      <Clock size={14} />
                      {format(new Date(evento.event_date), "d MMM yyyy, HH:mm", { locale: es })}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {evento.location && (
                      <div className="flex items-center gap-1 text-sm text-gray-500">
                        <MapPin size={14} />
                        {evento.location}
                      </div>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button onClick={() => openEditModal(evento)} className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg">
                        <Pencil size={18} />
                      </button>
                      <button onClick={() => openDeleteDialog(evento)} className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
                        <Trash2 size={18} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
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

      <Modal open={showModal} onClose={() => setShowModal(false)} title={selectedEvento ? 'Editar Evento' : 'Nuevo Evento'} size="lg">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Select label="Tipo" value={formData.type} onChange={(e) => setFormData({ ...formData, type: e.target.value as EventoType })} options={typeOptions} />
          <Input label="Fecha y Hora" type="datetime-local" value={formData.event_date} onChange={(e) => setFormData({ ...formData, event_date: e.target.value })} required />
          <Input label="Ubicacion" value={formData.location} onChange={(e) => setFormData({ ...formData, location: e.target.value })} placeholder="Ej: Sede social, Sector B" />
          <Textarea label="Descripcion" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={4} />
          <div className="flex items-center gap-2">
            <input type="checkbox" id="is_public" checked={formData.is_public} onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
            <label htmlFor="is_public" className="text-sm text-gray-700">Evento publico</label>
          </div>
          <div className="flex gap-3 pt-4">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">{selectedEvento ? 'Guardar' : 'Crear'}</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} onConfirm={handleDelete} title="Eliminar Evento" message={`¿Eliminar "${selectedEvento?.title}"?`} confirmText="Eliminar" isLoading={isSaving} />
    </div>
  )
}
