'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getMapaAreas, getMapaPuntos } from '@/lib/api'
import type { MapaArea, MapaPunto } from '@/types'
import { Loader2, MapPin, Square } from 'lucide-react'

export default function MapaPage() {
  const { getToken } = useAuth()
  const [areas, setAreas] = useState<MapaArea[]>([])
  const [puntos, setPuntos] = useState<MapaPunto[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchData = async () => {
      const token = getToken()
      if (!token) return

      setIsLoading(true)
      try {
        const [areasData, puntosData] = await Promise.all([
          getMapaAreas(token),
          getMapaPuntos(token),
        ])
        setAreas(areasData.areas)
        setPuntos(puntosData.puntos)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar')
      } finally {
        setIsLoading(false)
      }
    }
    fetchData()
  }, [])

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Mapa</h1>
        <p className="text-gray-500 mt-1">Gestion de areas y puntos de interes</p>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      {isLoading ? (
        <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
            <div className="px-6 py-4 border-b bg-gray-50 flex items-center gap-2">
              <Square className="h-5 w-5 text-gray-500" />
              <h2 className="font-medium">Areas ({areas.length})</h2>
            </div>
            <div className="divide-y divide-gray-100 max-h-96 overflow-y-auto">
              {areas.length === 0 ? (
                <div className="p-8 text-center text-gray-500">No hay areas</div>
              ) : (
                areas.map((area) => (
                  <div key={area.id} className="p-4 hover:bg-gray-50">
                    <div className="flex items-center gap-3">
                      <div className="w-4 h-4 rounded" style={{ backgroundColor: area.fill_color }} />
                      <div>
                        <p className="font-medium text-gray-900">{area.name}</p>
                        <p className="text-sm text-gray-500">{area.type}</p>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
            <div className="px-6 py-4 border-b bg-gray-50 flex items-center gap-2">
              <MapPin className="h-5 w-5 text-gray-500" />
              <h2 className="font-medium">Puntos de Interes ({puntos.length})</h2>
            </div>
            <div className="divide-y divide-gray-100 max-h-96 overflow-y-auto">
              {puntos.length === 0 ? (
                <div className="p-8 text-center text-gray-500">No hay puntos</div>
              ) : (
                puntos.map((punto) => (
                  <div key={punto.id} className="p-4 hover:bg-gray-50">
                    <div className="flex items-center gap-3">
                      <MapPin className="h-5 w-5 text-blue-500" />
                      <div>
                        <p className="font-medium text-gray-900">{punto.name}</p>
                        <p className="text-sm text-gray-500">{punto.type}</p>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
