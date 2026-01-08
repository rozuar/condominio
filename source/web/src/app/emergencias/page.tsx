import { getEmergencias, getActiveEmergencias } from '@/lib/api'
import EmergenciaCard from '@/components/emergencias/EmergenciaCard'
import EmergenciaBanner from '@/components/emergencias/EmergenciaBanner'
import { EMERGENCIA_PRIORITY_LABELS, EMERGENCIA_STATUS_LABELS, EmergenciaPriority, EmergenciaStatus } from '@/types'

export const revalidate = 30 // Revalidate every 30 seconds for emergencies

interface PageProps {
  searchParams: { status?: string; priority?: string; page?: string }
}

async function getData(status?: string, priority?: string, page = 1) {
  try {
    return await getEmergencias({ status, priority, page, per_page: 12 })
  } catch (error) {
    return { emergencias: [], total: 0, page: 1, per_page: 12 }
  }
}

async function getActive() {
  try {
    return await getActiveEmergencias()
  } catch (error) {
    return { emergencias: [], total: 0 }
  }
}

export default async function EmergenciasPage({ searchParams }: PageProps) {
  const { status, priority, page } = searchParams
  const [data, activeData] = await Promise.all([
    getData(status, priority, page ? parseInt(page) : 1),
    getActive()
  ])

  const activeEmergencias = activeData.emergencias || []
  const criticalEmergencias = activeEmergencias.filter(e => e.priority === 'critical')

  return (
    <div className="py-8">
      {/* Critical emergencies banner */}
      {criticalEmergencias.length > 0 && (
        <div className="mb-8">
          {criticalEmergencias.map((emergencia) => (
            <EmergenciaBanner key={emergencia.id} emergencia={emergencia} />
          ))}
        </div>
      )}

      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3 mb-2">
          <h1 className="text-3xl font-bold text-gray-900">Emergencias</h1>
          {activeEmergencias.length > 0 && (
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-800">
              {activeEmergencias.length} activa{activeEmergencias.length !== 1 ? 's' : ''}
            </span>
          )}
        </div>
        <p className="text-gray-600 mb-8">Avisos urgentes y alertas de la comunidad</p>

        {/* Filters */}
        <div className="flex flex-wrap gap-4 mb-8">
          {/* Status filter */}
          <div className="flex flex-wrap gap-2">
            <span className="text-sm text-gray-500 self-center mr-1">Estado:</span>
            <a
              href="/emergencias"
              className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                !status ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Todos
            </a>
            {(Object.keys(EMERGENCIA_STATUS_LABELS) as EmergenciaStatus[]).map((s) => (
              <a
                key={s}
                href={`/emergencias?status=${s}${priority ? `&priority=${priority}` : ''}`}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  status === s ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {EMERGENCIA_STATUS_LABELS[s]}
              </a>
            ))}
          </div>

          {/* Priority filter */}
          <div className="flex flex-wrap gap-2">
            <span className="text-sm text-gray-500 self-center mr-1">Prioridad:</span>
            <a
              href={`/emergencias${status ? `?status=${status}` : ''}`}
              className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                !priority ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Todas
            </a>
            {(Object.keys(EMERGENCIA_PRIORITY_LABELS) as EmergenciaPriority[]).map((p) => (
              <a
                key={p}
                href={`/emergencias?${status ? `status=${status}&` : ''}priority=${p}`}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  priority === p ? 'bg-primary text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {EMERGENCIA_PRIORITY_LABELS[p]}
              </a>
            ))}
          </div>
        </div>

        {/* List */}
        {data.emergencias.length > 0 ? (
          <div className="space-y-4">
            {data.emergencias.map((emergencia) => (
              <EmergenciaCard key={emergencia.id} emergencia={emergencia} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-green-50 rounded-lg">
            <div className="text-4xl mb-4">&#x2705;</div>
            <p className="text-green-800 font-medium">No hay emergencias activas</p>
            <p className="text-green-600 text-sm mt-1">La comunidad se encuentra sin alertas</p>
          </div>
        )}

        {/* Pagination */}
        {data.total > data.per_page && (
          <div className="flex justify-center gap-2 mt-8">
            {Array.from({ length: Math.ceil(data.total / data.per_page) }, (_, i) => i + 1).map((p) => (
              <a
                key={p}
                href={`/emergencias?${status ? `status=${status}&` : ''}${priority ? `priority=${priority}&` : ''}page=${p}`}
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
