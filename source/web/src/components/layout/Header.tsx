'use client'

import Link from 'next/link'
import Image from 'next/image'
import { useState } from 'react'
import { Menu, X, User, LogOut } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

const publicNavigation = [
  { name: 'Inicio', href: '/' },
  { name: 'Comunicados', href: '/comunicados' },
  { name: 'Calendario', href: '/calendario' },
  { name: 'Contacto Directiva', href: '/contacto' },
]

const privateNavigation = [
  { name: 'Tesorería', href: '/tesoreria' },
  { name: 'Actas', href: '/actas' },
  { name: 'Documentos', href: '/documentos' },
  { name: 'Votaciones', href: '/votaciones' },
]

export default function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const { user, isAuthenticated, logout, isLoading } = useAuth()

  const navigation = isAuthenticated
    ? [...publicNavigation, ...privateNavigation]
    : publicNavigation

  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur border-b border-gray-200">
      <nav className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Link href="/" className="flex items-center gap-2">
              <Image
                src="/logo.svg"
                alt="Comunidad Viña Pelvin"
                width={28}
                height={28}
                priority
              />
              <span className="text-base sm:text-lg font-semibold text-gray-900 tracking-tight">
                Comunidad Viña Pelvin
              </span>
            </Link>
          </div>

          {/* Desktop navigation */}
          <div className="hidden md:block">
            <div className="flex items-center gap-6">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className="text-sm text-gray-600 hover:text-gray-900 transition-colors"
                >
                  {item.name}
                </Link>
              ))}

              {!isLoading && (
                isAuthenticated ? (
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-gray-600 flex items-center gap-1">
                      <User size={16} />
                      {user?.name}
                    </span>
                    <button
                      onClick={logout}
                      className="flex items-center gap-1 text-gray-600 hover:text-gray-900"
                    >
                      <LogOut size={16} />
                      Salir
                    </button>
                  </div>
                ) : (
                  <Link
                    href="/auth/login"
                    className="bg-primary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-primary-dark transition-colors shadow-sm"
                  >
                    Ingresar
                  </Link>
                )
              )}
            </div>
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="p-2 text-gray-700"
            >
              {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>

        {/* Mobile navigation */}
        {mobileMenuOpen && (
          <div className="md:hidden pb-4">
            <div className="flex flex-col gap-2">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className="text-gray-700 hover:text-gray-900 py-2"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  {item.name}
                </Link>
              ))}

              {!isLoading && (
                isAuthenticated ? (
                  <>
                    <div className="border-t border-gray-200 my-2 pt-2">
                      <span className="text-sm text-gray-600">{user?.name}</span>
                    </div>
                    <button
                      onClick={() => {
                        logout()
                        setMobileMenuOpen(false)
                      }}
                      className="text-gray-700 hover:text-gray-900 py-2 text-left"
                    >
                      Cerrar sesión
                    </button>
                  </>
                ) : (
                  <Link
                    href="/auth/login"
                    className="bg-primary text-white px-4 py-2 rounded-lg font-medium text-center mt-2 shadow-sm"
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    Ingresar
                  </Link>
                )
              )}
            </div>
          </div>
        )}
      </nav>
    </header>
  )
}
