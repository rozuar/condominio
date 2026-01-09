'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getComunicados, createComunicado, updateComunicado, deleteComunicado } from '@/lib/api'
import type { Comunicado, ComunicadoType } from '@/types'
import { COMUNICADO_TYPE_LABELS, COMUNICADO_TYPE_COLORS } from '@/types'
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-react'
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
  { value: 'informativo', label: 'Informativo' },
  { value: 'seguridad', label: 'Seguridad' },
  { value: 'tesoreria', label: 'Tesoreria' },
  { value: 'asamblea', label: 'Asamblea' },
]

export default function ComunicadosPage() {
  const { getToken } = useAuth()
  const [comunicados, setComunicados] = useState<Comunicado[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  // Modal states
  const [showModal, setShowModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [selectedComunicado, setSelectedComunicado] = useState<Comunicado | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  // Form states
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    type: 'informativo' as ComunicadoType,
    is_public: true,
  })

  const fetchComunicados = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getComunicados(token, { page, per_page: 10 })
      setComunicados(data.comunicados)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar comunicados')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchComunicados()
  }, [page])

  const openCreateModal = () => {
    setSelectedComunicado(null)
    setFormData({ title: '', content: '', type: 'informativo', is_public: true })
    setShowModal(true)
  }

  const openEditModal = (comunicado: Comunicado) => {
    setSelectedComunicado(comunicado)
    setFormData({
      title: comunicado.title,
      content: comunicado.content,
      type: comunicado.type,
      is_public: comunicado.is_public,
    })
    setShowModal(true)
  }

  const openDeleteDialog = (comunicado: Comunicado) => {
    setSelectedComunicado(comunicado)
    setShowDeleteDialog(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      if (selectedComunicado) {
        await updateComunicado(selectedComunicado.id, token, formData)
      } else {
        await createComunicado(token, formData)
      }
      setShowModal(false)
      fetchComunicados()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedComunicado) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteComunicado(selectedComunicado.id, token)
      setShowDeleteDialog(false)
      setSelectedComunicado(null)
      fetchComunicados()
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
          <h1 className="text-2xl font-bold text-gray-900">Comunicados</h1>
          <p className="text-gray-500 mt-1">{total} comunicados en total</p>
        </div>
        <Button onClick={openCreateModal} icon={<Plus size={20} />}>
          Nuevo Comunicado
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
        ) : comunicados.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            No hay comunicados
          </div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Titulo
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Tipo
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Fecha
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Estado
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Acciones
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {comunicados.map((comunicado) => (
                <tr key={comunicado.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="text-sm font-medium text-gray-900">{comunicado.title}</div>
                    <div className="text-sm text-gray-500 truncate max-w-md">
                      {comunicado.content.substring(0, 100)}...
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={COMUNICADO_TYPE_COLORS[comunicado.type]}>
                      {COMUNICADO_TYPE_LABELS[comunicado.type]}
                    </Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {format(new Date(comunicado.created_at), "d MMM yyyy", { locale: es })}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={comunicado.is_public ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}>
                      {comunicado.is_public ? 'Publico' : 'Privado'}
                    </Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => openEditModal(comunicado)}
                        className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                      >
                        <Pencil size={18} />
                      </button>
                      <button
                        onClick={() => openDeleteDialog(comunicado)}
                        className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                      >
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
            <p className="text-sm text-gray-500">
              Pagina {page} de {totalPages}
            </p>
            <div className="flex gap-2">
              <Button
                variant="secondary"
                size="sm"
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
              >
                Anterior
              </Button>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => setPage(page + 1)}
                disabled={page === totalPages}
              >
                Siguiente
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      <Modal
        open={showModal}
        onClose={() => setShowModal(false)}
        title={selectedComunicado ? 'Editar Comunicado' : 'Nuevo Comunicado'}
        size="lg"
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input
            label="Titulo"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            required
            placeholder="Titulo del comunicado"
          />

          <Select
            label="Tipo"
            value={formData.type}
            onChange={(e) => setFormData({ ...formData, type: e.target.value as ComunicadoType })}
            options={typeOptions}
          />

          <Textarea
            label="Contenido"
            value={formData.content}
            onChange={(e) => setFormData({ ...formData, content: e.target.value })}
            required
            rows={6}
            placeholder="Escribe el contenido del comunicado..."
          />

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="is_public"
              checked={formData.is_public}
              onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })}
              className="h-4 w-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
            />
            <label htmlFor="is_public" className="text-sm text-gray-700">
              Comunicado publico (visible para todos)
            </label>
          </div>

          <div className="flex gap-3 pt-4">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setShowModal(false)}
              className="flex-1"
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              loading={isSaving}
              className="flex-1"
            >
              {selectedComunicado ? 'Guardar Cambios' : 'Crear Comunicado'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation */}
      <ConfirmDialog
        open={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        onConfirm={handleDelete}
        title="Eliminar Comunicado"
        message={`¿Estas seguro de eliminar "${selectedComunicado?.title}"? Esta accion no se puede deshacer.`}
        confirmText="Eliminar"
        isLoading={isSaving}
      />
    </div>
  )
}
