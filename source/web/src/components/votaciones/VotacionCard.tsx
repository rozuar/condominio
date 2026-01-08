'use client'

import Link from 'next/link'
import { formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import { Vote, Users, CheckCircle, Clock, Calendar } from 'lucide-react'
import {
  Votacion,
  VotacionStatus,
  VOTACION_STATUS_LABELS,
  VOTACION_STATUS_COLORS,
} from '@/types'

interface VotacionCardProps {
  votacion: Votacion
}

export default function VotacionCard({ votacion }: VotacionCardProps) {
  const isActive = votacion.status === 'active'
  const isClosed = votacion.status === 'closed'
  const hasVoted = votacion.has_voted

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('es-CL', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  }

  return (
    <Link
      href={`/votaciones/${votacion.id}`}
      className={`block bg-white rounded-lg shadow border hover:shadow-md transition-shadow ${
        isActive ? 'border-l-4 border-l-green-500' : ''
      }`}
    >
      <div className="p-6">
        <div className="flex flex-wrap items-center gap-2 mb-3">
          <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${VOTACION_STATUS_COLORS[votacion.status]}`}>
            {VOTACION_STATUS_LABELS[votacion.status]}
          </span>
          {hasVoted && (
            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              <CheckCircle className="w-3 h-3" />
              Ya votaste
            </span>
          )}
          {votacion.requires_quorum && (
            <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
              Requiere quorum ({votacion.quorum_percentage}%)
            </span>
          )}
        </div>

        <h3 className="text-lg font-semibold text-gray-900 mb-2">
          {votacion.title}
        </h3>

        {votacion.description && (
          <p className="text-gray-600 text-sm mb-4 line-clamp-2">
            {votacion.description}
          </p>
        )}

        <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500">
          {votacion.opciones && votacion.opciones.length > 0 && (
            <div className="flex items-center gap-1">
              <Vote className="w-4 h-4" />
              <span>{votacion.opciones.length} opciones</span>
            </div>
          )}

          {typeof votacion.total_votos === 'number' && (
            <div className="flex items-center gap-1">
              <Users className="w-4 h-4" />
              <span>{votacion.total_votos} voto{votacion.total_votos !== 1 ? 's' : ''}</span>
            </div>
          )}

          {votacion.end_date && (
            <div className="flex items-center gap-1">
              <Calendar className="w-4 h-4" />
              <span>
                {isActive ? 'Cierra' : 'Cerro'} {formatDate(votacion.end_date)}
              </span>
            </div>
          )}

          <div className="flex items-center gap-1 ml-auto">
            <Clock className="w-4 h-4" />
            <span>
              {formatDistanceToNow(new Date(votacion.created_at), {
                addSuffix: true,
                locale: es,
              })}
            </span>
          </div>
        </div>
      </div>
    </Link>
  )
}
