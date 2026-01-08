'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { TrendingUp, TrendingDown, Wallet, ArrowUpRight, ArrowDownRight } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getTesoreriaResumen, getMovimientos } from '@/lib/api'
import {
  Movimiento,
  MovimientoType,
  MOVIMIENTO_TYPE_LABELS,
  MOVIMIENTO_TYPE_COLORS,
} from '@/types'

const MONTHS = [
  { value: 1, label: 'Enero' },
  { value: 2, label: 'Febrero' },
  { value: 3, label: 'Marzo' },
  { value: 4, label: 'Abril' },
  { value: 5, label: 'Mayo' },
  { value: 6, label: 'Junio' },
  { value: 7, label: 'Julio' },
  { value: 8, label: 'Agosto' },
  { value: 9, label: 'Septiembre' },
  { value: 10, label: 'Octubre' },
  { value: 11, label: 'Noviembre' },
  { value: 12, label: 'Diciembre' },
]

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('es-CL', {
    style: 'currency',
    currency: 'CLP',
    minimumFractionDigits: 0,
  }).format(amount)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-CL', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export default function TesoreriaPage() {
  const router = useRouter()
  const { isAuthenticated, isLoading: authLoading, getToken } = useAuth()

  const [resumen, setResumen] = useState<{
    total_ingresos: number
    total_egresos: number
    balance: number
    movimientos_count: number
  } | null>(null)
  const [movimientos, setMovimientos] = useState<Movimiento[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Filters
  const currentYear = new Date().getFullYear()
  const [selectedYear, setSelectedYear] = useState<number | undefined>(undefined)
  const [selectedMonth, setSelectedMonth] = useState<number | undefined>(undefined)
  const [selectedType, setSelectedType] = useState<MovimientoType | ''>('')

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/auth/login')
    }
  }, [authLoading, isAuthenticated, router])

  useEffect(() => {
    if (!isAuthenticated) return

    const fetchData = async () => {
      setIsLoading(true)
      setError(null)

      const token = getToken()
      if (!token) {
        setError('No hay sesión activa')
        setIsLoading(false)
        return
      }

      try {
        const [resumenData, movimientosData] = await Promise.all([
          getTesoreriaResumen(token),
          getMovimientos(token, {
            year: selectedYear,
            month: selectedMonth,
            type: selectedType || undefined,
            per_page: 50,
          }),
        ])

        setResumen(resumenData)
        setMovimientos(movimientosData.movimientos || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error al cargar los datos')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken, selectedYear, selectedMonth, selectedType])

  if (authLoading || (!isAuthenticated && !error)) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-gray-500">Cargando...</div>
      </div>
    )
  }

  return (
    <div className="py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Tesorería</h1>
          <p className="text-gray-600 mt-2">Balance y movimientos financieros de la comunidad</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {/* Summary Cards */}
        {resumen && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-white rounded-lg shadow border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Total Ingresos</p>
                  <p className="text-2xl font-bold text-green-600">
                    {formatCurrency(resumen.total_ingresos)}
                  </p>
                </div>
                <div className="p-3 bg-green-100 rounded-full">
                  <TrendingUp className="h-6 w-6 text-green-600" />
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Total Egresos</p>
                  <p className="text-2xl font-bold text-red-600">
                    {formatCurrency(resumen.total_egresos)}
                  </p>
                </div>
                <div className="p-3 bg-red-100 rounded-full">
                  <TrendingDown className="h-6 w-6 text-red-600" />
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Balance Actual</p>
                  <p className={`text-2xl font-bold ${resumen.balance >= 0 ? 'text-primary' : 'text-red-600'}`}>
                    {formatCurrency(resumen.balance)}
                  </p>
                </div>
                <div className="p-3 bg-primary/10 rounded-full">
                  <Wallet className="h-6 w-6 text-primary" />
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Filters */}
        <div className="bg-white rounded-lg shadow border p-4 mb-6">
          <div className="flex flex-wrap gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Año</label>
              <select
                value={selectedYear || ''}
                onChange={(e) => setSelectedYear(e.target.value ? parseInt(e.target.value) : undefined)}
                className="px-3 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
              >
                <option value="">Todos</option>
                {[currentYear, currentYear - 1, currentYear - 2].map((year) => (
                  <option key={year} value={year}>
                    {year}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Mes</label>
              <select
                value={selectedMonth || ''}
                onChange={(e) => setSelectedMonth(e.target.value ? parseInt(e.target.value) : undefined)}
                className="px-3 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
              >
                <option value="">Todos</option>
                {MONTHS.map((month) => (
                  <option key={month.value} value={month.value}>
                    {month.label}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Tipo</label>
              <select
                value={selectedType}
                onChange={(e) => setSelectedType(e.target.value as MovimientoType | '')}
                className="px-3 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
              >
                <option value="">Todos</option>
                <option value="ingreso">Ingresos</option>
                <option value="egreso">Egresos</option>
              </select>
            </div>
          </div>
        </div>

        {/* Movements List */}
        <div className="bg-white rounded-lg shadow border overflow-hidden">
          <div className="px-6 py-4 border-b">
            <h2 className="text-lg font-semibold">Movimientos</h2>
          </div>

          {isLoading ? (
            <div className="p-8 text-center text-gray-500">Cargando movimientos...</div>
          ) : movimientos.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No hay movimientos para los filtros seleccionados
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Fecha
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Descripción
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Categoría
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Tipo
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Monto
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {movimientos.map((mov) => (
                    <tr key={mov.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(mov.date)}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900">
                        {mov.description}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {mov.category}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium ${MOVIMIENTO_TYPE_COLORS[mov.type]}`}>
                          {mov.type === 'ingreso' ? (
                            <ArrowUpRight className="h-3 w-3" />
                          ) : (
                            <ArrowDownRight className="h-3 w-3" />
                          )}
                          {MOVIMIENTO_TYPE_LABELS[mov.type]}
                        </span>
                      </td>
                      <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium text-right ${mov.type === 'ingreso' ? 'text-green-600' : 'text-red-600'}`}>
                        {mov.type === 'ingreso' ? '+' : '-'}{formatCurrency(mov.amount)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
