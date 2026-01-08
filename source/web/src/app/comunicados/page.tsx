import { getComunicados } from '@/lib/api'
import ComunicadoCard from '@/components/comunicados/ComunicadoCard'
import { COMUNICADO_TYPE_LABELS, ComunicadoType } from '@/types'

export const revalidate = 60

interface PageProps {
  searchParams: { type?: string; page?: string }
}

async function getData(type?: string, page = 1) {
  try {
    return await getComunicados({ type, page, per_page: 12 })
  } catch (error) {
    return { comunicados: [], total: 0, page: 1, per_page: 12 }
  }
}

export default async function ComunicadosPage({ searchParams }: PageProps) {
  const { type, page } = searchParams
  const data = await getData(type, page ? parseInt(page) : 1)

  return (
    <div className="py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Comunicados</h1>
        <p className="text-gray-600 mb-8">Mantente informado sobre las novedades de la comunidad</p>

        {/* Filters */}
        <div className="flex flex-wrap gap-2 mb-8">
          <a
            href="/comunicados"
            className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
              !type ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            Todos
          </a>
          {(Object.keys(COMUNICADO_TYPE_LABELS) as ComunicadoType[]).map((t) => (
            <a
              key={t}
              href={`/comunicados?type=${t}`}
              className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                type === t ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {COMUNICADO_TYPE_LABELS[t]}
            </a>
          ))}
        </div>

        {/* List */}
        {data.comunicados.length > 0 ? (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {data.comunicados.map((comunicado) => (
              <ComunicadoCard key={comunicado.id} comunicado={comunicado} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-gray-500">No hay comunicados disponibles</p>
          </div>
        )}

        {/* Pagination */}
        {data.total > data.per_page && (
          <div className="flex justify-center gap-2 mt-8">
            {Array.from({ length: Math.ceil(data.total / data.per_page) }, (_, i) => i + 1).map((p) => (
              <a
                key={p}
                href={`/comunicados?${type ? `type=${type}&` : ''}page=${p}`}
                className={`w-10 h-10 flex items-center justify-center rounded ${
                  p === data.page ? 'bg-primary text-white' : 'bg-gray-100 hover:bg-gray-200'
                }`}
              >
                {p}
              </a>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
