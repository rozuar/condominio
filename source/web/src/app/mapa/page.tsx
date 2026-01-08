'use client'

import { useState, useEffect } from 'react'
import { Map, MapPin, Layers, Info, ChevronDown, ChevronUp, X } from 'lucide-react'
import { getMapaData, MapaData } from '@/lib/api'
import MapWrapper from '@/components/mapa/MapWrapper'
import type { MapaArea, MapaPunto, AreaType } from '@/types'

const AREA_TYPE_LABELS: Record<AreaType, string> = {
  parcela: 'Parcelas',
  area_comun: 'Areas Comunes',
  acceso: 'Accesos',
  canal: 'Canales',
  camino: 'Caminos',
}

const AREA_TYPE_COLORS: Record<AreaType, string> = {
  parcela: 'bg-green-100 text-green-800',
  area_comun: 'bg-blue-100 text-blue-800',
  acceso: 'bg-yellow-100 text-yellow-800',
  canal: 'bg-cyan-100 text-cyan-800',
  camino: 'bg-gray-100 text-gray-800',
}

export default function MapaPage() {
  const [data, setData] = useState<MapaData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedArea, setSelectedArea] = useState<MapaArea | null>(null)
  const [selectedPunto, setSelectedPunto] = useState<MapaPunto | null>(null)
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set(['parcelas']))

  useEffect(() => {
    const fetchData = async () => {
      try {
        const result = await getMapaData()
        setData(result)
      } catch (err) {
        setError('Error al cargar los datos del mapa')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  const toggleSection = (section: string) => {
    setExpandedSections(prev => {
      const newSet = new Set(prev)
      if (newSet.has(section)) {
        newSet.delete(section)
      } else {
        newSet.add(section)
      }
      return newSet
    })
  }

  const groupAreasByType = (areas: MapaArea[]) => {
    const grouped: Record<AreaType, MapaArea[]> = {
      parcela: [],
      area_comun: [],
      acceso: [],
      canal: [],
      camino: [],
    }
    areas.forEach(area => {
      if (grouped[area.type]) {
        grouped[area.type].push(area)
      }
    })
    return grouped
  }

  if (loading) {
    return (
      <div className="py-8">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="animate-pulse">
            <div className="h-8 w-32 bg-gray-200 rounded mb-4" />
            <div className="h-[600px] bg-gray-200 rounded-lg" />
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="py-8">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center py-12 bg-red-50 rounded-lg">
            <p className="text-red-700">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  const groupedAreas = data ? groupAreasByType(data.areas) : null

  return (
    <div className="py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3 mb-2">
          <Map className="w-8 h-8 text-primary" />
          <h1 className="text-3xl font-bold text-gray-900">Mapa de la Comunidad</h1>
        </div>
        <p className="text-gray-600 mb-8">
          Explora las parcelas, areas comunes y puntos de interes de Vina Pelvin
        </p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Sidebar with areas and puntos */}
          <div className="lg:col-span-1 space-y-4">
            {/* Areas section */}
            <div className="bg-white rounded-lg shadow-md">
              <div className="p-4 border-b flex items-center gap-2">
                <Layers className="w-5 h-5 text-primary" />
                <h2 className="font-semibold">Areas del Mapa</h2>
              </div>
              <div className="divide-y">
                {groupedAreas && Object.entries(groupedAreas).map(([type, areas]) => {
                  if (areas.length === 0) return null
                  const isExpanded = expandedSections.has(type)

                  return (
                    <div key={type}>
                      <button
                        onClick={() => toggleSection(type)}
                        className="w-full p-3 flex items-center justify-between hover:bg-gray-50"
                      >
                        <div className="flex items-center gap-2">
                          <span className={`px-2 py-0.5 rounded text-xs font-medium ${AREA_TYPE_COLORS[type as AreaType]}`}>
                            {areas.length}
                          </span>
                          <span className="font-medium text-sm">
                            {AREA_TYPE_LABELS[type as AreaType]}
                          </span>
                        </div>
                        {isExpanded ? (
                          <ChevronUp className="w-4 h-4 text-gray-400" />
                        ) : (
                          <ChevronDown className="w-4 h-4 text-gray-400" />
                        )}
                      </button>
                      {isExpanded && (
                        <div className="bg-gray-50 px-3 py-2 max-h-48 overflow-y-auto">
                          {areas.map(area => (
                            <button
                              key={area.id}
                              onClick={() => {
                                setSelectedArea(area)
                                setSelectedPunto(null)
                              }}
                              className={`w-full text-left px-2 py-1.5 rounded text-sm hover:bg-gray-100 ${
                                selectedArea?.id === area.id ? 'bg-primary/10 text-primary' : ''
                              }`}
                            >
                              {area.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>
                  )
                })}
              </div>
            </div>

            {/* Puntos de interes section */}
            {data && data.puntos.length > 0 && (
              <div className="bg-white rounded-lg shadow-md">
                <div className="p-4 border-b flex items-center gap-2">
                  <MapPin className="w-5 h-5 text-primary" />
                  <h2 className="font-semibold">Puntos de Interes</h2>
                </div>
                <div className="p-3 space-y-1 max-h-64 overflow-y-auto">
                  {data.puntos.map(punto => (
                    <button
                      key={punto.id}
                      onClick={() => {
                        setSelectedPunto(punto)
                        setSelectedArea(null)
                      }}
                      className={`w-full text-left px-2 py-1.5 rounded text-sm hover:bg-gray-100 flex items-center gap-2 ${
                        selectedPunto?.id === punto.id ? 'bg-primary/10 text-primary' : ''
                      }`}
                    >
                      <MapPin className="w-4 h-4" />
                      {punto.name}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Map area */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-md overflow-hidden">
              {/* Leaflet Map */}
              <div className="aspect-[4/3] relative">
                <MapWrapper
                  areas={data?.areas || []}
                  puntos={data?.puntos || []}
                  selectedArea={selectedArea}
                  selectedPunto={selectedPunto}
                  onSelectArea={setSelectedArea}
                  onSelectPunto={setSelectedPunto}
                />

                {/* Overlay showing selected item info */}
                {(selectedArea || selectedPunto) && (
                  <div className="absolute top-4 right-4 bg-white rounded-lg shadow-lg p-4 max-w-xs z-[1000]">
                    <div className="flex items-start gap-2">
                      <Info className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                      <div className="flex-1">
                        {selectedArea && (
                          <>
                            <div className="flex items-start justify-between gap-2">
                              <h3 className="font-semibold text-gray-900">{selectedArea.name}</h3>
                              <button
                                onClick={() => setSelectedArea(null)}
                                className="text-gray-400 hover:text-gray-600"
                              >
                                <X className="w-4 h-4" />
                              </button>
                            </div>
                            <span className={`inline-block mt-1 px-2 py-0.5 rounded text-xs font-medium ${AREA_TYPE_COLORS[selectedArea.type]}`}>
                              {AREA_TYPE_LABELS[selectedArea.type]}
                            </span>
                            {selectedArea.description && (
                              <p className="text-sm text-gray-600 mt-2">{selectedArea.description}</p>
                            )}
                            {selectedArea.center_lat && selectedArea.center_lng && (
                              <p className="text-xs text-gray-400 mt-2">
                                Ubicacion: {selectedArea.center_lat.toFixed(6)}, {selectedArea.center_lng.toFixed(6)}
                              </p>
                            )}
                          </>
                        )}
                        {selectedPunto && (
                          <>
                            <div className="flex items-start justify-between gap-2">
                              <h3 className="font-semibold text-gray-900">{selectedPunto.name}</h3>
                              <button
                                onClick={() => setSelectedPunto(null)}
                                className="text-gray-400 hover:text-gray-600"
                              >
                                <X className="w-4 h-4" />
                              </button>
                            </div>
                            <span className="inline-block mt-1 px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                              {selectedPunto.type}
                            </span>
                            {selectedPunto.description && (
                              <p className="text-sm text-gray-600 mt-2">{selectedPunto.description}</p>
                            )}
                            <p className="text-xs text-gray-400 mt-2">
                              Ubicacion: {selectedPunto.lat.toFixed(6)}, {selectedPunto.lng.toFixed(6)}
                            </p>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                )}
              </div>

              {/* Map controls info */}
              <div className="p-4 bg-gray-50 border-t">
                <div className="flex flex-wrap gap-4 text-sm text-gray-600">
                  <div className="flex items-center gap-2">
                    <div className="w-4 h-4 rounded bg-green-500" />
                    <span>Parcelas ({groupedAreas?.parcela.length || 0})</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <div className="w-4 h-4 rounded bg-blue-500" />
                    <span>Areas Comunes ({groupedAreas?.area_comun.length || 0})</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <MapPin className="w-4 h-4 text-purple-500" />
                    <span>Puntos de Interes ({data?.puntos.length || 0})</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Statistics */}
            <div className="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div className="bg-white rounded-lg shadow p-4 text-center">
                <p className="text-2xl font-bold text-primary">{groupedAreas?.parcela.length || 0}</p>
                <p className="text-sm text-gray-600">Parcelas</p>
              </div>
              <div className="bg-white rounded-lg shadow p-4 text-center">
                <p className="text-2xl font-bold text-blue-600">{groupedAreas?.area_comun.length || 0}</p>
                <p className="text-sm text-gray-600">Areas Comunes</p>
              </div>
              <div className="bg-white rounded-lg shadow p-4 text-center">
                <p className="text-2xl font-bold text-yellow-600">{(groupedAreas?.acceso.length || 0) + (groupedAreas?.camino.length || 0)}</p>
                <p className="text-sm text-gray-600">Accesos y Caminos</p>
              </div>
              <div className="bg-white rounded-lg shadow p-4 text-center">
                <p className="text-2xl font-bold text-purple-600">{data?.puntos.length || 0}</p>
                <p className="text-sm text-gray-600">Puntos de Interes</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
