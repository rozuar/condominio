'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import {
  ArrowLeft,
  Calendar,
  X,
  ChevronLeft,
  ChevronRight,
  Play,
  ImageIcon,
} from 'lucide-react'
import { getGaleria, GaleriaWithItems } from '@/lib/api'
import type { GaleriaItem } from '@/types'

function formatDate(dateStr: string | undefined) {
  if (!dateStr) return null
  return new Date(dateStr).toLocaleDateString('es-CL', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export default function GaleriaDetailPage() {
  const params = useParams()
  const id = params.id as string

  const [galeria, setGaleria] = useState<GaleriaWithItems | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [lightboxOpen, setLightboxOpen] = useState(false)
  const [currentIndex, setCurrentIndex] = useState(0)

  useEffect(() => {
    const fetchGaleria = async () => {
      try {
        const data = await getGaleria(id)
        setGaleria(data)
      } catch (err) {
        setError('Error al cargar la galeria')
      } finally {
        setLoading(false)
      }
    }

    fetchGaleria()
  }, [id])

  const openLightbox = (index: number) => {
    setCurrentIndex(index)
    setLightboxOpen(true)
  }

  const closeLightbox = () => {
    setLightboxOpen(false)
  }

  const goToPrevious = () => {
    if (!galeria) return
    setCurrentIndex((prev) => (prev === 0 ? galeria.items.length - 1 : prev - 1))
  }

  const goToNext = () => {
    if (!galeria) return
    setCurrentIndex((prev) => (prev === galeria.items.length - 1 ? 0 : prev + 1))
  }

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!lightboxOpen) return
      if (e.key === 'Escape') closeLightbox()
      if (e.key === 'ArrowLeft') goToPrevious()
      if (e.key === 'ArrowRight') goToNext()
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [lightboxOpen])

  if (loading) {
    return (
      <div className="py-8">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="animate-pulse">
            <div className="h-8 w-32 bg-gray-200 rounded mb-4" />
            <div className="h-10 w-2/3 bg-gray-200 rounded mb-2" />
            <div className="h-6 w-1/3 bg-gray-200 rounded mb-8" />
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {[...Array(8)].map((_, i) => (
                <div key={i} className="aspect-square bg-gray-200 rounded-lg" />
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !galeria) {
    return (
      <div className="py-8">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <Link
            href="/galeria"
            className="inline-flex items-center text-primary hover:text-primary-light mb-6"
          >
            <ArrowLeft className="w-4 h-4 mr-1" />
            Volver a Galeria
          </Link>
          <div className="text-center py-12 bg-red-50 rounded-lg">
            <p className="text-red-700">{error || 'Galeria no encontrada'}</p>
          </div>
        </div>
      </div>
    )
  }

  const currentItem = galeria.items[currentIndex]

  return (
    <div className="py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Back link */}
        <Link
          href="/galeria"
          className="inline-flex items-center text-primary hover:text-primary-light mb-6"
        >
          <ArrowLeft className="w-4 h-4 mr-1" />
          Volver a Galeria
        </Link>

        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">{galeria.title}</h1>
          {galeria.description && (
            <p className="text-gray-600 mt-2">{galeria.description}</p>
          )}
          <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
            {galeria.event_date && (
              <div className="flex items-center gap-1.5">
                <Calendar className="w-4 h-4" />
                <span>{formatDate(galeria.event_date)}</span>
              </div>
            )}
            <span>{galeria.items.length} {galeria.items.length === 1 ? 'foto' : 'fotos'}</span>
          </div>
        </div>

        {/* Gallery Grid */}
        {galeria.items.length > 0 ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {galeria.items.map((item, index) => (
              <button
                key={item.id}
                onClick={() => openLightbox(index)}
                className="group aspect-square bg-gray-100 rounded-lg overflow-hidden relative"
              >
                {item.file_type === 'video' ? (
                  <>
                    {item.thumbnail_url ? (
                      <img
                        src={item.thumbnail_url}
                        alt={item.caption || 'Video thumbnail'}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center bg-gray-200">
                        <Play className="w-12 h-12 text-gray-400" />
                      </div>
                    )}
                    <div className="absolute inset-0 flex items-center justify-center">
                      <div className="w-12 h-12 bg-black/50 rounded-full flex items-center justify-center group-hover:bg-black/70 transition-colors">
                        <Play className="w-6 h-6 text-white fill-white" />
                      </div>
                    </div>
                  </>
                ) : (
                  <img
                    src={item.thumbnail_url || item.file_url}
                    alt={item.caption || 'Gallery image'}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                )}
              </button>
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-gray-50 rounded-lg">
            <ImageIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-600">Esta galeria no tiene fotos</p>
          </div>
        )}
      </div>

      {/* Lightbox */}
      {lightboxOpen && currentItem && (
        <div className="fixed inset-0 z-50 bg-black flex items-center justify-center">
          {/* Close button */}
          <button
            onClick={closeLightbox}
            className="absolute top-4 right-4 text-white/80 hover:text-white z-10 p-2"
          >
            <X className="w-8 h-8" />
          </button>

          {/* Navigation buttons */}
          {galeria.items.length > 1 && (
            <>
              <button
                onClick={goToPrevious}
                className="absolute left-4 top-1/2 -translate-y-1/2 text-white/80 hover:text-white z-10 p-2"
              >
                <ChevronLeft className="w-10 h-10" />
              </button>
              <button
                onClick={goToNext}
                className="absolute right-4 top-1/2 -translate-y-1/2 text-white/80 hover:text-white z-10 p-2"
              >
                <ChevronRight className="w-10 h-10" />
              </button>
            </>
          )}

          {/* Content */}
          <div className="max-w-6xl max-h-[90vh] w-full mx-4">
            {currentItem.file_type === 'video' ? (
              <video
                src={currentItem.file_url}
                controls
                autoPlay
                className="max-w-full max-h-[85vh] mx-auto"
              />
            ) : (
              <img
                src={currentItem.file_url}
                alt={currentItem.caption || 'Gallery image'}
                className="max-w-full max-h-[85vh] mx-auto object-contain"
              />
            )}

            {/* Caption and counter */}
            <div className="text-center mt-4">
              {currentItem.caption && (
                <p className="text-white/90 mb-2">{currentItem.caption}</p>
              )}
              <p className="text-white/60 text-sm">
                {currentIndex + 1} / {galeria.items.length}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
