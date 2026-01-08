import Link from 'next/link'
import { format, formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import { AlertTriangle, AlertCircle, Info, CheckCircle, Clock } from 'lucide-react'
import {
  Emergencia,
  EMERGENCIA_PRIORITY_LABELS,
  EMERGENCIA_PRIORITY_COLORS,
  EMERGENCIA_STATUS_LABELS,
  EmergenciaPriority,
} from '@/types'

interface EmergenciaCardProps {
  emergencia: Emergencia
}

const priorityIcons: Record<EmergenciaPriority, React.ReactNode> = {
  critical: <AlertTriangle className="w-5 h-5" />,
  high: <AlertCircle className="w-5 h-5" />,
  medium: <Info className="w-5 h-5" />,
  low: <Info className="w-5 h-5" />,
}

const priorityBorders: Record<EmergenciaPriority, string> = {
  critical: 'border-l-4 border-l-red-500',
  high: 'border-l-4 border-l-orange-500',
  medium: 'border-l-4 border-l-yellow-500',
  low: 'border-l-4 border-l-blue-500',
}

export default function EmergenciaCard({ emergencia }: EmergenciaCardProps) {
  const createdDate = format(new Date(emergencia.created_at), "d 'de' MMMM, yyyy HH:mm", { locale: es })
  const timeAgo = formatDistanceToNow(new Date(emergencia.created_at), { addSuffix: true, locale: es })

  const isActive = emergencia.status === 'active'
  const isResolved = emergencia.status === 'resolved'

  return (
    <article
      className={`bg-white border rounded-lg p-6 ${priorityBorders[emergencia.priority]} ${
        isActive ? 'shadow-md' : 'opacity-75'
      }`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          {/* Header */}
          <div className="flex flex-wrap items-center gap-2 mb-3">
            <span className={`inline-flex items-center gap-1 text-xs font-medium px-2 py-1 rounded ${EMERGENCIA_PRIORITY_COLORS[emergencia.priority]}`}>
              {priorityIcons[emergencia.priority]}
              {EMERGENCIA_PRIORITY_LABELS[emergencia.priority]}
            </span>
            <span className={`text-xs font-medium px-2 py-1 rounded ${
              isActive ? 'bg-green-100 text-green-800' :
              isResolved ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
            }`}>
              {EMERGENCIA_STATUS_LABELS[emergencia.status]}
            </span>
            <span className="text-sm text-gray-500 flex items-center gap-1">
              <Clock className="w-3 h-3" />
              {timeAgo}
            </span>
          </div>

          {/* Title */}
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            <Link href={`/emergencias/${emergencia.id}`} className="hover:text-primary">
              {emergencia.title}
            </Link>
          </h3>

          {/* Content */}
          <p className="text-gray-600 text-sm mb-4 whitespace-pre-line">
            {emergencia.content.length > 300
              ? emergencia.content.substring(0, 300) + '...'
              : emergencia.content}
          </p>

          {/* Footer */}
          <div className="flex flex-wrap items-center gap-4 text-xs text-gray-500">
            <span>Publicado: {createdDate}</span>
            {emergencia.creator_name && (
              <span>Por: {emergencia.creator_name}</span>
            )}
            {emergencia.expires_at && (
              <span className="text-orange-600">
                Expira: {format(new Date(emergencia.expires_at), "d MMM yyyy HH:mm", { locale: es })}
              </span>
            )}
            {isResolved && emergencia.resolved_at && (
              <span className="text-green-600 flex items-center gap-1">
                <CheckCircle className="w-3 h-3" />
                Resuelto: {format(new Date(emergencia.resolved_at), "d MMM yyyy HH:mm", { locale: es })}
              </span>
            )}
          </div>
        </div>
      </div>
    </article>
  )
}
