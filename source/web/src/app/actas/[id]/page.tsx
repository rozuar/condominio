'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Calendar, Users } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getActa } from '@/lib/api'
import {
  Acta,
  ACTA_TYPE_LABELS,
  ACTA_TYPE_COLORS,
} from '@/types'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-CL', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export default function ActaDetailPage() {
  const router = useRouter()
  const params = useParams()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()

  const [acta, setActa] = useState<Acta | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/auth/login')
    }
  }, [authLoading, isAuthenticated, router])

  useEffect(() => {
    if (!isAuthenticated || !params.id) return

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
        const data = await getActa(params.id as string, token)
        setActa(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar el acta')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, params.id])

  if (authLoading || (!isAuthenticated && !error)) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando...</div>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando acta...</div>
      </div>
    )
  }

  if (error || !acta) {
    return (
      <div className="py-8">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error || 'Acta no encontrada'}
          </div>
          <Link
            href="/actas"
            className="inline-flex items-center gap-2 mt-4 text-primary hover:underline"
          >
            <ArrowLeft className="h-4 w-4" />
            Volver a actas
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <Link
          href="/actas"
          className="inline-flex items-center gap-2 text-gray-600 hover:text-primary mb-6"
        >
          <ArrowLeft className="h-4 w-4" />
          Volver a actas
        </Link>

        <article className="bg-white rounded-lg shadow border overflow-hidden">
          <div className="p-6 sm:p-8">
            <div className="flex items-center gap-3 mb-4">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${ACTA_TYPE_COLORS[acta.type]}`}>
                {ACTA_TYPE_LABELS[acta.type]}
              </span>
            </div>

            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">
              {acta.title}
            </h1>

            <div className="flex flex-wrap gap-4 text-gray-500 mb-8 pb-6 border-b">
              <span className="flex items-center gap-2">
                <Calendar className="h-5 w-5" />
                {formatDate(acta.meeting_date)}
              </span>
              <span className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                {acta.attendees_count} asistentes
              </span>
            </div>

            <div className="prose max-w-none">
              {acta.content.split('\n').map((paragraph, idx) => (
                <p key={idx} className="mb-4 text-gray-700 leading-relaxed">
                  {paragraph}
                </p>
              ))}
            </div>
          </div>
        </article>
      </div>
    </div>
  )
}
