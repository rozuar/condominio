'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import { format, formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import {
  ArrowLeft,
  Vote,
  Users,
  CheckCircle,
  Clock,
  Calendar,
  AlertCircle,
  Ban,
} from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getVotacion, getVotacionResultados, emitirVoto } from '@/lib/api'
import {
  Votacion,
  VotacionResultado,
  VOTACION_STATUS_LABELS,
  VOTACION_STATUS_COLORS,
} from '@/types'

export default function VotacionDetailPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()

  const [votacion, setVotacion] = useState<Votacion | null>(null)
  const [resultado, setResultado] = useState<VotacionResultado | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedOpcion, setSelectedOpcion] = useState<string | null>(null)
  const [isAbstention, setIsAbstention] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitSuccess, setSubmitSuccess] = useState(false)

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/auth/login')
    }
  }, [authLoading, isAuthenticated, router])

  useEffect(() => {
    if (!isAuthenticated || !id) return

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
        const [votacionData, resultadoData] = await Promise.all([
          getVotacion(id, token),
          getVotacionResultados(id, token),
        ])

        setVotacion(votacionData)
        setResultado(resultadoData)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar la votacion')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, id])

  const handleVote = async () => {
    if (!votacion) return
    if (!isAbstention && !selectedOpcion) {
      setSubmitError('Selecciona una opcion o abstente')
      return
    }

    const token = getToken()
    if (!token) {
      setSubmitError('No hay sesion activa')
      return
    }

    setIsSubmitting(true)
    setSubmitError(null)

    try {
      await emitirVoto(votacion.id, token, {
        opcion_id: isAbstention ? undefined : selectedOpcion || undefined,
        is_abstention: isAbstention,
      })

      setSubmitSuccess(true)

      // Reload data
      const [votacionData, resultadoData] = await Promise.all([
        getVotacion(id, token),
        getVotacionResultados(id, token),
      ])
      setVotacion(votacionData)
      setResultado(resultadoData)
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'Error al emitir el voto')
    } finally {
      setIsSubmitting(false)
    }
  }

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
        <div className="text-gray-500">Cargando votacion...</div>
      </div>
    )
  }

  if (error || !votacion) {
    return (
      <div className="py-8">
        <div className="max-w-4xl mx-auto px-4">
          <Link
            href="/votaciones"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-primary mb-6"
          >
            <ArrowLeft className="w-4 h-4" />
            Volver a votaciones
          </Link>
          <div className="p-6 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error || 'Votacion no encontrada'}
          </div>
        </div>
      </div>
    )
  }

  const isActive = votacion.status === 'active'
  const canVote = isActive && !votacion.has_voted
  const showResults = votacion.has_voted || !isActive

  return (
    <div className="py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Back link */}
        <Link
          href="/votaciones"
          className="inline-flex items-center gap-2 text-gray-600 hover:text-primary mb-6"
        >
          <ArrowLeft className="w-4 h-4" />
          Volver a votaciones
        </Link>

        {/* Main card */}
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          {/* Header */}
          <div className="p-6 border-b bg-gray-50">
            <div className="flex flex-wrap items-center gap-2 mb-4">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${VOTACION_STATUS_COLORS[votacion.status]}`}>
                {VOTACION_STATUS_LABELS[votacion.status]}
              </span>
              {votacion.has_voted && (
                <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                  <CheckCircle className="w-4 h-4" />
                  Ya votaste
                </span>
              )}
            </div>

            <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mb-2">
              {votacion.title}
            </h1>

            {votacion.description && (
              <p className="text-gray-600">{votacion.description}</p>
            )}
          </div>

          {/* Content */}
          <div className="p-6">
            {/* Meta info */}
            <div className="flex flex-wrap gap-6 text-sm text-gray-600 mb-8">
              {votacion.opciones && (
                <div className="flex items-center gap-2">
                  <Vote className="w-4 h-4" />
                  <span>{votacion.opciones.length} opciones</span>
                </div>
              )}
              {resultado && (
                <div className="flex items-center gap-2">
                  <Users className="w-4 h-4" />
                  <span>{resultado.total_votos} votos de {resultado.total_vecinos} vecinos</span>
                </div>
              )}
              {votacion.start_date && (
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  <span>Inicio: {format(new Date(votacion.start_date), "d MMM yyyy", { locale: es })}</span>
                </div>
              )}
              {votacion.end_date && (
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  <span>Cierre: {format(new Date(votacion.end_date), "d MMM yyyy", { locale: es })}</span>
                </div>
              )}
            </div>

            {/* Quorum info */}
            {votacion.requires_quorum && resultado && (
              <div className={`mb-6 p-4 rounded-lg ${resultado.quorum_alcanzado ? 'bg-green-50 border border-green-200' : 'bg-yellow-50 border border-yellow-200'}`}>
                <div className="flex items-center gap-2">
                  {resultado.quorum_alcanzado ? (
                    <CheckCircle className="w-5 h-5 text-green-600" />
                  ) : (
                    <AlertCircle className="w-5 h-5 text-yellow-600" />
                  )}
                  <span className={resultado.quorum_alcanzado ? 'text-green-800' : 'text-yellow-800'}>
                    Quorum: {resultado.participacion.toFixed(1)}% de participacion
                    ({votacion.quorum_percentage}% requerido)
                  </span>
                </div>
              </div>
            )}

            {/* Voting form */}
            {canVote && (
              <div className="mb-8">
                <h2 className="text-lg font-semibold mb-4">Emite tu voto</h2>

                {submitError && (
                  <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
                    {submitError}
                  </div>
                )}

                {submitSuccess && (
                  <div className="mb-4 p-3 bg-green-50 border border-green-200 text-green-700 rounded-lg text-sm">
                    Tu voto ha sido registrado correctamente
                  </div>
                )}

                <div className="space-y-3">
                  {votacion.opciones?.map((opcion) => (
                    <label
                      key={opcion.id}
                      className={`flex items-center p-4 border rounded-lg cursor-pointer transition-colors ${
                        selectedOpcion === opcion.id && !isAbstention
                          ? 'border-primary bg-primary/5'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <input
                        type="radio"
                        name="opcion"
                        value={opcion.id}
                        checked={selectedOpcion === opcion.id && !isAbstention}
                        onChange={() => {
                          setSelectedOpcion(opcion.id)
                          setIsAbstention(false)
                        }}
                        className="w-4 h-4 text-primary focus:ring-primary"
                      />
                      <span className="ml-3 font-medium">{opcion.label}</span>
                      {opcion.description && (
                        <span className="ml-2 text-sm text-gray-500">
                          - {opcion.description}
                        </span>
                      )}
                    </label>
                  ))}

                  {votacion.allow_abstention && (
                    <label
                      className={`flex items-center p-4 border rounded-lg cursor-pointer transition-colors ${
                        isAbstention
                          ? 'border-gray-500 bg-gray-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <input
                        type="radio"
                        name="opcion"
                        checked={isAbstention}
                        onChange={() => {
                          setIsAbstention(true)
                          setSelectedOpcion(null)
                        }}
                        className="w-4 h-4 text-gray-500 focus:ring-gray-500"
                      />
                      <Ban className="w-4 h-4 ml-3 text-gray-500" />
                      <span className="ml-2 font-medium text-gray-600">Abstenerse</span>
                    </label>
                  )}
                </div>

                <button
                  onClick={handleVote}
                  disabled={isSubmitting || (!selectedOpcion && !isAbstention)}
                  className="mt-6 w-full sm:w-auto px-6 py-3 bg-primary text-white font-medium rounded-lg hover:bg-primary-light disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {isSubmitting ? 'Enviando...' : 'Confirmar voto'}
                </button>
              </div>
            )}

            {/* Results */}
            {showResults && resultado && (
              <div>
                <h2 className="text-lg font-semibold mb-4">Resultados</h2>

                <div className="space-y-4">
                  {resultado.resultados.map((r, index) => (
                    <div key={r.opcion_id} className="relative">
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-medium">{r.label}</span>
                        <span className="text-sm text-gray-600">
                          {r.count} voto{r.count !== 1 ? 's' : ''} ({r.percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <div className="h-8 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full transition-all duration-500 ${
                            index === 0 ? 'bg-primary' : 'bg-primary/60'
                          }`}
                          style={{ width: `${r.percentage}%` }}
                        />
                      </div>
                    </div>
                  ))}

                  {resultado.total_abstenciones > 0 && (
                    <div className="pt-4 border-t">
                      <div className="flex items-center justify-between text-gray-600">
                        <span>Abstenciones</span>
                        <span>{resultado.total_abstenciones}</span>
                      </div>
                    </div>
                  )}
                </div>

                <div className="mt-6 p-4 bg-gray-50 rounded-lg">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-gray-500">Total votos:</span>
                      <span className="ml-2 font-medium">{resultado.total_votos}</span>
                    </div>
                    <div>
                      <span className="text-gray-500">Participacion:</span>
                      <span className="ml-2 font-medium">{resultado.participacion.toFixed(1)}%</span>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="px-6 py-4 border-t bg-gray-50 text-sm text-gray-500">
            <div className="flex items-center gap-2">
              <Clock className="w-4 h-4" />
              <span>
                Creada {formatDistanceToNow(new Date(votacion.created_at), {
                  addSuffix: true,
                  locale: es,
                })}
                {votacion.creator_name && ` por ${votacion.creator_name}`}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
