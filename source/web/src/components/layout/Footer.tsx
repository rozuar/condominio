export default function Footer() {
  return (
    <footer className="bg-gray-100 border-t">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <h3 className="font-bold text-primary mb-2">Comunidad Viña Pelvin</h3>
            <p className="text-sm text-gray-600">
              73 parcelas en armonía con la naturaleza
            </p>
          </div>
          <div>
            <h4 className="font-semibold text-gray-800 mb-2">Enlaces</h4>
            <ul className="space-y-1 text-sm text-gray-600">
              <li><a href="/comunicados" className="hover:text-primary">Comunicados</a></li>
              <li><a href="/calendario" className="hover:text-primary">Calendario</a></li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-gray-800 mb-2">Contacto</h4>
            <p className="text-sm text-gray-600">
              directiva@vinapelvin.cl
            </p>
          </div>
        </div>
        <div className="border-t mt-8 pt-6 text-center text-sm text-gray-500">
          © {new Date().getFullYear()} Comunidad Viña Pelvin. Todos los derechos reservados.
        </div>
      </div>
    </footer>
  )
}
