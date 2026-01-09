'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getVotaciones, getVotacion, createVotacion, updateVotacion, publishVotacion, closeVotacion, cancelVotacion, deleteVotacion, getVotacionResultados } from '@/lib/api'
import type { Votacion, VotacionStatus, VotacionResultado } from '@/types'
import { VOTACION_STATUS_LABELS, VOTACION_STATUS_COLORS } from '@/types'
import { Plus, Pencil, Trash2, Loader2, Play, Square, Ban, BarChart3 } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import Input from '@/components/ui/Input'
import Textarea from '@/components/ui/Textarea'

export default function VotacionesPage() {
  const { getToken } = useAuth()
  const [votaciones, setVotaciones] = useState<Votacion[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [showResultsModal, setShowResultsModal] = useState(false)
  const [actionDialog, setActionDialog] = useState<{ open: boolean; type: 'publish' | 'close' | 'cancel' | null }>({ open: false, type: null })
  const [selectedVotacion, setSelectedVotacion] = useState<Votacion | null>(null)
  const [resultados, setResultados] = useState<VotacionResultado | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    requires_quorum: false,
    quorum_percentage: 50,
    allow_abstention: true,
    opciones: [{ label: '', description: '' }, { label: '', description: '' }],
  })

  const fetchVotaciones = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getVotaciones(token, { page, per_page: 10 })
      setVotaciones(data.votaciones)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchVotaciones()
  }, [page])

  const openCreateModal = () => {
    setSelectedVotacion(null)
    setFormData({ title: '', description: '', requires_quorum: false, quorum_percentage: 50, allow_abstention: true, opciones: [{ label: '', description: '' }, { label: '', description: '' }] })
    setShowModal(true)
  }

  const openEditModal = async (votacion: Votacion) => {
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      const full = await getVotacion(votacion.id, token)
      setSelectedVotacion(full)
      setFormData({
        title: full.title,
        description: full.description || '',
        requires_quorum: full.requires_quorum,
        quorum_percentage: full.quorum_percentage,
        allow_abstention: full.allow_abstention,
        opciones: (full.opciones?.length ? full.opciones : [{ label: '', description: '' }, { label: '', description: '' }]).map((o) => ({
          label: o.label,
          description: o.description || '',
        })),
      })
      setShowModal(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar votacion')
    } finally {
      setIsSaving(false)
    }
  }

  const openResultsModal = async (votacion: Votacion) => {
    const token = getToken()
    if (!token) return

    setSelectedVotacion(votacion)
    try {
      const data = await getVotacionResultados(votacion.id, token)
      setResultados(data)
      setShowResultsModal(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar resultados')
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      const data = {
        ...formData,
        opciones: formData.opciones.filter(o => o.label.trim()).map((o, i) => ({ ...o, order_index: i })),
      }
      if (selectedVotacion) {
        await updateVotacion(selectedVotacion.id, token, data)
      } else {
        await createVotacion(token, data)
      }
      setShowModal(false)
      fetchVotaciones()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar')
    } finally {
      setIsSaving(false)
    }
  }

  const handleAction = async () => {
    if (!selectedVotacion || !actionDialog.type) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      if (actionDialog.type === 'publish') await publishVotacion(selectedVotacion.id, token)
      else if (actionDialog.type === 'close') await closeVotacion(selectedVotacion.id, token)
      else if (actionDialog.type === 'cancel') await cancelVotacion(selectedVotacion.id, token)
      setActionDialog({ open: false, type: null })
      fetchVotaciones()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedVotacion) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteVotacion(selectedVotacion.id, token)
      setShowDeleteDialog(false)
      fetchVotaciones()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar')
    } finally {
      setIsSaving(false)
    }
  }

  const addOpcion = () => setFormData({ ...formData, opciones: [...formData.opciones, { label: '', description: '' }] })
  const removeOpcion = (index: number) => setFormData({ ...formData, opciones: formData.opciones.filter((_, i) => i !== index) })
  const updateOpcion = (index: number, field: string, value: string) => {
    const opciones = [...formData.opciones]
    opciones[index] = { ...opciones[index], [field]: value }
    setFormData({ ...formData, opciones })
  }

  const totalPages = Math.ceil(total / 10)

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Votaciones</h1>
          <p className="text-gray-500 mt-1">{total} votaciones en total</p>
        </div>
        <Button onClick={openCreateModal} icon={<Plus size={20} />}>Nueva Votacion</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
        ) : votaciones.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay votaciones</div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Votacion</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Opciones</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Votos</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {votaciones.map((votacion) => (
                <tr key={votacion.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="text-sm font-medium text-gray-900">{votacion.title}</div>
                    {votacion.description && <div className="text-sm text-gray-500 truncate max-w-md">{votacion.description.substring(0, 60)}...</div>}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={VOTACION_STATUS_COLORS[votacion.status]}>{VOTACION_STATUS_LABELS[votacion.status]}</Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{votacion.opciones?.length ?? 0} opciones</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{votacion.total_votos || 0}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex items-center justify-end gap-1">
                      {votacion.status === 'draft' && (
                        <button onClick={() => { setSelectedVotacion(votacion); setActionDialog({ open: true, type: 'publish' }); }} className="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg" title="Publicar">
                          <Play size={18} />
                        </button>
                      )}
                      {votacion.status === 'active' && (
                        <button onClick={() => { setSelectedVotacion(votacion); setActionDialog({ open: true, type: 'close' }); }} className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg" title="Cerrar">
                          <Square size={18} />
                        </button>
                      )}
                      {(votacion.status === 'active' || votacion.status === 'closed') && (
                        <button onClick={() => openResultsModal(votacion)} className="p-2 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg" title="Resultados">
                          <BarChart3 size={18} />
                        </button>
                      )}
                      {votacion.status === 'draft' && (
                        <button onClick={() => openEditModal(votacion)} className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg"><Pencil size={18} /></button>
                      )}
                      {votacion.status === 'active' && (
                        <button onClick={() => { setSelectedVotacion(votacion); setActionDialog({ open: true, type: 'cancel' }); }} className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg" title="Cancelar">
                          <Ban size={18} />
                        </button>
                      )}
                      {votacion.status === 'draft' && (
                        <button onClick={() => { setSelectedVotacion(votacion); setShowDeleteDialog(true); }} className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg"><Trash2 size={18} /></button>
                      )}
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

      {/* Create/Edit Modal */}
      <Modal open={showModal} onClose={() => setShowModal(false)} title={selectedVotacion ? 'Editar Votacion' : 'Nueva Votacion'} size="lg">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Textarea label="Descripcion" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={3} />

          <div className="space-y-3">
            <label className="block text-sm font-medium text-gray-700">Opciones de voto</label>
            {formData.opciones.map((opcion, index) => (
              <div key={index} className="flex gap-2">
                <Input placeholder={`Opcion ${index + 1}`} value={opcion.label} onChange={(e) => updateOpcion(index, 'label', e.target.value)} className="flex-1" />
                {formData.opciones.length > 2 && (
                  <Button type="button" variant="ghost" onClick={() => removeOpcion(index)}>X</Button>
                )}
              </div>
            ))}
            <Button type="button" variant="secondary" size="sm" onClick={addOpcion}>+ Agregar opcion</Button>
          </div>

          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <input type="checkbox" id="requires_quorum" checked={formData.requires_quorum} onChange={(e) => setFormData({ ...formData, requires_quorum: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
              <label htmlFor="requires_quorum" className="text-sm text-gray-700">Requiere quorum</label>
            </div>
            {formData.requires_quorum && (
              <Input type="number" label="Porcentaje de quorum" value={formData.quorum_percentage} onChange={(e) => setFormData({ ...formData, quorum_percentage: parseInt(e.target.value) })} min={1} max={100} />
            )}
            <div className="flex items-center gap-2">
              <input type="checkbox" id="allow_abstention" checked={formData.allow_abstention} onChange={(e) => setFormData({ ...formData, allow_abstention: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
              <label htmlFor="allow_abstention" className="text-sm text-gray-700">Permitir abstencion</label>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">{selectedVotacion ? 'Guardar' : 'Crear'}</Button>
          </div>
        </form>
      </Modal>

      {/* Results Modal */}
      <Modal open={showResultsModal} onClose={() => setShowResultsModal(false)} title="Resultados de Votacion" size="lg">
        {resultados && (
          <div className="space-y-6">
            <div className="grid grid-cols-3 gap-4">
              <div className="p-4 bg-gray-50 rounded-lg text-center">
                <p className="text-2xl font-bold">{resultados.total_votos}</p>
                <p className="text-sm text-gray-500">Votos totales</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg text-center">
                <p className="text-2xl font-bold">{resultados.participacion.toFixed(1)}%</p>
                <p className="text-sm text-gray-500">Participacion</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg text-center">
                <p className="text-2xl font-bold">{resultados.quorum_alcanzado ? 'Si' : 'No'}</p>
                <p className="text-sm text-gray-500">Quorum alcanzado</p>
              </div>
            </div>

            <div className="space-y-3">
              {resultados.resultados.map((r) => (
                <div key={r.opcion_id} className="p-4 border rounded-lg">
                  <div className="flex justify-between mb-2">
                    <span className="font-medium">{r.label}</span>
                    <span className="text-gray-500">{r.count} votos ({r.percentage.toFixed(1)}%)</span>
                  </div>
                  <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div className="h-full bg-blue-600 rounded-full" style={{ width: `${r.percentage}%` }} />
                  </div>
                </div>
              ))}
              {resultados.total_abstenciones > 0 && (
                <div className="p-4 border rounded-lg bg-gray-50">
                  <div className="flex justify-between">
                    <span className="font-medium text-gray-500">Abstenciones</span>
                    <span className="text-gray-500">{resultados.total_abstenciones}</span>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </Modal>

      <ConfirmDialog open={actionDialog.open} onClose={() => setActionDialog({ open: false, type: null })} onConfirm={handleAction} title={actionDialog.type === 'publish' ? 'Publicar Votacion' : actionDialog.type === 'close' ? 'Cerrar Votacion' : 'Cancelar Votacion'} message={actionDialog.type === 'publish' ? 'La votacion sera visible y los vecinos podran votar.' : actionDialog.type === 'close' ? 'Se cerrara la votacion y no se podran emitir mas votos.' : 'La votacion sera cancelada.'} confirmText={actionDialog.type === 'publish' ? 'Publicar' : actionDialog.type === 'close' ? 'Cerrar' : 'Cancelar'} isLoading={isSaving} variant={actionDialog.type === 'cancel' ? 'danger' : 'warning'} />
      <ConfirmDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} onConfirm={handleDelete} title="Eliminar Votacion" message={`¿Eliminar "${selectedVotacion?.title}"?`} confirmText="Eliminar" isLoading={isSaving} />
    </div>
  )
}
