'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { FileText, File, FileCheck, FileQuestion, Eye } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getDocumentos } from '@/lib/api'
import {
  Documento,
  DocumentoCategory,
  DOCUMENTO_CATEGORY_LABELS,
  DOCUMENTO_CATEGORY_COLORS,
} from '@/types'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-CL', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const CATEGORY_ICONS: Record<DocumentoCategory, typeof FileText> = {
  reglamento: FileCheck,
  protocolo: File,
  formulario: FileText,
  otro: FileQuestion,
}

export default function DocumentosPage() {
  const router = useRouter()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()

  const [documentos, setDocumentos] = useState<Documento[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedCategory, setSelectedCategory] = useState<DocumentoCategory | ''>('')

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/auth/login')
    }
  }, [authLoading, isAuthenticated, router])

  useEffect(() => {
    if (!isAuthenticated) return

    const fetchData = async () => {
      setIsLoading(true)
      setError(null)

      const token = getToken()
      if (!token) {
        setError('No hay sesión activa')
        setIsLoading(false)
        return
      }

      try {
        const data = await getDocumentos(token, {
          category: selectedCategory || undefined,
          per_page: 50,
        })
        setDocumentos(data.documentos || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar los documentos')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, selectedCategory])

  if (authLoading || (!isAuthenticated && !error)) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando...</div>
      </div>
    )
  }

  const categories: { value: DocumentoCategory | ''; label: string }[] = [
    { value: '', label: 'Todos' },
    { value: 'reglamento', label: 'Reglamentos' },
    { value: 'protocolo', label: 'Protocolos' },
    { value: 'formulario', label: 'Formularios' },
    { value: 'otro', label: 'Otros' },
  ]

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Documentos</h1>
          <p className="text-gray-600 mt-2">Documentos oficiales de la comunidad</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {/* Category Filter */}
        <div className="mb-6">
          <div className="flex flex-wrap gap-2">
            {categories.map((cat) => (
              <button
                key={cat.value}
                onClick={() => setSelectedCategory(cat.value)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                  selectedCategory === cat.value
                    ? 'bg-primary text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {cat.label}
              </button>
            ))}
          </div>
        </div>

        {/* Documents List */}
        {isLoading ? (
          <div className="text-center py-12 text-gray-500">Cargando documentos...</div>
        ) : documentos.length === 0 ? (
          <div className="text-center py-12">
            <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No hay documentos disponibles</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {documentos.map((doc) => {
              const Icon = CATEGORY_ICONS[doc.category]
              return (
                <div
                  key={doc.id}
                  className="bg-white rounded-lg shadow border overflow-hidden hover:shadow-md transition-shadow"
                >
                  <div className="p-6">
                    <div className="flex items-start gap-4">
                      <div className={`p-3 rounded-lg ${DOCUMENTO_CATEGORY_COLORS[doc.category].replace('text-', 'bg-').split(' ')[0]}`}>
                        <Icon className={`h-6 w-6 ${DOCUMENTO_CATEGORY_COLORS[doc.category].split(' ')[1]}`} />
                      </div>
                      <div className="flex-1 min-w-0">
                        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium mb-2 ${DOCUMENTO_CATEGORY_COLORS[doc.category]}`}>
                          {DOCUMENTO_CATEGORY_LABELS[doc.category]}
                        </span>
                        <h3 className="font-semibold text-gray-900 truncate">
                          {doc.title}
                        </h3>
                        {doc.description && (
                          <p className="text-sm text-gray-500 mt-1 line-clamp-2">
                            {doc.description}
                          </p>
                        )}
                        <p className="text-xs text-gray-400 mt-2">
                          Agregado el {formatDate(doc.created_at)}
                        </p>
                      </div>
                    </div>

                    <div className="mt-4 pt-4 border-t">
                      <button
                        onClick={() => {
                          if (doc.file_url) {
                            window.open(doc.file_url, '_blank', 'noopener,noreferrer')
                          }
                        }}
                        disabled={!doc.file_url}
                        className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-primary/10 text-primary rounded-lg font-medium hover:bg-primary/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <Eye className="h-4 w-4" />
                        Ver documento
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
