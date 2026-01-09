'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getActas, createActa } from '@/lib/api'
import type { Acta, ActaType } from '@/types'
import { ACTA_TYPE_LABELS, ACTA_TYPE_COLORS } from '@/types'
import { Plus, Loader2, FileText } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'

export default function ActasPage() {
  const { getToken } = useAuth()
  const [actas, setActas] = useState<Acta[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [formData, setFormData] = useState({ title: '', content: '', meeting_date: '', type: 'ordinaria' as ActaType, attendees_count: 0 })

  const fetchActas = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getActas(token, { page, per_page: 10 })
      setActas(data.actas)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => { fetchActas() }, [page])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await createActa(token, formData)
      setShowModal(false)
      setFormData({ title: '', content: '', meeting_date: '', type: 'ordinaria', attendees_count: 0 })
      fetchActas()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al crear')
    } finally {
      setIsSaving(false)
    }
  }

  const totalPages = Math.ceil(total / 10)

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Actas</h1>
          <p className="text-gray-500 mt-1">{total} actas</p>
        </div>
        <Button onClick={() => setShowModal(true)} icon={<Plus size={20} />}>Nueva Acta</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
        ) : actas.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay actas</div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Acta</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha Reunion</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Asistentes</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {actas.map((acta) => (
                <tr key={acta.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <FileText className="h-5 w-5 text-gray-400" />
                      <div>
                        <div className="text-sm font-medium text-gray-900">{acta.title}</div>
                        <div className="text-sm text-gray-500 truncate max-w-md">{acta.content.substring(0, 80)}...</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap"><Badge color={ACTA_TYPE_COLORS[acta.type]}>{ACTA_TYPE_LABELS[acta.type]}</Badge></td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{format(new Date(acta.meeting_date), "d MMM yyyy", { locale: es })}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{acta.attendees_count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        {totalPages > 1 && (
          <div className="px-6 py-4 border-t flex items-center justify-between">
            <p className="text-sm text-gray-500">Pagina {page} de {totalPages}</p>
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" onClick={() => setPage(page - 1)} disabled={page === 1}>Anterior</Button>
              <Button variant="secondary" size="sm" onClick={() => setPage(page + 1)} disabled={page === totalPages}>Siguiente</Button>
            </div>
          </div>
        )}
      </div>

      <Modal open={showModal} onClose={() => setShowModal(false)} title="Nueva Acta" size="lg">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Select label="Tipo" value={formData.type} onChange={(e) => setFormData({ ...formData, type: e.target.value as ActaType })} options={[{ value: 'ordinaria', label: 'Ordinaria' }, { value: 'extraordinaria', label: 'Extraordinaria' }]} />
          <Input type="date" label="Fecha de Reunion" value={formData.meeting_date} onChange={(e) => setFormData({ ...formData, meeting_date: e.target.value })} required />
          <Input type="number" label="Numero de Asistentes" value={formData.attendees_count} onChange={(e) => setFormData({ ...formData, attendees_count: parseInt(e.target.value) })} required />
          <Textarea label="Contenido del Acta" value={formData.content} onChange={(e) => setFormData({ ...formData, content: e.target.value })} rows={8} required />
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">Crear Acta</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
