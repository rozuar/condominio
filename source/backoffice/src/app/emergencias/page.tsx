'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getEmergencias, createEmergencia, updateEmergencia, resolveEmergencia, deleteEmergencia } from '@/lib/api'
import type { Emergencia, EmergenciaPriority, EmergenciaStatus } from '@/types'
import { EMERGENCIA_PRIORITY_LABELS, EMERGENCIA_PRIORITY_COLORS, EMERGENCIA_STATUS_LABELS, EMERGENCIA_STATUS_COLORS } from '@/types'
import { Plus, Pencil, Trash2, Loader2, CheckCircle, AlertTriangle } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'

const priorityOptions = [
  { value: 'low', label: 'Baja' },
  { value: 'medium', label: 'Media' },
  { value: 'high', label: 'Alta' },
  { value: 'critical', label: 'Critica' },
]

export default function EmergenciasPage() {
  const { getToken } = useAuth()
  const [emergencias, setEmergencias] = useState<Emergencia[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [showResolveDialog, setShowResolveDialog] = useState(false)
  const [selectedEmergencia, setSelectedEmergencia] = useState<Emergencia | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  const [formData, setFormData] = useState({
    title: '',
    content: '',
    priority: 'medium' as EmergenciaPriority,
    notify_email: true,
    notify_push: true,
  })

  const fetchEmergencias = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getEmergencias(token, { page, per_page: 10 })
      setEmergencias(data.emergencias)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchEmergencias()
  }, [page])

  const openCreateModal = () => {
    setSelectedEmergencia(null)
    setFormData({ title: '', content: '', priority: 'medium', notify_email: true, notify_push: true })
    setShowModal(true)
  }

  const openEditModal = (emergencia: Emergencia) => {
    setSelectedEmergencia(emergencia)
    setFormData({
      title: emergencia.title,
      content: emergencia.content,
      priority: emergencia.priority,
      notify_email: emergencia.notify_email,
      notify_push: emergencia.notify_push,
    })
    setShowModal(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      if (selectedEmergencia) {
        await updateEmergencia(selectedEmergencia.id, token, formData)
      } else {
        await createEmergencia(token, formData)
      }
      setShowModal(false)
      fetchEmergencias()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar')
    } finally {
      setIsSaving(false)
    }
  }

  const handleResolve = async () => {
    if (!selectedEmergencia) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await resolveEmergencia(selectedEmergencia.id, token)
      setShowResolveDialog(false)
      fetchEmergencias()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al resolver')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedEmergencia) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteEmergencia(selectedEmergencia.id, token)
      setShowDeleteDialog(false)
      fetchEmergencias()
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
          <h1 className="text-2xl font-bold text-gray-900">Emergencias</h1>
          <p className="text-gray-500 mt-1">{total} emergencias en total</p>
        </div>
        <Button onClick={openCreateModal} icon={<Plus size={20} />} variant="danger">
          Nueva Emergencia
        </Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
        ) : emergencias.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay emergencias</div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Emergencia</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prioridad</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {emergencias.map((emergencia) => (
                <tr key={emergencia.id} className={`hover:bg-gray-50 ${emergencia.status !== 'active' ? 'opacity-60' : ''}`}>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <AlertTriangle className={`h-5 w-5 ${emergencia.priority === 'critical' ? 'text-red-600' : 'text-amber-500'}`} />
                      <div>
                        <div className="text-sm font-medium text-gray-900">{emergencia.title}</div>
                        <div className="text-sm text-gray-500 truncate max-w-md">{emergencia.content.substring(0, 60)}...</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={EMERGENCIA_PRIORITY_COLORS[emergencia.priority]}>{EMERGENCIA_PRIORITY_LABELS[emergencia.priority]}</Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={EMERGENCIA_STATUS_COLORS[emergencia.status]}>{EMERGENCIA_STATUS_LABELS[emergencia.status]}</Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {format(new Date(emergencia.created_at), "d MMM yyyy, HH:mm", { locale: es })}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex items-center justify-end gap-2">
                      {emergencia.status === 'active' && (
                        <button onClick={() => { setSelectedEmergencia(emergencia); setShowResolveDialog(true); }} className="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg" title="Marcar como resuelta">
                          <CheckCircle size={18} />
                        </button>
                      )}
                      <button onClick={() => openEditModal(emergencia)} className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg">
                        <Pencil size={18} />
                      </button>
                      <button onClick={() => { setSelectedEmergencia(emergencia); setShowDeleteDialog(true); }} className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg">
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

      <Modal open={showModal} onClose={() => setShowModal(false)} title={selectedEmergencia ? 'Editar Emergencia' : 'Nueva Emergencia'} size="lg">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Select label="Prioridad" value={formData.priority} onChange={(e) => setFormData({ ...formData, priority: e.target.value as EmergenciaPriority })} options={priorityOptions} />
          <Textarea label="Contenido" value={formData.content} onChange={(e) => setFormData({ ...formData, content: e.target.value })} required rows={4} />
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <input type="checkbox" id="notify_email" checked={formData.notify_email} onChange={(e) => setFormData({ ...formData, notify_email: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
              <label htmlFor="notify_email" className="text-sm text-gray-700">Notificar por email</label>
            </div>
            <div className="flex items-center gap-2">
              <input type="checkbox" id="notify_push" checked={formData.notify_push} onChange={(e) => setFormData({ ...formData, notify_push: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
              <label htmlFor="notify_push" className="text-sm text-gray-700">Notificar push</label>
            </div>
          </div>
          <div className="flex gap-3 pt-4">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" variant="danger" loading={isSaving} className="flex-1">{selectedEmergencia ? 'Guardar' : 'Crear'}</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog open={showResolveDialog} onClose={() => setShowResolveDialog(false)} onConfirm={handleResolve} title="Resolver Emergencia" message="¿Marcar esta emergencia como resuelta?" confirmText="Resolver" isLoading={isSaving} variant="warning" />
      <ConfirmDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} onConfirm={handleDelete} title="Eliminar Emergencia" message={`¿Eliminar "${selectedEmergencia?.title}"?`} confirmText="Eliminar" isLoading={isSaving} />
    </div>
  )
}
