import { getEventos } from '@/lib/api'
import EventoCard from '@/components/calendario/EventoCard'
import { EVENTO_TYPE_LABELS, EventoType } from '@/types'

export const revalidate = 60

interface PageProps {
  searchParams: { type?: string }
}

async function getData(type?: string) {
  try {
    return await getEventos({ type, upcoming: true, per_page: 20 })
  } catch (error) {
    return { eventos: [], total: 0, page: 1, per_page: 20 }
  }
}

export default async function CalendarioPage({ searchParams }: PageProps) {
  const { type } = searchParams
  const data = await getData(type)

  return (
    <div className="py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Calendario</h1>
        <p className="text-gray-600 mb-8">Próximos eventos y actividades de la comunidad</p>

        {/* Filters */}
        <div className="flex flex-wrap gap-2 mb-8">
          <a
            href="/calendario"
            className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
              !type ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            Todos
          </a>
          {(Object.keys(EVENTO_TYPE_LABELS) as EventoType[]).map((t) => (
            <a
              key={t}
              href={`/calendario?type=${t}`}
              className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                type === t ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {EVENTO_TYPE_LABELS[t]}
            </a>
          ))}
        </div>

        {/* List */}
        {data.eventos.length > 0 ? (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {data.eventos.map((evento) => (
              <EventoCard key={evento.id} evento={evento} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-gray-500">No hay eventos próximos</p>
          </div>
        )}
      </div>
    </div>
  )
}
