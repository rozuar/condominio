'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { login as apiLogin, getMe, refreshTokens } from '@/lib/api'
import type { User } from '@/types'

const BACKOFFICE_ALLOWED_ROLES = new Set(['directiva', 'admin'])

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  getToken: () => string | null
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    if (typeof window === 'undefined') {
      setIsLoading(false)
      return
    }

    const accessToken = localStorage.getItem('access_token')
    const refreshToken = localStorage.getItem('refresh_token')

    if (!accessToken) {
      setIsLoading(false)
      return
    }

    try {
      const { user } = await getMe(accessToken)
      if (!BACKOFFICE_ALLOWED_ROLES.has(user.role)) {
        throw new Error('Acceso no autorizado')
      }
      setUser(user)
    } catch {
      if (refreshToken) {
        try {
          const tokens = await refreshTokens(refreshToken)
          localStorage.setItem('access_token', tokens.access_token)
          localStorage.setItem('refresh_token', tokens.refresh_token)
          const { user } = await getMe(tokens.access_token)
          if (!BACKOFFICE_ALLOWED_ROLES.has(user.role)) {
            throw new Error('Acceso no autorizado')
          }
          setUser(user)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
        }
      }
    } finally {
      setIsLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    const response = await apiLogin(email, password)

    if (!BACKOFFICE_ALLOWED_ROLES.has(response.user.role)) {
      throw new Error('Solo usuarios con rol directiva/admin pueden acceder al backoffice')
    }

    localStorage.setItem('access_token', response.access_token)
    localStorage.setItem('refresh_token', response.refresh_token)
    setUser(response.user)
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
  }

  const getToken = () => {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('access_token')
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        logout,
        getToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
