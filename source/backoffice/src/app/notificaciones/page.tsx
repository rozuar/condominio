'use client'

import { useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { broadcastNotificacion } from '@/lib/api'
import { Send, Bell, Users } from 'lucide-react'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Textarea from '@/components/ui/Textarea'

export default function NotificacionesPage() {
  const { getToken } = useAuth()
  const [isSending, setIsSending] = useState(false)
  const [success, setSuccess] = useState('')
  const [error, setError] = useState('')

  const [formData, setFormData] = useState({ title: '', body: '', type: 'general' })

  const handleBroadcast = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSending(true)
    setError('')
    setSuccess('')

    try {
      const result = await broadcastNotificacion(token, formData)
      setSuccess(`Notificacion enviada a ${result.count} usuarios`)
      setFormData({ title: '', body: '', type: 'general' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al enviar')
    } finally {
      setIsSending(false)
    }
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Notificaciones</h1>
        <p className="text-gray-500 mt-1">Enviar notificaciones a los vecinos</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-white rounded-xl shadow-sm border p-6">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-blue-100 rounded-lg">
              <Users className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <h2 className="text-lg font-semibold">Notificacion Masiva</h2>
              <p className="text-sm text-gray-500">Enviar a todos los vecinos registrados</p>
            </div>
          </div>

          {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}
          {success && <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg text-green-700">{success}</div>}

          <form onSubmit={handleBroadcast} className="space-y-6">
            <Input label="Titulo" value={formData.title} onChange={(e) => setFormData({ ...formData, title: e.target.value })} required placeholder="Titulo de la notificacion" />
            <Textarea label="Mensaje" value={formData.body} onChange={(e) => setFormData({ ...formData, body: e.target.value })} required rows={4} placeholder="Escribe el mensaje..." />
            <Input label="Tipo" value={formData.type} onChange={(e) => setFormData({ ...formData, type: e.target.value })} placeholder="general, urgente, recordatorio" />
            <Button type="submit" loading={isSending} icon={<Send size={18} />} className="w-full">
              Enviar Notificacion
            </Button>
          </form>
        </div>

        <div className="bg-white rounded-xl shadow-sm border p-6">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-amber-100 rounded-lg">
              <Bell className="h-6 w-6 text-amber-600" />
            </div>
            <div>
              <h2 className="text-lg font-semibold">Tipos de Notificaciones</h2>
              <p className="text-sm text-gray-500">Categorias disponibles</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="font-medium">general</p>
              <p className="text-sm text-gray-500">Avisos generales de la comunidad</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="font-medium">urgente</p>
              <p className="text-sm text-gray-500">Notificaciones importantes que requieren atencion</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="font-medium">recordatorio</p>
              <p className="text-sm text-gray-500">Recordatorios de pagos, reuniones, etc.</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="font-medium">evento</p>
              <p className="text-sm text-gray-500">Notificaciones sobre eventos proximos</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
