'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { FileText, Calendar, Users, ChevronRight } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getActas } from '@/lib/api'
import {
  Acta,
  ActaType,
  ACTA_TYPE_LABELS,
  ACTA_TYPE_COLORS,
} from '@/types'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-CL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export default function ActasPage() {
  const router = useRouter()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()

  const [actas, setActas] = useState<Acta[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedType, setSelectedType] = useState<ActaType | ''>('')

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
        const data = await getActas(token, {
          type: selectedType || undefined,
          per_page: 50,
        })
        setActas(data.actas || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar las actas')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, selectedType])

  if (authLoading || (!isAuthenticated && !error)) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando...</div>
      </div>
    )
  }

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Actas de Reuniones</h1>
          <p className="text-gray-600 mt-2">Registro oficial de las reuniones de la comunidad</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {/* Filter */}
        <div className="mb-6">
          <div className="flex gap-2">
            <button
              onClick={() => setSelectedType('')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                selectedType === ''
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Todas
            </button>
            <button
              onClick={() => setSelectedType('ordinaria')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                selectedType === 'ordinaria'
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Ordinarias
            </button>
            <button
              onClick={() => setSelectedType('extraordinaria')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                selectedType === 'extraordinaria'
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Extraordinarias
            </button>
          </div>
        </div>

        {/* Actas List */}
        {isLoading ? (
          <div className="text-center py-12 text-gray-500">Cargando actas...</div>
        ) : actas.length === 0 ? (
          <div className="text-center py-12">
            <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No hay actas disponibles</p>
          </div>
        ) : (
          <div className="space-y-4">
            {actas.map((acta) => (
              <Link
                key={acta.id}
                href={`/actas/${acta.id}`}
                className="block bg-white rounded-lg shadow border hover:shadow-md transition-shadow"
              >
                <div className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${ACTA_TYPE_COLORS[acta.type]}`}>
                          {ACTA_TYPE_LABELS[acta.type]}
                        </span>
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900 mb-2">
                        {acta.title}
                      </h3>
                      <div className="flex flex-wrap gap-4 text-sm text-gray-500">
                        <span className="flex items-center gap-1">
                          <Calendar className="h-4 w-4" />
                          {formatDate(acta.meeting_date)}
                        </span>
                        <span className="flex items-center gap-1">
                          <Users className="h-4 w-4" />
                          {acta.attendees_count} asistentes
                        </span>
                      </div>
                    </div>
                    <ChevronRight className="h-5 w-5 text-gray-400 flex-shrink-0" />
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
