'use client'

import { useState, useEffect } from 'react'
import {
  Mail,
  Send,
  CheckCircle,
  Clock,
  MessageSquare,
  User,
  AlertCircle,
} from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { enviarMensajeContacto, getMisMensajes } from '@/lib/api'
import {
  MensajeContacto,
  CONTACTO_STATUS_LABELS,
} from '@/types'

export default function ContactoPage() {
  const { isAuthenticated, getToken, user } = useAuth()

  // Form state
  const [nombre, setNombre] = useState('')
  const [email, setEmail] = useState('')
  const [asunto, setAsunto] = useState('')
  const [mensaje, setMensaje] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitSuccess, setSubmitSuccess] = useState(false)

  // Messages state (for authenticated users)
  const [misMensajes, setMisMensajes] = useState<MensajeContacto[]>([])
  const [loadingMensajes, setLoadingMensajes] = useState(false)

  // Pre-fill form for authenticated users
  useEffect(() => {
    if (user) {
      setNombre(user.name || '')
      setEmail(user.email || '')
    }
  }, [user])

  // Load user's messages
  useEffect(() => {
    if (!isAuthenticated) return

    const fetchMensajes = async () => {
      setLoadingMensajes(true)
      const token = getToken()
      if (!token) return

      try {
        const data = await getMisMensajes(token)
        setMisMensajes(data.mensajes || [])
      } catch (err) {
        console.error('Error loading messages:', err)
      } finally {
        setLoadingMensajes(false)
      }
    }

    fetchMensajes()
  }, [isAuthenticated, getToken, submitSuccess])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    setSubmitError(null)

    try {
      await enviarMensajeContacto({
        nombre: nombre.trim(),
        email: email.trim(),
        asunto: asunto.trim(),
        mensaje: mensaje.trim(),
      })

      setSubmitSuccess(true)
      setAsunto('')
      setMensaje('')
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'Error al enviar el mensaje')
    } finally {
      setIsSubmitting(false)
    }
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('es-CL', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <Clock className="w-4 h-4 text-yellow-500" />
      case 'read':
        return <CheckCircle className="w-4 h-4 text-blue-500" />
      case 'replied':
        return <MessageSquare className="w-4 h-4 text-green-500" />
      default:
        return <Mail className="w-4 h-4 text-gray-500" />
    }
  }

  return (
    <div className="py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Contacto Directiva</h1>
          <p className="text-gray-600 mt-2">
            Envia un mensaje a la directiva de la comunidad
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Contact Form */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
              <Mail className="w-5 h-5 text-primary" />
              Enviar Mensaje
            </h2>

            {submitSuccess && (
              <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg flex items-start gap-3">
                <CheckCircle className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-green-800 font-medium">Mensaje enviado correctamente</p>
                  <p className="text-green-700 text-sm mt-1">
                    La directiva revisara tu mensaje y te respondera a la brevedad.
                  </p>
                </div>
              </div>
            )}

            {submitError && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                <p className="text-red-700">{submitError}</p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Nombre *
                </label>
                <input
                  type="text"
                  value={nombre}
                  onChange={(e) => setNombre(e.target.value)}
                  required
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                  placeholder="Tu nombre completo"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email *
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                  placeholder="tu@email.com"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Asunto *
                </label>
                <input
                  type="text"
                  value={asunto}
                  onChange={(e) => setAsunto(e.target.value)}
                  required
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                  placeholder="Asunto del mensaje"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Mensaje *
                </label>
                <textarea
                  value={mensaje}
                  onChange={(e) => setMensaje(e.target.value)}
                  required
                  rows={5}
                  className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary resize-none"
                  placeholder="Escribe tu mensaje aqui..."
                />
              </div>

              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-primary text-white font-medium rounded-lg hover:bg-primary-light disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isSubmitting ? (
                  'Enviando...'
                ) : (
                  <>
                    <Send className="w-4 h-4" />
                    Enviar Mensaje
                  </>
                )}
              </button>
            </form>
          </div>

          {/* Info & Messages */}
          <div className="space-y-6">
            {/* Contact Info */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <h2 className="text-xl font-semibold mb-4">Informacion de Contacto</h2>
              <div className="space-y-4 text-gray-600">
                <div className="flex items-start gap-3">
                  <Mail className="w-5 h-5 text-primary mt-0.5" />
                  <div>
                    <p className="font-medium text-gray-900">Email</p>
                    <p>directiva@vinapelvin.cl</p>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <User className="w-5 h-5 text-primary mt-0.5" />
                  <div>
                    <p className="font-medium text-gray-900">Directiva</p>
                    <p>Comunidad Vina Pelvin</p>
                  </div>
                </div>
              </div>
              <div className="mt-6 p-4 bg-blue-50 rounded-lg">
                <p className="text-sm text-blue-800">
                  Los mensajes son revisados regularmente por la directiva.
                  Recibirás una respuesta a tu correo electronico.
                </p>
              </div>
            </div>

            {/* My Messages (for authenticated users) */}
            {isAuthenticated && (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                  <MessageSquare className="w-5 h-5 text-primary" />
                  Mis Mensajes
                </h2>

                {loadingMensajes ? (
                  <p className="text-gray-500 text-center py-4">Cargando...</p>
                ) : misMensajes.length === 0 ? (
                  <p className="text-gray-500 text-center py-4">
                    No has enviado mensajes aun
                  </p>
                ) : (
                  <div className="space-y-3 max-h-96 overflow-y-auto">
                    {misMensajes.map((msg) => (
                      <div
                        key={msg.id}
                        className="p-3 border rounded-lg hover:bg-gray-50"
                      >
                        <div className="flex items-start justify-between gap-2">
                          <div className="flex-1 min-w-0">
                            <p className="font-medium text-gray-900 truncate">
                              {msg.asunto}
                            </p>
                            <p className="text-sm text-gray-500 truncate">
                              {msg.mensaje}
                            </p>
                          </div>
                          <div className="flex items-center gap-1 flex-shrink-0">
                            {getStatusIcon(msg.status)}
                            <span className="text-xs text-gray-500">
                              {CONTACTO_STATUS_LABELS[msg.status]}
                            </span>
                          </div>
                        </div>
                        <p className="text-xs text-gray-400 mt-1">
                          {formatDate(msg.created_at)}
                        </p>
                        {msg.respuesta && (
                          <div className="mt-2 p-2 bg-green-50 rounded text-sm">
                            <p className="text-green-800 font-medium text-xs mb-1">Respuesta:</p>
                            <p className="text-green-700">{msg.respuesta}</p>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
