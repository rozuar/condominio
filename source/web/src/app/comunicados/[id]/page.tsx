import { notFound } from 'next/navigation'
import Link from 'next/link'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { ArrowLeft } from 'lucide-react'
import { getComunicado } from '@/lib/api'
import { Comunicado, COMUNICADO_TYPE_LABELS, COMUNICADO_TYPE_COLORS } from '@/types'

interface PageProps {
  params: { id: string }
}

async function getData(id: string): Promise<Comunicado | null> {
  try {
    return await getComunicado(id) as Comunicado
  } catch (error) {
    return null
  }
}

export default async function ComunicadoDetailPage({ params }: PageProps) {
  const comunicado = await getData(params.id)

  if (!comunicado) {
    notFound()
  }

  const publishedDate = comunicado.published_at
    ? format(new Date(comunicado.published_at), "d 'de' MMMM 'de' yyyy, HH:mm", { locale: es })
    : null

  return (
    <div className="py-8">
      <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
        <Link
          href="/comunicados"
          className="inline-flex items-center gap-1 text-gray-600 hover:text-primary mb-6"
        >
          <ArrowLeft size={16} />
          Volver a comunicados
        </Link>

        <article>
          <header className="mb-8">
            <div className="flex items-center gap-3 mb-4">
              <span className={`text-sm font-medium px-3 py-1 rounded ${COMUNICADO_TYPE_COLORS[comunicado.type]}`}>
                {COMUNICADO_TYPE_LABELS[comunicado.type]}
              </span>
            </div>

            <h1 className="text-3xl font-bold text-gray-900 mb-4">
              {comunicado.title}
            </h1>

            <div className="flex items-center gap-4 text-sm text-gray-500">
              {publishedDate && <span>Publicado: {publishedDate}</span>}
              {comunicado.author_name && <span>Por: {comunicado.author_name}</span>}
            </div>
          </header>

          <div className="prose max-w-none">
            {comunicado.content.split('\n').map((paragraph: string, index: number) => (
              <p key={index} className="mb-4 text-gray-700 leading-relaxed">
                {paragraph}
              </p>
            ))}
          </div>
        </article>
      </div>
    </div>
  )
}
