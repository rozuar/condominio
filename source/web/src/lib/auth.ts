const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export function getGoogleAuthUrl(): string {
  return `${API_URL}/auth/google`
}

export interface User {
  id: string
  email: string
  name: string
  role: 'visitor' | 'vecino' | 'directiva'
  parcela_id?: number
  email_verified: boolean
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Login failed')
  }

  return response.json()
}

export async function register(email: string, password: string, name: string): Promise<AuthResponse> {
  const response = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, name }),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Registration failed')
  }

  return response.json()
}

export async function refreshTokens(refreshToken: string): Promise<AuthResponse> {
  const response = await fetch(`${API_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })

  if (!response.ok) {
    throw new Error('Token refresh failed')
  }

  return response.json()
}

export async function getMe(token: string): Promise<User> {
  const response = await fetch(`${API_URL}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })

  if (!response.ok) {
    throw new Error('Failed to get user')
  }

  return response.json()
}

// Storage helpers
export function saveTokens(accessToken: string, refreshToken: string) {
  if (typeof window !== 'undefined') {
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
  }
}

export function getAccessToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('access_token')
  }
  return null
}

export function getRefreshToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('refresh_token')
  }
  return null
}

export function clearTokens() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }
}
