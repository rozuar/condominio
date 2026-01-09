'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getPeriodos, createPeriodo, getGastosPeriodo, registrarPago, getResumenPeriodo } from '@/lib/api'
import type { PeriodoGasto, GastoComun, ResumenGastos } from '@/types'
import { PAGO_STATUS_LABELS, PAGO_STATUS_COLORS } from '@/types'
import { Plus, Loader2, DollarSign, CheckCircle } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'

const MESES = ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio', 'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre']

export default function GastosPage() {
  const { getToken } = useAuth()
  const [periodos, setPeriodos] = useState<PeriodoGasto[]>([])
  const [selectedPeriodo, setSelectedPeriodo] = useState<PeriodoGasto | null>(null)
  const [gastos, setGastos] = useState<GastoComun[]>([])
  const [resumen, setResumen] = useState<ResumenGastos | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingGastos, setIsLoadingGastos] = useState(false)
  const [error, setError] = useState('')

  const [showNewPeriodoModal, setShowNewPeriodoModal] = useState(false)
  const [showPagoModal, setShowPagoModal] = useState(false)
  const [selectedGasto, setSelectedGasto] = useState<GastoComun | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  const [newPeriodo, setNewPeriodo] = useState({ year: new Date().getFullYear(), month: new Date().getMonth() + 1, monto_base: 0, fecha_vencimiento: '' })
  const [pagoData, setPagoData] = useState({ monto_pagado: 0, metodo_pago: '', referencia_pago: '' })

  const fetchPeriodos = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const data = await getPeriodos(token, { per_page: 50 })
      setPeriodos(data.periodos)
      if (data.periodos.length > 0 && !selectedPeriodo) {
        selectPeriodo(data.periodos[0])
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  const selectPeriodo = async (periodo: PeriodoGasto) => {
    const token = getToken()
    if (!token) return

    setSelectedPeriodo(periodo)
    setIsLoadingGastos(true)
    try {
      const [gastosData, resumenData] = await Promise.all([
        getGastosPeriodo(periodo.id, token),
        getResumenPeriodo(periodo.id, token),
      ])
      setGastos(gastosData.gastos)
      setResumen(resumenData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar gastos')
    } finally {
      setIsLoadingGastos(false)
    }
  }

  useEffect(() => {
    fetchPeriodos()
  }, [])

  const handleCreatePeriodo = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await createPeriodo(token, newPeriodo)
      setShowNewPeriodoModal(false)
      fetchPeriodos()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al crear periodo')
    } finally {
      setIsSaving(false)
    }
  }

  const handleRegistrarPago = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedGasto) return
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await registrarPago(selectedGasto.id, token, pagoData)
      setShowPagoModal(false)
      if (selectedPeriodo) selectPeriodo(selectedPeriodo)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al registrar pago')
    } finally {
      setIsSaving(false)
    }
  }

  const openPagoModal = (gasto: GastoComun) => {
    setSelectedGasto(gasto)
    setPagoData({ monto_pagado: gasto.monto - gasto.monto_pagado, metodo_pago: '', referencia_pago: '' })
    setShowPagoModal(true)
  }

  const formatMoney = (amount: number) => new Intl.NumberFormat('es-CL', { style: 'currency', currency: 'CLP' }).format(amount)

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Gastos Comunes</h1>
          <p className="text-gray-500 mt-1">Gestion de periodos y pagos</p>
        </div>
        <Button onClick={() => setShowNewPeriodoModal(true)} icon={<Plus size={20} />}>Nuevo Periodo</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Periodos sidebar */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="px-4 py-3 bg-gray-50 border-b">
              <h2 className="font-medium text-gray-900">Periodos</h2>
            </div>
            {isLoading ? (
              <div className="p-8 text-center"><Loader2 className="h-6 w-6 animate-spin text-blue-600 mx-auto" /></div>
            ) : (
              <div className="divide-y divide-gray-100 max-h-96 overflow-y-auto">
                {periodos.map((periodo) => (
                  <button
                    key={periodo.id}
                    onClick={() => selectPeriodo(periodo)}
                    className={`w-full px-4 py-3 text-left hover:bg-gray-50 transition-colors ${selectedPeriodo?.id === periodo.id ? 'bg-blue-50 border-l-2 border-blue-600' : ''}`}
                  >
                    <p className="font-medium text-gray-900">{MESES[periodo.month - 1]} {periodo.year}</p>
                    <p className="text-sm text-gray-500">{formatMoney(periodo.monto_base)}</p>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Main content */}
        <div className="lg:col-span-3 space-y-6">
          {selectedPeriodo && resumen && (
            <>
              {/* Resumen cards */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-white rounded-xl shadow-sm border p-4">
                  <p className="text-sm text-gray-500">Total Parcelas</p>
                  <p className="text-2xl font-bold">{resumen.total_parcelas}</p>
                </div>
                <div className="bg-white rounded-xl shadow-sm border p-4">
                  <p className="text-sm text-gray-500">Pagados</p>
                  <p className="text-2xl font-bold text-green-600">{resumen.total_pagados}</p>
                </div>
                <div className="bg-white rounded-xl shadow-sm border p-4">
                  <p className="text-sm text-gray-500">Pendientes</p>
                  <p className="text-2xl font-bold text-yellow-600">{resumen.total_pendientes}</p>
                </div>
                <div className="bg-white rounded-xl shadow-sm border p-4">
                  <p className="text-sm text-gray-500">Recaudado</p>
                  <p className="text-2xl font-bold text-blue-600">{resumen.porcentaje_recaudo.toFixed(0)}%</p>
                </div>
              </div>

              {/* Gastos table */}
              <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
                {isLoadingGastos ? (
                  <div className="p-8 text-center"><Loader2 className="h-6 w-6 animate-spin text-blue-600 mx-auto" /></div>
                ) : (
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Parcela</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Propietario</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Monto</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Acciones</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {gastos.map((gasto) => (
                        <tr key={gasto.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap font-medium">{gasto.parcela_numero}</td>
                          <td className="px-6 py-4 whitespace-nowrap text-gray-500">{gasto.user_name || '-'}</td>
                          <td className="px-6 py-4 whitespace-nowrap">{formatMoney(gasto.monto)}</td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <Badge color={PAGO_STATUS_COLORS[gasto.status]}>{PAGO_STATUS_LABELS[gasto.status]}</Badge>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-right">
                            {gasto.status !== 'paid' && (
                              <button onClick={() => openPagoModal(gasto)} className="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg" title="Registrar pago">
                                <DollarSign size={18} />
                              </button>
                            )}
                            {gasto.status === 'paid' && <CheckCircle className="h-5 w-5 text-green-500 inline" />}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            </>
          )}
        </div>
      </div>

      {/* New Periodo Modal */}
      <Modal open={showNewPeriodoModal} onClose={() => setShowNewPeriodoModal(false)} title="Nuevo Periodo" size="md">
        <form onSubmit={handleCreatePeriodo} className="space-y-6">
          <div className="grid grid-cols-2 gap-4">
            <Input type="number" label="Ano" value={newPeriodo.year} onChange={(e) => setNewPeriodo({ ...newPeriodo, year: parseInt(e.target.value) })} required min={2020} max={2030} />
            <Input type="number" label="Mes" value={newPeriodo.month} onChange={(e) => setNewPeriodo({ ...newPeriodo, month: parseInt(e.target.value) })} required min={1} max={12} />
          </div>
          <Input type="number" label="Monto Base" value={newPeriodo.monto_base} onChange={(e) => setNewPeriodo({ ...newPeriodo, monto_base: parseInt(e.target.value) })} required />
          <Input type="date" label="Fecha Vencimiento" value={newPeriodo.fecha_vencimiento} onChange={(e) => setNewPeriodo({ ...newPeriodo, fecha_vencimiento: e.target.value })} required />
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowNewPeriodoModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">Crear Periodo</Button>
          </div>
        </form>
      </Modal>

      {/* Pago Modal */}
      <Modal open={showPagoModal} onClose={() => setShowPagoModal(false)} title="Registrar Pago" size="md">
        <form onSubmit={handleRegistrarPago} className="space-y-6">
          <p className="text-gray-600">Parcela: <strong>{selectedGasto?.parcela_numero}</strong></p>
          <Input type="number" label="Monto a Pagar" value={pagoData.monto_pagado} onChange={(e) => setPagoData({ ...pagoData, monto_pagado: parseInt(e.target.value) })} required />
          <Input label="Metodo de Pago" value={pagoData.metodo_pago} onChange={(e) => setPagoData({ ...pagoData, metodo_pago: e.target.value })} placeholder="Transferencia, Efectivo, etc." />
          <Input label="Referencia" value={pagoData.referencia_pago} onChange={(e) => setPagoData({ ...pagoData, referencia_pago: e.target.value })} placeholder="Numero de operacion" />
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowPagoModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">Registrar Pago</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
