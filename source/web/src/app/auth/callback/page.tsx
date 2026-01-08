'use client'

import { Suspense } from 'react'
import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import { Loader2, CheckCircle, XCircle } from 'lucide-react'

type CallbackStatus = 'loading' | 'success' | 'error'

function CallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { loginWithTokens } = useAuth()
  const [status, setStatus] = useState<CallbackStatus>('loading')
  const [errorMessage, setErrorMessage] = useState('')

  useEffect(() => {
    const handleCallback = async () => {
      // Check for error from OAuth
      const error = searchParams.get('error')
      if (error) {
        setStatus('error')
        setErrorMessage(decodeURIComponent(error))
        setTimeout(() => router.push('/auth/login'), 3000)
        return
      }

      // Get tokens from URL
      const accessToken = searchParams.get('access_token')
      const refreshToken = searchParams.get('refresh_token')

      if (!accessToken || !refreshToken) {
        setStatus('error')
        setErrorMessage('No se recibieron los tokens de autenticacion')
        setTimeout(() => router.push('/auth/login'), 3000)
        return
      }

      try {
        // Login with tokens
        await loginWithTokens(accessToken, refreshToken)
        setStatus('success')

        // Redirect to home after brief success message
        setTimeout(() => router.push('/'), 1500)
      } catch (err) {
        setStatus('error')
        setErrorMessage(err instanceof Error ? err.message : 'Error al procesar la autenticacion')
        setTimeout(() => router.push('/auth/login'), 3000)
      }
    }

    handleCallback()
  }, [searchParams, loginWithTokens, router])

  return (
    <div className="min-h-[80vh] flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full text-center">
        {status === 'loading' && (
          <div className="bg-white p-8 rounded-lg shadow border">
            <Loader2 className="w-16 h-16 text-primary mx-auto animate-spin" />
            <h1 className="text-2xl font-bold text-gray-900 mt-6">
              Procesando autenticacion
            </h1>
            <p className="text-gray-600 mt-2">
              Por favor espera un momento...
            </p>
          </div>
        )}

        {status === 'success' && (
          <div className="bg-white p-8 rounded-lg shadow border">
            <CheckCircle className="w-16 h-16 text-green-500 mx-auto" />
            <h1 className="text-2xl font-bold text-gray-900 mt-6">
              Autenticacion exitosa
            </h1>
            <p className="text-gray-600 mt-2">
              Bienvenido! Redirigiendo...
            </p>
          </div>
        )}

        {status === 'error' && (
          <div className="bg-white p-8 rounded-lg shadow border">
            <XCircle className="w-16 h-16 text-red-500 mx-auto" />
            <h1 className="text-2xl font-bold text-gray-900 mt-6">
              Error de autenticacion
            </h1>
            <p className="text-red-600 mt-2">
              {errorMessage}
            </p>
            <p className="text-gray-500 text-sm mt-4">
              Redirigiendo al login...
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

function LoadingFallback() {
  return (
    <div className="min-h-[80vh] flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full text-center">
        <div className="bg-white p-8 rounded-lg shadow border">
          <Loader2 className="w-16 h-16 text-primary mx-auto animate-spin" />
          <h1 className="text-2xl font-bold text-gray-900 mt-6">
            Cargando...
          </h1>
        </div>
      </div>
    </div>
  )
}

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <CallbackContent />
    </Suspense>
  )
}
