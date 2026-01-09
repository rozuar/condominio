'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { getTesoreriaResumen, getMovimientos, createMovimiento } from '@/lib/api'
import type { TesoreriaResumen, Movimiento, MovimientoType } from '@/types'
import { MOVIMIENTO_TYPE_LABELS, MOVIMIENTO_TYPE_COLORS } from '@/types'
import { Plus, Loader2, TrendingUp, TrendingDown, Wallet } from 'lucide-react'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Modal from '@/components/ui/Modal'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'

export default function TesoreriaPage() {
  const { getToken } = useAuth()
  const [resumen, setResumen] = useState<TesoreriaResumen | null>(null)
  const [movimientos, setMovimientos] = useState<Movimiento[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [formData, setFormData] = useState({ description: '', amount: 0, type: 'ingreso' as MovimientoType, category: '', date: '' })

  const fetchData = async () => {
    const token = getToken()
    if (!token) return

    setIsLoading(true)
    try {
      const [resumenData, movimientosData] = await Promise.all([
        getTesoreriaResumen(token),
        getMovimientos(token, { page, per_page: 10 }),
      ])
      setResumen(resumenData)
      setMovimientos(movimientosData.movimientos)
      setTotal(movimientosData.total)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [page])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = getToken()
    if (!token) return

    setIsSaving(true)
    try {
      await createMovimiento(token, formData)
      setShowModal(false)
      setFormData({ description: '', amount: 0, type: 'ingreso', category: '', date: '' })
      fetchData()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al crear')
    } finally {
      setIsSaving(false)
    }
  }

  const formatMoney = (amount: number) => new Intl.NumberFormat('es-CL', { style: 'currency', currency: 'CLP' }).format(amount)
  const totalPages = Math.ceil(total / 10)

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Tesoreria</h1>
          <p className="text-gray-500 mt-1">Gestion de ingresos y egresos</p>
        </div>
        <Button onClick={() => setShowModal(true)} icon={<Plus size={20} />}>Nuevo Movimiento</Button>
      </div>

      {error && <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">{error}</div>}

      {isLoading ? (
        <div className="flex items-center justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-blue-600" /></div>
      ) : (
        <>
          {resumen && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
              <div className="bg-white rounded-xl shadow-sm border p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-500">Ingresos</p>
                    <p className="text-2xl font-bold text-green-600">{formatMoney(resumen.total_ingresos)}</p>
                  </div>
                  <div className="p-3 bg-green-100 rounded-full"><TrendingUp className="h-6 w-6 text-green-600" /></div>
                </div>
              </div>
              <div className="bg-white rounded-xl shadow-sm border p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-500">Egresos</p>
                    <p className="text-2xl font-bold text-red-600">{formatMoney(resumen.total_egresos)}</p>
                  </div>
                  <div className="p-3 bg-red-100 rounded-full"><TrendingDown className="h-6 w-6 text-red-600" /></div>
                </div>
              </div>
              <div className="bg-white rounded-xl shadow-sm border p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-500">Balance</p>
                    <p className={`text-2xl font-bold ${resumen.balance >= 0 ? 'text-blue-600' : 'text-red-600'}`}>{formatMoney(resumen.balance)}</p>
                  </div>
                  <div className="p-3 bg-blue-100 rounded-full"><Wallet className="h-6 w-6 text-blue-600" /></div>
                </div>
              </div>
            </div>
          )}

          <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Descripcion</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Categoria</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Monto</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {movimientos.map((mov) => (
                  <tr key={mov.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{format(new Date(mov.date), "d MMM yyyy", { locale: es })}</td>
                    <td className="px-6 py-4 text-sm text-gray-900">{mov.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{mov.category}</td>
                    <td className="px-6 py-4 whitespace-nowrap"><Badge color={MOVIMIENTO_TYPE_COLORS[mov.type]}>{MOVIMIENTO_TYPE_LABELS[mov.type]}</Badge></td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium text-right ${mov.type === 'ingreso' ? 'text-green-600' : 'text-red-600'}`}>
                      {mov.type === 'ingreso' ? '+' : '-'}{formatMoney(mov.amount)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {totalPages > 1 && (
              <div className="px-6 py-4 border-t flex items-center justify-between">
                <p className="text-sm text-gray-500">Pagina {page} de {totalPages}</p>
                <div className="flex gap-2">
                  <Button variant="secondary" size="sm" onClick={() => setPage(page - 1)} disabled={page === 1}>Anterior</Button>
                  <Button variant="secondary" size="sm" onClick={() => setPage(page + 1)} disabled={page === totalPages}>Siguiente</Button>
                </div>
              </div>
            )}
          </div>
        </>
      )}

      <Modal open={showModal} onClose={() => setShowModal(false)} title="Nuevo Movimiento">
        <form onSubmit={handleSubmit} className="space-y-6">
          <Select label="Tipo" value={formData.type} onChange={(e) => setFormData({ ...formData, type: e.target.value as MovimientoType })} options={[{ value: 'ingreso', label: 'Ingreso' }, { value: 'egreso', label: 'Egreso' }]} />
          <Input label="Descripcion" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} required />
          <Input type="number" label="Monto" value={formData.amount} onChange={(e) => setFormData({ ...formData, amount: parseInt(e.target.value) })} required />
          <Input label="Categoria" value={formData.category} onChange={(e) => setFormData({ ...formData, category: e.target.value })} placeholder="Ej: Mantenimiento, Servicios" />
          <Input type="date" label="Fecha" value={formData.date} onChange={(e) => setFormData({ ...formData, date: e.target.value })} required />
          <div className="flex gap-3">
            <Button type="button" variant="secondary" onClick={() => setShowModal(false)} className="flex-1">Cancelar</Button>
            <Button type="submit" loading={isSaving} className="flex-1">Crear</Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
