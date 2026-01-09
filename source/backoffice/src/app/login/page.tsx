'use client'

import { useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { Loader2, AlertCircle, Shield, Crown, Wallet, ClipboardList } from 'lucide-react'

export default function LoginPage() {
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      await login(email, password)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al iniciar sesion')
    } finally {
      setIsLoading(false)
    }
  }

  const devProfiles = [
    {
      label: 'Admin',
      email: 'admin@vinapelvin.cl',
      password: 'admin123',
      icon: <Shield className="h-4 w-4" />,
    },
    {
      label: 'Presidente',
      email: 'presidente@vinapelvin.cl',
      password: 'admin123',
      icon: <Crown className="h-4 w-4" />,
    },
    {
      label: 'Tesorero',
      email: 'tesorero@vinapelvin.cl',
      password: 'admin123',
      icon: <Wallet className="h-4 w-4" />,
    },
    {
      label: 'Secretaria',
      email: 'secretaria@vinapelvin.cl',
      password: 'admin123',
      icon: <ClipboardList className="h-4 w-4" />,
    },
  ] as const

  const handleDevProfileLogin = async (profile: (typeof devProfiles)[number]) => {
    setEmail(profile.email)
    setPassword(profile.password)
    setError('')
    setIsLoading(true)

    try {
      await login(profile.email, profile.password)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al iniciar sesion')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 px-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-xl shadow-xl p-8">
          <div className="text-center mb-8">
            <h1 className="text-2xl font-bold text-gray-900">Backoffice</h1>
            <p className="text-gray-500 mt-1">Comunidad Vina Pelvin</p>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg flex items-start gap-3">
              <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                placeholder="admin@vinapelvin.cl"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                Contrasena
              </label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                placeholder="********"
              />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-blue-600 text-white py-2.5 px-4 rounded-lg font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
              {isLoading ? 'Ingresando...' : 'Ingresar'}
            </button>
          </form>

          {process.env.NODE_ENV === 'development' && (
            <div className="mt-6 pt-6 border-t border-gray-200">
              <p className="text-xs text-gray-500 text-center mb-3">Modo desarrollo</p>
              <div className="grid grid-cols-2 gap-2">
                {devProfiles.map((profile) => (
                  <button
                    key={profile.email}
                    type="button"
                    onClick={() => handleDevProfileLogin(profile)}
                    disabled={isLoading}
                    className="w-full bg-gray-100 text-gray-700 py-2 px-3 rounded-lg text-sm hover:bg-gray-200 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                    title={`Ingresar como ${profile.label}`}
                  >
                    {profile.icon}
                    {profile.label}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>

        <p className="text-center text-slate-500 text-sm mt-6">
          Solo usuarios con rol directiva pueden acceder
        </p>
      </div>
    </div>
  )
}
