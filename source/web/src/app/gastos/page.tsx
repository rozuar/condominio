'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import {
  Receipt,
  AlertCircle,
  CheckCircle,
  Clock,
  CreditCard,
  Calendar,
  Home,
} from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { getMiEstadoCuenta } from '@/lib/api'
import {
  MiEstadoCuenta,
  GastoComun,
  PAGO_STATUS_LABELS,
  PAGO_STATUS_COLORS,
} from '@/types'

const MONTHS = [
  '', 'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
  'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
]

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('es-CL', {
    style: 'currency',
    currency: 'CLP',
    minimumFractionDigits: 0,
  }).format(amount)
}

function formatPeriodo(gasto: GastoComun): string {
  if (gasto.periodo) {
    return `${MONTHS[gasto.periodo.month]} ${gasto.periodo.year}`
  }
  return 'Periodo'
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-CL', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

export default function GastosPage() {
  const router = useRouter()
  const { isAuthenticated, isLoading: authLoading, getToken, user } = useAuth()

  const [estadoCuenta, setEstadoCuenta] = useState<MiEstadoCuenta | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

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
        setError('No hay sesion activa')
        setIsLoading(false)
        return
      }

      try {
        const data = await getMiEstadoCuenta(token)
        setEstadoCuenta(data)
      } catch (err) {
        if (err instanceof Error && err.message.includes('no associated parcela')) {
          setError('Tu cuenta no tiene una parcela asociada. Contacta a la directiva.')
        } else {
          setError(err instanceof Error ? err.message : 'Error al cargar los datos')
        }
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAuthenticated, getToken])

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
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Gastos Comunes</h1>
          <p className="text-gray-600 mt-2">Estado de cuenta de tu parcela</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center gap-3">
            <AlertCircle className="w-5 h-5 flex-shrink-0" />
            {error}
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12 text-gray-500">Cargando estado de cuenta...</div>
        ) : estadoCuenta ? (
          <>
            {/* Parcela Info */}
            <div className="bg-white rounded-lg shadow border p-6 mb-8">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-primary/10 rounded-full">
                  <Home className="w-8 h-8 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Parcela</p>
                  <p className="text-2xl font-bold text-gray-900">N {estadoCuenta.parcela_numero}</p>
                </div>
              </div>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
              {/* Pendiente */}
              <div className={`rounded-lg shadow border p-6 ${estadoCuenta.total_pendiente > 0 ? 'bg-red-50 border-red-200' : 'bg-green-50 border-green-200'}`}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600">Total Pendiente</p>
                    <p className={`text-3xl font-bold ${estadoCuenta.total_pendiente > 0 ? 'text-red-600' : 'text-green-600'}`}>
                      {formatCurrency(estadoCuenta.total_pendiente)}
                    </p>
                    {estadoCuenta.gastos_pendientes.length > 0 && (
                      <p className="text-sm text-gray-500 mt-1">
                        {estadoCuenta.gastos_pendientes.length} periodo{estadoCuenta.gastos_pendientes.length !== 1 ? 's' : ''} pendiente{estadoCuenta.gastos_pendientes.length !== 1 ? 's' : ''}
                      </p>
                    )}
                  </div>
                  <div className={`p-3 rounded-full ${estadoCuenta.total_pendiente > 0 ? 'bg-red-100' : 'bg-green-100'}`}>
                    {estadoCuenta.total_pendiente > 0 ? (
                      <AlertCircle className="w-8 h-8 text-red-600" />
                    ) : (
                      <CheckCircle className="w-8 h-8 text-green-600" />
                    )}
                  </div>
                </div>
              </div>

              {/* Pagado */}
              <div className="bg-white rounded-lg shadow border p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-500">Total Pagado (ultimos 12 meses)</p>
                    <p className="text-3xl font-bold text-gray-900">
                      {formatCurrency(estadoCuenta.total_pagado)}
                    </p>
                    <p className="text-sm text-gray-500 mt-1">
                      {estadoCuenta.gastos_pagados.length} pago{estadoCuenta.gastos_pagados.length !== 1 ? 's' : ''} registrado{estadoCuenta.gastos_pagados.length !== 1 ? 's' : ''}
                    </p>
                  </div>
                  <div className="p-3 bg-blue-100 rounded-full">
                    <CreditCard className="w-8 h-8 text-blue-600" />
                  </div>
                </div>
              </div>
            </div>

            {/* Pending Gastos */}
            {estadoCuenta.gastos_pendientes.length > 0 && (
              <div className="bg-white rounded-lg shadow border overflow-hidden mb-8">
                <div className="px-6 py-4 border-b bg-red-50">
                  <h2 className="text-lg font-semibold text-red-800 flex items-center gap-2">
                    <AlertCircle className="w-5 h-5" />
                    Gastos Pendientes
                  </h2>
                </div>
                <div className="divide-y">
                  {estadoCuenta.gastos_pendientes.map((gasto) => (
                    <div key={gasto.id} className="p-4 hover:bg-gray-50">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div className="p-2 bg-gray-100 rounded">
                            <Calendar className="w-5 h-5 text-gray-600" />
                          </div>
                          <div>
                            <p className="font-medium text-gray-900">
                              {formatPeriodo(gasto)}
                            </p>
                            <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${PAGO_STATUS_COLORS[gasto.status]}`}>
                              {PAGO_STATUS_LABELS[gasto.status]}
                            </span>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-lg font-bold text-red-600">
                            {formatCurrency(gasto.monto - gasto.monto_pagado)}
                          </p>
                          {gasto.monto_pagado > 0 && (
                            <p className="text-xs text-gray-500">
                              Abonado: {formatCurrency(gasto.monto_pagado)}
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
                <div className="px-6 py-4 border-t bg-gray-50">
                  <p className="text-sm text-gray-600">
                    Para pagar, realiza una transferencia a la cuenta de la comunidad y notifica a la directiva con el comprobante.
                  </p>
                </div>
              </div>
            )}

            {/* No pending message */}
            {estadoCuenta.gastos_pendientes.length === 0 && (
              <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center mb-8">
                <CheckCircle className="w-12 h-12 text-green-600 mx-auto mb-4" />
                <p className="text-green-800 font-medium text-lg">Estas al dia con tus gastos comunes</p>
                <p className="text-green-600 mt-1">No tienes pagos pendientes</p>
              </div>
            )}

            {/* Payment History */}
            {estadoCuenta.gastos_pagados.length > 0 && (
              <div className="bg-white rounded-lg shadow border overflow-hidden">
                <div className="px-6 py-4 border-b">
                  <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
                    <Receipt className="w-5 h-5" />
                    Historial de Pagos
                  </h2>
                </div>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          Periodo
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          Fecha Pago
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          Metodo
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                          Monto
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {estadoCuenta.gastos_pagados.map((gasto) => (
                        <tr key={gasto.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {formatPeriodo(gasto)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {gasto.fecha_pago ? formatDate(gasto.fecha_pago) : '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {gasto.metodo_pago || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-right text-green-600">
                            {formatCurrency(gasto.monto_pagado)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </>
        ) : null}
      </div>
    </div>
  )
}
