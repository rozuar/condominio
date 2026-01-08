import Link from 'next/link'
import { AlertTriangle, X } from 'lucide-react'
import { Emergencia } from '@/types'

interface EmergenciaBannerProps {
  emergencia: Emergencia
}

export default function EmergenciaBanner({ emergencia }: EmergenciaBannerProps) {
  return (
    <div className="bg-red-600 text-white">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-3">
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <AlertTriangle className="w-6 h-6 flex-shrink-0 animate-pulse" />
            <div>
              <p className="font-bold text-sm uppercase tracking-wide">Alerta Critica</p>
              <p className="font-medium">{emergencia.title}</p>
            </div>
          </div>
          <Link
            href={`/emergencias/${emergencia.id}`}
            className="flex-shrink-0 bg-white text-red-600 px-4 py-2 rounded font-medium text-sm hover:bg-red-50 transition-colors"
          >
            Ver detalles
          </Link>
        </div>
      </div>
    </div>
  )
}
