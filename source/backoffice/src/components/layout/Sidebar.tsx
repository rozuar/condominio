'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard,
  FileText,
  Calendar,
  AlertTriangle,
  Vote,
  Wallet,
  MessageSquare,
  Image,
  Map,
  Receipt,
  ScrollText,
  FolderOpen,
  Bell,
  LogOut,
  ChevronLeft,
  Menu,
} from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { useState } from 'react'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Comunicados', href: '/comunicados', icon: FileText },
  { name: 'Eventos', href: '/eventos', icon: Calendar },
  { name: 'Emergencias', href: '/emergencias', icon: AlertTriangle },
  { name: 'Votaciones', href: '/votaciones', icon: Vote },
  { name: 'Gastos Comunes', href: '/gastos', icon: Wallet },
  { name: 'Contacto', href: '/contacto', icon: MessageSquare },
  { name: 'Galerias', href: '/galerias', icon: Image },
  { name: 'Mapa', href: '/mapa', icon: Map },
  { name: 'Tesoreria', href: '/tesoreria', icon: Receipt },
  { name: 'Actas', href: '/actas', icon: ScrollText },
  { name: 'Documentos', href: '/documentos', icon: FolderOpen },
  { name: 'Notificaciones', href: '/notificaciones', icon: Bell },
]

export default function Sidebar() {
  const pathname = usePathname()
  const { user, logout } = useAuth()
  const [collapsed, setCollapsed] = useState(false)

  return (
    <aside
      className={`fixed inset-y-0 left-0 z-50 flex flex-col bg-slate-900 transition-all duration-300 ${
        collapsed ? 'w-16' : 'w-64'
      }`}
    >
      {/* Header */}
      <div className="flex h-16 items-center justify-between px-4 border-b border-slate-800">
        {!collapsed && (
          <span className="text-lg font-semibold text-white">Backoffice</span>
        )}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-2 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg"
        >
          {collapsed ? <Menu size={20} /> : <ChevronLeft size={20} />}
        </button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto py-4 px-2">
        <ul className="space-y-1">
          {navigation.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
            return (
              <li key={item.name}>
                <Link
                  href={item.href}
                  className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                    isActive
                      ? 'bg-blue-600 text-white'
                      : 'text-slate-400 hover:text-white hover:bg-slate-800'
                  }`}
                  title={collapsed ? item.name : undefined}
                >
                  <item.icon size={20} />
                  {!collapsed && <span>{item.name}</span>}
                </Link>
              </li>
            )
          })}
        </ul>
      </nav>

      {/* User section */}
      <div className="border-t border-slate-800 p-4">
        {!collapsed && user && (
          <div className="mb-3">
            <p className="text-sm font-medium text-white truncate">{user.name}</p>
            <p className="text-xs text-slate-400 truncate">{user.email}</p>
          </div>
        )}
        <button
          onClick={logout}
          className={`flex items-center gap-3 w-full px-3 py-2 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition-colors ${
            collapsed ? 'justify-center' : ''
          }`}
          title={collapsed ? 'Cerrar sesion' : undefined}
        >
          <LogOut size={20} />
          {!collapsed && <span>Cerrar sesion</span>}
        </button>
      </div>
    </aside>
  )
}
