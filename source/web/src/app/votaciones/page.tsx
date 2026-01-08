'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Vote, Plus } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getVotaciones, getActiveVotaciones } from '@/lib/api'
import VotacionCard from '@/components/votaciones/VotacionCard'
import {
  Votacion,
  VotacionStatus,
  VOTACION_STATUS_LABELS,
  VOTACION_STATUS_COLORS,
} from '@/types'

export default function VotacionesPage() {
  const router = useRouter()
  const { isAuthenticated, isLoading: authLoading, getToken, user } = useAuth()

  const [votaciones, setVotaciones] = useState<Votacion[]>([])
  const [activeVotaciones, setActiveVotaciones] = useState<Votacion[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedStatus, setSelectedStatus] = useState<VotacionStatus | ''>('')
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const perPage = 10

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
        setError('No hay sesion activa')
        setIsLoading(false)
        return
      }

      try {
        const [votacionesData, activeData] = await Promise.all([
          getVotaciones(token, {
            status: selectedStatus || undefined,
            page,
            per_page: perPage,
          }),
          getActiveVotaciones(token),
        ])

        setVotaciones(votacionesData.votaciones || [])
        setTotal(votacionesData.total)
        setActiveVotaciones(activeData.votaciones || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar las votaciones')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, selectedStatus, page])

  const handleStatusChange = (status: VotacionStatus | '') => {
    setSelectedStatus(status)
    setPage(1)
  }

  if (authLoading || (!isAuthenticated && !error)) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando...</div>
      </div>
    )
  }

  const pendingVotes = activeVotaciones.filter(v => !v.has_voted)
  const totalPages = Math.ceil(total / perPage)

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold text-gray-900">Votaciones</h1>
              {pendingVotes.length > 0 && (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                  {pendingVotes.length} pendiente{pendingVotes.length !== 1 ? 's' : ''}
                </span>
              )}
            </div>
            <p className="text-gray-600 mt-2">Participa en las decisiones de la comunidad</p>
          </div>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {/* Active Votations Alert */}
        {pendingVotes.length > 0 && (
          <div className="mb-8 p-4 bg-green-50 border border-green-200 rounded-lg">
            <div className="flex items-center gap-3">
              <Vote className="w-6 h-6 text-green-600" />
              <div>
                <p className="font-medium text-green-800">
                  Tienes {pendingVotes.length} votacion{pendingVotes.length !== 1 ? 'es' : ''} activa{pendingVotes.length !== 1 ? 's' : ''} pendiente{pendingVotes.length !== 1 ? 's' : ''}
                </p>
                <p className="text-sm text-green-700">
                  Tu voto es importante para las decisiones de la comunidad
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Filters */}
        <div className="flex flex-wrap gap-2 mb-6">
          <button
            onClick={() => handleStatusChange('')}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              selectedStatus === ''
                ? 'bg-primary text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            Todas
          </button>
          {(['active', 'closed', 'cancelled'] as VotacionStatus[]).map((status) => (
            <button
              key={status}
              onClick={() => handleStatusChange(status)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedStatus === status
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {VOTACION_STATUS_LABELS[status]}
            </button>
          ))}
        </div>

        {/* List */}
        {isLoading ? (
          <div className="text-center py-12 text-gray-500">Cargando votaciones...</div>
        ) : votaciones.length === 0 ? (
          <div className="text-center py-12 bg-gray-50 rounded-lg">
            <Vote className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-600 font-medium">No hay votaciones disponibles</p>
            <p className="text-gray-500 text-sm mt-1">
              {selectedStatus
                ? `No hay votaciones con estado "${VOTACION_STATUS_LABELS[selectedStatus]}"`
                : 'Las votaciones apareceran aqui cuando la directiva las publique'}
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {votaciones.map((votacion) => (
              <VotacionCard key={votacion.id} votacion={votacion} />
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center gap-2 mt-8">
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
              <button
                key={p}
                onClick={() => setPage(p)}
                className={`w-10 h-10 flex items-center justify-center rounded ${
                  p === page
                    ? 'bg-primary text-white'
                    : 'bg-gray-100 hover:bg-gray-200'
                }`}
              >
                {p}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
