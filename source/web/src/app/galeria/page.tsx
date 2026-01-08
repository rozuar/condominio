import { getGalerias } from '@/lib/api'
import { Image, Calendar, ImageIcon } from 'lucide-react'
import Link from 'next/link'

export const revalidate = 60 // Revalidate every minute

interface PageProps {
  searchParams: { page?: string }
}

async function getData(page = 1) {
  try {
    return await getGalerias({ page, per_page: 12, is_public: true })
  } catch (error) {
    return { galerias: [], total: 0, page: 1, per_page: 12 }
  }
}

function formatDate(dateStr: string | undefined) {
  if (!dateStr) return null
  return new Date(dateStr).toLocaleDateString('es-CL', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export default async function GaleriaPage({ searchParams }: PageProps) {
  const { page } = searchParams
  const data = await getData(page ? parseInt(page) : 1)

  return (
    <div className="py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3 mb-2">
          <Image className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-bold text-gray-900">Galeria</h1>
        </div>
        <p className="text-gray-600 mb-8">Fotos y videos de eventos de la comunidad</p>

        {/* Gallery Grid */}
        {data.galerias.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {data.galerias.map((galeria) => (
              <Link
                key={galeria.id}
                href={`/galeria/${galeria.id}`}
                className="group bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
              >
                {/* Cover Image */}
                <div className="aspect-video bg-gray-100 relative">
                  {galeria.cover_image_url ? (
                    <img
                      src={galeria.cover_image_url}
                      alt={galeria.title}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <ImageIcon className="w-16 h-16 text-gray-300" />
                    </div>
                  )}
                  {/* Items count badge */}
                  {galeria.items_count !== undefined && galeria.items_count > 0 && (
                    <div className="absolute bottom-2 right-2 bg-black/70 text-white px-2 py-1 rounded text-sm">
                      {galeria.items_count} {galeria.items_count === 1 ? 'foto' : 'fotos'}
                    </div>
                  )}
                </div>

                {/* Info */}
                <div className="p-4">
                  <h2 className="font-semibold text-gray-900 group-hover:text-primary transition-colors line-clamp-1">
                    {galeria.title}
                  </h2>
                  {galeria.description && (
                    <p className="text-sm text-gray-600 mt-1 line-clamp-2">
                      {galeria.description}
                    </p>
                  )}
                  {galeria.event_date && (
                    <div className="flex items-center gap-1.5 text-sm text-gray-500 mt-2">
                      <Calendar className="w-4 h-4" />
                      <span>{formatDate(galeria.event_date)}</span>
                    </div>
                  )}
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-gray-50 rounded-lg">
            <ImageIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-600 font-medium">No hay galerias disponibles</p>
            <p className="text-gray-500 text-sm mt-1">Pronto se agregaran fotos de eventos</p>
          </div>
        )}

        {/* Pagination */}
        {data.total > data.per_page && (
          <div className="flex justify-center gap-2 mt-8">
            {Array.from({ length: Math.ceil(data.total / data.per_page) }, (_, i) => i + 1).map((p) => (
              <a
                key={p}
                href={`/galeria?page=${p}`}
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
