import Link from 'next/link'
import { ArrowRight, FileText, Calendar, Shield, CreditCard, Mail } from 'lucide-react'
import { getLatestComunicados, getUpcomingEventos } from '@/lib/api'
import ComunicadoCard from '@/components/comunicados/ComunicadoCard'
import EventoCard from '@/components/calendario/EventoCard'

export const revalidate = 60

async function getData() {
  try {
    const [comunicadosRes, eventosRes] = await Promise.all([
      getLatestComunicados(3),
      getUpcomingEventos(3),
    ])
    return {
      comunicados: comunicadosRes.comunicados || [],
      eventos: eventosRes.eventos || [],
    }
  } catch (error) {
    return { comunicados: [], eventos: [] }
  }
}

const quickLinks = [
  { name: 'Comunicados', href: '/comunicados', icon: FileText, description: 'Últimas noticias' },
  { name: 'Calendario', href: '/calendario', icon: Calendar, description: 'Próximos eventos' },
  { name: 'Emergencias', href: '/emergencias', icon: Shield, description: 'Avisos urgentes' },
  { name: 'Gastos Comunes', href: '/gastos', icon: CreditCard, description: 'Pagar en línea' },
  { name: 'Contacto Directiva', href: '/contacto', icon: Mail, description: 'Enviar mensaje' },
]

export default async function HomePage() {
  const { comunicados, eventos } = await getData()

  return (
    <div>
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-gray-950 via-gray-900 to-primary text-white py-16 md:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Comunidad Viña Pelvin
            </h1>
            <p className="text-xl text-white/90 max-w-2xl mx-auto">
              73 parcelas en armonía con la naturaleza. Portal oficial de nuestra comunidad.
            </p>
          </div>
        </div>
      </section>

      {/* Quick Links */}
      <section className="py-10">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {quickLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className="group flex flex-col items-center p-4 card hover:shadow-md transition-shadow text-center"
              >
                <link.icon className="w-8 h-8 text-primary mb-2 group-hover:text-primary-dark transition-colors" />
                <span className="font-medium text-gray-900">{link.name}</span>
                <span className="text-xs text-gray-500">{link.description}</span>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Latest Comunicados */}
      <section className="py-12">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Últimos Comunicados</h2>
            <Link href="/comunicados" className="text-primary font-medium flex items-center gap-1 hover:underline">
              Ver todos <ArrowRight size={16} />
            </Link>
          </div>

          {comunicados.length > 0 ? (
            <div className="grid md:grid-cols-3 gap-6">
              {comunicados.map((comunicado) => (
                <ComunicadoCard key={comunicado.id} comunicado={comunicado} />
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">No hay comunicados recientes</p>
          )}
        </div>
      </section>

      {/* Upcoming Events */}
      <section className="py-12 bg-gray-50">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Próximos Eventos</h2>
            <Link href="/calendario" className="text-primary font-medium flex items-center gap-1 hover:underline">
              Ver calendario <ArrowRight size={16} />
            </Link>
          </div>

          {eventos.length > 0 ? (
            <div className="grid md:grid-cols-3 gap-4">
              {eventos.map((evento) => (
                <EventoCard key={evento.id} evento={evento} />
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">No hay eventos próximos</p>
          )}
        </div>
      </section>
    </div>
  )
}
