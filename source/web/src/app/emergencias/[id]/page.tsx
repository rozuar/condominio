import { notFound } from 'next/navigation'
import Link from 'next/link'
import { format, formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import { ArrowLeft, AlertTriangle, AlertCircle, Info, CheckCircle, Clock, User } from 'lucide-react'
import { getEmergencia } from '@/lib/api'
import {
  EMERGENCIA_PRIORITY_LABELS,
  EMERGENCIA_PRIORITY_COLORS,
  EMERGENCIA_STATUS_LABELS,
  EmergenciaPriority,
  EmergenciaStatus,
} from '@/types'

export const revalidate = 30

interface PageProps {
  params: { id: string }
}

const priorityIcons: Record<EmergenciaPriority, React.ReactNode> = {
  critical: <AlertTriangle className="w-6 h-6" />,
  high: <AlertCircle className="w-6 h-6" />,
  medium: <Info className="w-6 h-6" />,
  low: <Info className="w-6 h-6" />,
}

async function getData(id: string) {
  try {
    return await getEmergencia(id)
  } catch (error) {
    return null
  }
}

export default async function EmergenciaDetailPage({ params }: PageProps) {
  const emergencia = await getData(params.id)

  if (!emergencia) {
    notFound()
  }

  const createdDate = format(new Date(emergencia.created_at), "EEEE d 'de' MMMM, yyyy 'a las' HH:mm", { locale: es })
  const timeAgo = formatDistanceToNow(new Date(emergencia.created_at), { addSuffix: true, locale: es })

  const status = emergencia.status as EmergenciaStatus
  const priority = emergencia.priority as EmergenciaPriority
  const isActive = status === 'active'
  const isResolved = status === 'resolved'
  const isCritical = priority === 'critical'

  return (
    <div className={`min-h-screen ${isCritical && isActive ? 'bg-red-50' : 'bg-gray-50'}`}>
      {/* Critical header */}
      {isCritical && isActive && (
        <div className="bg-red-600 text-white py-4">
          <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 flex items-center gap-3">
            <AlertTriangle className="w-8 h-8 animate-pulse" />
            <div>
              <p className="font-bold uppercase tracking-wide">Alerta Critica Activa</p>
              <p className="text-red-100 text-sm">Esta emergencia requiere atencion inmediata</p>
            </div>
          </div>
        </div>
      )}

      <div className="py-8">
        <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8">
          {/* Back link */}
          <Link
            href="/emergencias"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-primary mb-6"
          >
            <ArrowLeft className="w-4 h-4" />
            Volver a emergencias
          </Link>

          {/* Main card */}
          <article className="bg-white rounded-lg shadow-lg overflow-hidden">
            {/* Header */}
            <div className={`p-6 ${isCritical ? 'bg-red-100' : 'bg-gray-50'} border-b`}>
              <div className="flex flex-wrap items-center gap-3 mb-4">
                <span className={`inline-flex items-center gap-2 text-sm font-medium px-3 py-1 rounded-full ${EMERGENCIA_PRIORITY_COLORS[priority]}`}>
                  {priorityIcons[priority]}
                  Prioridad {EMERGENCIA_PRIORITY_LABELS[priority]}
                </span>
                <span className={`text-sm font-medium px-3 py-1 rounded-full ${
                  isActive ? 'bg-green-100 text-green-800' :
                  isResolved ? 'bg-blue-100 text-blue-800' : 'bg-gray-200 text-gray-800'
                }`}>
                  {EMERGENCIA_STATUS_LABELS[status]}
                </span>
              </div>

              <h1 className="text-2xl md:text-3xl font-bold text-gray-900">
                {emergencia.title}
              </h1>
            </div>

            {/* Content */}
            <div className="p-6">
              <div className="prose prose-gray max-w-none mb-8">
                <p className="whitespace-pre-line text-gray-700 leading-relaxed">
                  {emergencia.content}
                </p>
              </div>

              {/* Meta info */}
              <div className="border-t pt-6 space-y-3">
                <div className="flex items-center gap-2 text-sm text-gray-600">
                  <Clock className="w-4 h-4" />
                  <span>Publicado {timeAgo}</span>
                  <span className="text-gray-400">({createdDate})</span>
                </div>

                {emergencia.creator_name && (
                  <div className="flex items-center gap-2 text-sm text-gray-600">
                    <User className="w-4 h-4" />
                    <span>Creado por: {emergencia.creator_name}</span>
                  </div>
                )}

                {emergencia.expires_at && (
                  <div className="flex items-center gap-2 text-sm text-orange-600">
                    <Clock className="w-4 h-4" />
                    <span>
                      Expira: {format(new Date(emergencia.expires_at), "d 'de' MMMM, yyyy 'a las' HH:mm", { locale: es })}
                    </span>
                  </div>
                )}

                {isResolved && emergencia.resolved_at && (
                  <div className="flex items-center gap-2 text-sm text-green-600">
                    <CheckCircle className="w-4 h-4" />
                    <span>
                      Resuelto el {format(new Date(emergencia.resolved_at), "d 'de' MMMM, yyyy 'a las' HH:mm", { locale: es })}
                      {emergencia.resolver_name && ` por ${emergencia.resolver_name}`}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>
  )
}
