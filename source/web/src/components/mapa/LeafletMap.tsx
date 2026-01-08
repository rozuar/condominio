'use client'

import { useEffect, useRef } from 'react'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import type { MapaArea, MapaPunto, AreaType } from '@/types'

// Fix for default marker icons in webpack
delete (L.Icon.Default.prototype as unknown as { _getIconUrl?: unknown })._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
})

interface LeafletMapProps {
  areas: MapaArea[]
  puntos: MapaPunto[]
  selectedArea: MapaArea | null
  selectedPunto: MapaPunto | null
  onSelectArea: (area: MapaArea | null) => void
  onSelectPunto: (punto: MapaPunto | null) => void
}

// Default center for Vina Pelvin (approximate location near Melipilla, Chile)
const DEFAULT_CENTER: L.LatLngExpression = [-33.68, -71.21]
const DEFAULT_ZOOM = 15

// Colors for area types
const AREA_COLORS: Record<AreaType, { fill: string; stroke: string }> = {
  parcela: { fill: '#22c55e', stroke: '#16a34a' },
  area_comun: { fill: '#3b82f6', stroke: '#2563eb' },
  acceso: { fill: '#eab308', stroke: '#ca8a04' },
  canal: { fill: '#06b6d4', stroke: '#0891b2' },
  camino: { fill: '#6b7280', stroke: '#4b5563' },
}

// Custom marker icons for different punto types
const createPuntoIcon = (type: string): L.DivIcon => {
  const colors: Record<string, string> = {
    porteria: '#ef4444',
    sede: '#8b5cf6',
    agua: '#06b6d4',
    electricidad: '#f59e0b',
    emergencia: '#dc2626',
    default: '#6366f1',
  }

  const color = colors[type.toLowerCase()] || colors.default

  return L.divIcon({
    className: 'custom-marker',
    html: `
      <div style="
        background-color: ${color};
        width: 24px;
        height: 24px;
        border-radius: 50%;
        border: 3px solid white;
        box-shadow: 0 2px 4px rgba(0,0,0,0.3);
        display: flex;
        align-items: center;
        justify-content: center;
      ">
        <div style="
          width: 8px;
          height: 8px;
          background-color: white;
          border-radius: 50%;
        "></div>
      </div>
    `,
    iconSize: [24, 24],
    iconAnchor: [12, 12],
    popupAnchor: [0, -12],
  })
}

export default function LeafletMap({
  areas,
  puntos,
  selectedArea,
  selectedPunto,
  onSelectArea,
  onSelectPunto,
}: LeafletMapProps) {
  const mapRef = useRef<L.Map | null>(null)
  const mapContainerRef = useRef<HTMLDivElement>(null)
  const areasLayerRef = useRef<L.FeatureGroup | null>(null)
  const puntosLayerRef = useRef<L.LayerGroup | null>(null)
  const polygonRefs = useRef<Map<string, L.Polygon>>(new Map())
  const markerRefs = useRef<Map<string, L.Marker>>(new Map())

  // Initialize map
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return

    // Create map
    const map = L.map(mapContainerRef.current, {
      center: DEFAULT_CENTER,
      zoom: DEFAULT_ZOOM,
      zoomControl: true,
    })

    // Add tile layer (OpenStreetMap)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
      maxZoom: 19,
    }).addTo(map)

    // Create layer groups
    areasLayerRef.current = L.featureGroup().addTo(map)
    puntosLayerRef.current = L.layerGroup().addTo(map)

    mapRef.current = map

    return () => {
      map.remove()
      mapRef.current = null
    }
  }, [])

  // Update areas on map
  useEffect(() => {
    if (!mapRef.current || !areasLayerRef.current) return

    // Clear existing areas
    areasLayerRef.current.clearLayers()
    polygonRefs.current.clear()

    // Add areas
    areas.forEach((area) => {
      const colors = AREA_COLORS[area.type] || AREA_COLORS.parcela

      // GeoJSON coordinates are [lng, lat], Leaflet needs [lat, lng]
      const latLngs: L.LatLngExpression[][] = area.coordinates.map((ring) =>
        ring.map(([lng, lat]) => [lat, lng] as L.LatLngExpression)
      )

      const polygon = L.polygon(latLngs, {
        fillColor: area.fill_color || colors.fill,
        fillOpacity: 0.4,
        color: area.stroke_color || colors.stroke,
        weight: 2,
      })

      polygon.on('click', () => {
        onSelectArea(area)
        onSelectPunto(null)
      })

      polygon.bindTooltip(area.name, {
        permanent: false,
        direction: 'center',
      })

      polygon.addTo(areasLayerRef.current!)
      polygonRefs.current.set(area.id, polygon)
    })

    // Fit bounds if we have areas
    if (areas.length > 0 && areasLayerRef.current.getLayers().length > 0) {
      const bounds = areasLayerRef.current.getBounds()
      if (bounds.isValid()) {
        mapRef.current.fitBounds(bounds, { padding: [20, 20] })
      }
    }
  }, [areas, onSelectArea, onSelectPunto])

  // Update puntos on map
  useEffect(() => {
    if (!mapRef.current || !puntosLayerRef.current) return

    // Clear existing puntos
    puntosLayerRef.current.clearLayers()
    markerRefs.current.clear()

    // Add puntos
    puntos.forEach((punto) => {
      const marker = L.marker([punto.lat, punto.lng], {
        icon: createPuntoIcon(punto.type),
      })

      marker.on('click', () => {
        onSelectPunto(punto)
        onSelectArea(null)
      })

      marker.bindPopup(`
        <div style="min-width: 150px;">
          <strong>${punto.name}</strong>
          <br/>
          <span style="color: #666; font-size: 12px;">${punto.type}</span>
          ${punto.description ? `<br/><span style="font-size: 12px;">${punto.description}</span>` : ''}
        </div>
      `)

      marker.addTo(puntosLayerRef.current!)
      markerRefs.current.set(punto.id, marker)
    })
  }, [puntos, onSelectArea, onSelectPunto])

  // Highlight selected area
  useEffect(() => {
    polygonRefs.current.forEach((polygon, id) => {
      if (selectedArea?.id === id) {
        polygon.setStyle({
          weight: 4,
          fillOpacity: 0.6,
        })
        polygon.bringToFront()

        // Pan to area
        if (mapRef.current) {
          const bounds = polygon.getBounds()
          mapRef.current.panTo(bounds.getCenter())
        }
      } else {
        polygon.setStyle({
          weight: 2,
          fillOpacity: 0.4,
        })
      }
    })
  }, [selectedArea])

  // Highlight selected punto
  useEffect(() => {
    if (selectedPunto && mapRef.current) {
      const marker = markerRefs.current.get(selectedPunto.id)
      if (marker) {
        mapRef.current.panTo([selectedPunto.lat, selectedPunto.lng])
        marker.openPopup()
      }
    }
  }, [selectedPunto])

  return (
    <div
      ref={mapContainerRef}
      className="w-full h-full"
      style={{ minHeight: '400px' }}
    />
  )
}
