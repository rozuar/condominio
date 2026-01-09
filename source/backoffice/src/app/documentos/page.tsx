'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getDocumentos, createDocumento } from '@/lib/api'
import type { Documento, DocumentoCategory } from '@/types'
import { DOCUMENTO_CATEGORY_LABELS, DOCUMENTO_CATEGORY_COLORS } from '@/types'
import { Plus, Loader2, FileText, ExternalLink } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Textarea from '@/components/ui/Textarea'

export default function DocumentosPage() {
  const { getToken } = useAuth()
  const [documentos, setDocumentos] = useState<Documento[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [formData, setFormData] = useState({ title: '', description: '', file_url: '', category: 'otro' as DocumentoCategory, is_public: true })

  const fetchDocumentos = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getDocumentos(token, { page, per_page: 10 })
      setDocumentos(data.documentos)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => { fetchDocumentos() }, [page])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await createDocumento(token, formData)
      setShowModal(false)
      setFormData({ title: '', description: '', file_url: '', category: 'otro', is_public: true })
      fetchDocumentos()
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
          <h1 className="text-2xl font-bold text-gray-900">Documentos</h1>
          <p className="text-gray-500 mt-1">{total} documentos</p>
        </div>
        <Button onClick={() => setShowModal(true)} icon={<Plus size={20} />}>Nuevo Documento</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
        ) : documentos.length === 0 ? (
          <div className="text-center py-12 text-gray-500">No hay documentos</div>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Documento</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Categoria</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Visibilidad</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {documentos.map((doc) => (
                <tr key={doc.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <FileText className="h-5 w-5 text-gray-400" />
                      <div>
                        <div className="text-sm font-medium text-gray-900">{doc.title}</div>
                        {doc.description && <div className="text-sm text-gray-500 truncate max-w-md">{doc.description}</div>}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap"><Badge color={DOCUMENTO_CATEGORY_COLORS[doc.category]}>{DOCUMENTO_CATEGORY_LABELS[doc.category]}</Badge></td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={doc.is_public ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}>{doc.is_public ? 'Publico' : 'Privado'}</Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{format(new Date(doc.created_at), "d MMM yyyy", { locale: es })}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    {/^https?:\/\//.test(doc.file_url) ? (
                      <a
                        href={doc.file_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800"
                      >
                        <ExternalLink size={16} /> Ver
                      </a>
                    ) : (
                      <span
                        className="inline-flex items-center gap-1 text-gray-400 cursor-not-allowed"
                        title="El seed trae URLs relativas (ej: /documentos/...). Sube un link público (https://...) para evitar 404."
                      >
                        <ExternalLink size={16} /> Sin link
                      </span>
                    )}
                  </td>
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

      <Modal open={showModal} onClose={() => setShowModal(false)} title="Nuevo Documento">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required />
          <Select label="Categoria" value={formData.category} onChange={(e) => setFormData({ ...formData, category: e.target.value as DocumentoCategory })} options={[{ value: 'reglamento', label: 'Reglamento' }, { value: 'protocolo', label: 'Protocolo' }, { value: 'formulario', label: 'Formulario' }, { value: 'otro', label: 'Otro' }]} />
          <Input label="URL del Archivo" value={formData.file_url} onChange={(e) => setFormData({ ...formData, file_url: e.target.value })} required placeholder="https://..." />
          <Textarea label="Descripcion" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={3} />
          <div className="flex items-center gap-2">
            <input type="checkbox" id="is_public" checked={formData.is_public} onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })} className="h-4 w-4 text-blue-600 rounded" />
            <label htmlFor="is_public" className="text-sm text-gray-700">Documento publico</label>
          </div>
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">Crear</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
