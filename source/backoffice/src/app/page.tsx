'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getDashboardStats } from '@/lib/api'
import {
  FileText,
  Calendar,
  AlertTriangle,
  Vote,
  MessageSquare,
  Loader2,
} from 'lucide-react'
import Link from 'next/link'

interface Stats {
  totalComunicados: number
  totalEventos: number
  emergenciasActivas: number
  votacionesActivas: number
  mensajesPendientes: number
}

export default function DashboardPage() {
  const { user, getToken } = useAuth()
  const [stats, setStats] = useState<Stats | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchStats = async () => {
      const token = getToken()
      if (!token) return

      try {
        const data = await getDashboardStats(token)
        setStats(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar estadisticas')
      } finally {
        setIsLoading(false)
      }
    }

    fetchStats()
  }, [getToken])

  const statCards = [
    {
      name: 'Comunicados',
      value: stats?.totalComunicados || 0,
      icon: FileText,
      href: '/comunicados',
      color: 'bg-blue-500',
    },
    {
      name: 'Eventos',
      value: stats?.totalEventos || 0,
      icon: Calendar,
      href: '/eventos',
      color: 'bg-purple-500',
    },
    {
      name: 'Emergencias Activas',
      value: stats?.emergenciasActivas || 0,
      icon: AlertTriangle,
      href: '/emergencias',
      color: 'bg-red-500',
    },
    {
      name: 'Votaciones Activas',
      value: stats?.votacionesActivas || 0,
      icon: Vote,
      href: '/votaciones',
      color: 'bg-green-500',
    },
    {
      name: 'Mensajes Pendientes',
      value: stats?.mensajesPendientes || 0,
      icon: MessageSquare,
      href: '/contacto',
      color: 'bg-amber-500',
    },
  ]

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 mt-1">
          Bienvenido, {user?.name}
        </p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          {error}
        </div>
      )}

      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">
          {statCards.map((card) => (
            <Link
              key={card.name}
              href={card.href}
              className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">{card.name}</p>
                  <p className="text-3xl font-bold text-gray-900 mt-1">{card.value}</p>
                </div>
                <div className={`${card.color} p-3 rounded-lg`}>
                  <card.icon className="h-6 w-6 text-white" />
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      <div className="mt-12 grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Acciones Rapidas</h2>
          <div className="space-y-3">
            <Link
              href="/comunicados?action=new"
              className="block p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div className="flex items-center gap-3">
                <FileText className="h-5 w-5 text-blue-600" />
                <span className="font-medium">Crear Comunicado</span>
              </div>
            </Link>
            <Link
              href="/eventos?action=new"
              className="block p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div className="flex items-center gap-3">
                <Calendar className="h-5 w-5 text-purple-600" />
                <span className="font-medium">Crear Evento</span>
              </div>
            </Link>
            <Link
              href="/emergencias?action=new"
              className="block p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div className="flex items-center gap-3">
                <AlertTriangle className="h-5 w-5 text-red-600" />
                <span className="font-medium">Crear Emergencia</span>
              </div>
            </Link>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Informacion</h2>
          <div className="space-y-4 text-sm text-gray-600">
            <p>
              Este es el panel de administracion de la Comunidad Vina Pelvin.
              Desde aqui puedes gestionar todos los contenidos del portal.
            </p>
            <p>
              Usa el menu lateral para navegar entre las diferentes secciones.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
