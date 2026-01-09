'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getGalerias, createGaleria, updateGaleria, deleteGaleria } from '@/lib/api'
import type { Galeria } from '@/types'
import { Plus, Pencil, Trash2, Loader2, Image } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import Input from '@/components/ui/Input'
import Textarea from '@/components/ui/Textarea'

export default function GaleriasPage() {
  const { getToken } = useAuth()
  const [galerias, setGalerias] = useState<Galeria[]>([])
  const [total, setTotal] = useState(0)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [selectedGaleria, setSelectedGaleria] = useState<Galeria | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  const [formData, setFormData] = useState({ title: '', description: '', is_public: true })

  const fetchGalerias = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getGalerias(token, { per_page: 50 })
      setGalerias(data.galerias)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => { fetchGalerias() }, [])

  const openCreateModal = () => {
    setSelectedGaleria(null)
    setFormData({ title: '', description: '', is_public: true })
    setShowModal(true)
  }

  const openEditModal = (galeria: Galeria) => {
    setSelectedGaleria(galeria)
    setFormData({ title: galeria.title, description: galeria.description || '', is_public: galeria.is_public })
    setShowModal(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      if (selectedGaleria) {
        await updateGaleria(selectedGaleria.id, token, formData)
      } else {
        await createGaleria(token, formData)
      }
      setShowModal(false)
      fetchGalerias()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedGaleria) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await deleteGaleria(selectedGaleria.id, token)
      setShowDeleteDialog(false)
      fetchGalerias()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar')
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Galerias</h1>
          <p className="text-gray-500 mt-1">{total} galerias</p>
        </div>
        <Button onClick={openCreateModal} icon={<Plus size={20} />}>Nueva Galeria</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      {isLoading ? (
        <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
      ) : galerias.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border text-gray-500">No hay galerias</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {galerias.map((galeria) => (
            <div key={galeria.id} className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <div className="aspect-video bg-gray-100 flex items-center justify-center">
                {galeria.cover_image_url ? (
                  <img src={galeria.cover_image_url} alt={galeria.title} className="w-full h-full object-cover" />
                ) : (
                  <Image className="h-12 w-12 text-gray-300" />
                )}
              </div>
              <div className="p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="font-medium text-gray-900">{galeria.title}</h3>
                    <p className="text-sm text-gray-500 mt-1">{galeria.items_count || 0} elementos</p>
                  </div>
                  <Badge color={galeria.is_public ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}>
                    {galeria.is_public ? 'Publica' : 'Privada'}
                  </Badge>
                </div>
                <div className="flex gap-2 mt-4">
                  <Button variant="secondary" size="sm" onClick={() => openEditModal(galeria)} icon={<Pencil size={16} />}>Editar</Button>
                  <Button variant="ghost" size="sm" onClick={() => { setSelectedGaleria(galeria); setShowDeleteDialog(true); }} icon={<Trash2 size={16} />}>Eliminar</Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal open={showModal} onClose={() => setShowModal(false)} title={selectedGaleria ? 'Editar Galeria' : 'Nueva Galeria'}>
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Textarea label="Descripcion" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={3} />
          <div className="flex items-center gap-2">
            <input type="checkbox" id="is_public" checked={formData.is_public} onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
            <label htmlFor="is_public" className="text-sm text-gray-700">Galeria publica</label>
          </div>
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">{selectedGaleria ? 'Guardar' : 'Crear'}</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} onConfirm={handleDelete} title="Eliminar Galeria" message={`¿Eliminar "${selectedGaleria?.title}"?`} confirmText="Eliminar" isLoading={isSaving} />
    </div>
  )
}
