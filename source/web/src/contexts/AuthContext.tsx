'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import {
  User,
  login as apiLogin,
  register as apiRegister,
  getMe,
  saveTokens,
  getAccessToken,
  getRefreshToken,
  clearTokens,
  refreshTokens,
} from '@/lib/auth'

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  loginWithTokens: (accessToken: string, refreshToken: string) => Promise<void>
  logout: () => void
  getToken: () => string | null
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const initAuth = async () => {
      const token = getAccessToken()
      if (token) {
        try {
          const userData = await getMe(token)
          setUser(userData)
        } catch (error) {
          // Try refresh token
          const refresh = getRefreshToken()
          if (refresh) {
            try {
              const response = await refreshTokens(refresh)
              saveTokens(response.access_token, response.refresh_token)
              setUser(response.user)
            } catch {
              clearTokens()
            }
          } else {
            clearTokens()
          }
        }
      }
      setIsLoading(false)
    }

    initAuth()
  }, [])

  const login = async (email: string, password: string) => {
    const response = await apiLogin(email, password)
    saveTokens(response.access_token, response.refresh_token)
    setUser(response.user)
  }

  const register = async (email: string, password: string, name: string) => {
    const response = await apiRegister(email, password, name)
    saveTokens(response.access_token, response.refresh_token)
    setUser(response.user)
  }

  const loginWithTokens = async (accessToken: string, refreshToken: string) => {
    saveTokens(accessToken, refreshToken)
    const userData = await getMe(accessToken)
    setUser(userData)
  }

  const logout = () => {
    clearTokens()
    setUser(null)
  }

  const getToken = () => getAccessToken()

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        loginWithTokens,
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
