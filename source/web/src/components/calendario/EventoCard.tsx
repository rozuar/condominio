import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { Calendar, MapPin } from 'lucide-react'
import { Evento, EVENTO_TYPE_LABELS, EVENTO_TYPE_COLORS } from '@/types'

interface EventoCardProps {
  evento: Evento
}

export default function EventoCard({ evento }: EventoCardProps) {
  const eventDate = new Date(evento.event_date)
  const day = format(eventDate, 'd')
  const month = format(eventDate, 'MMM', { locale: es })
  const time = format(eventDate, 'HH:mm')

  return (
    <article className="bg-white border rounded-lg p-4 hover:shadow-md transition-shadow flex gap-4">
      <div className="flex-shrink-0 w-16 text-center">
        <div className="bg-primary text-white rounded-t px-2 py-1 text-xs uppercase">
          {month}
        </div>
        <div className="border border-t-0 rounded-b px-2 py-2">
          <span className="text-2xl font-bold text-gray-900">{day}</span>
        </div>
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className={`text-xs font-medium px-2 py-0.5 rounded ${EVENTO_TYPE_COLORS[evento.type]}`}>
            {EVENTO_TYPE_LABELS[evento.type]}
          </span>
        </div>

        <h3 className="font-semibold text-gray-900 truncate">{evento.title}</h3>

        <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
          <span className="flex items-center gap-1">
            <Calendar size={14} />
            {time} hrs
          </span>
          {evento.location && (
            <span className="flex items-center gap-1 truncate">
              <MapPin size={14} />
              {evento.location}
            </span>
          )}
        </div>
      </div>
    </article>
  )
}
