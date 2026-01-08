import Link from 'next/link'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { Comunicado, COMUNICADO_TYPE_LABELS, COMUNICADO_TYPE_COLORS } from '@/types'

interface ComunicadoCardProps {
  comunicado: Comunicado
}

export default function ComunicadoCard({ comunicado }: ComunicadoCardProps) {
  const publishedDate = comunicado.published_at
    ? format(new Date(comunicado.published_at), "d 'de' MMMM, yyyy", { locale: es })
    : null

  const excerpt = comunicado.content.length > 150
    ? comunicado.content.substring(0, 150) + '...'
    : comunicado.content

  return (
    <article className="bg-white border rounded-lg p-6 hover:shadow-md transition-shadow">
      <div className="flex items-center gap-2 mb-3">
        <span className={`text-xs font-medium px-2 py-1 rounded ${COMUNICADO_TYPE_COLORS[comunicado.type]}`}>
          {COMUNICADO_TYPE_LABELS[comunicado.type]}
        </span>
        {publishedDate && (
          <span className="text-sm text-gray-500">{publishedDate}</span>
        )}
      </div>

      <h3 className="text-lg font-semibold text-gray-900 mb-2">
        <Link href={`/comunicados/${comunicado.id}`} className="hover:text-primary">
          {comunicado.title}
        </Link>
      </h3>

      <p className="text-gray-600 text-sm mb-4">{excerpt}</p>

      <Link
        href={`/comunicados/${comunicado.id}`}
        className="text-primary font-medium text-sm hover:underline"
      >
        Leer más →
      </Link>
    </article>
  )
}
