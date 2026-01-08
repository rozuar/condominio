'use client'

import dynamic from 'next/dynamic'
import { Loader2 } from 'lucide-react'
import type { MapaArea, MapaPunto } from '@/types'

// Dynamically import the map component with SSR disabled
const LeafletMap = dynamic(() => import('./LeafletMap'), {
  ssr: false,
  loading: () => (
    <div className="w-full h-full flex items-center justify-center bg-gray-100">
      <div className="text-center">
        <Loader2 className="w-8 h-8 text-primary mx-auto animate-spin" />
        <p className="text-gray-600 mt-2">Cargando mapa...</p>
      </div>
    </div>
  ),
})

interface MapWrapperProps {
  areas: MapaArea[]
  puntos: MapaPunto[]
  selectedArea: MapaArea | null
  selectedPunto: MapaPunto | null
  onSelectArea: (area: MapaArea | null) => void
  onSelectPunto: (punto: MapaPunto | null) => void
}

export default function MapWrapper(props: MapWrapperProps) {
  return <LeafletMap {...props} />
}
