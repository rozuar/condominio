import Link from 'next/link'

export default function RegistroPage() {
  return (
    <div className="min-h-[80vh] flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Registro deshabilitado</h1>
          <p className="text-gray-600 mt-2">
            Las cuentas deben ser creadas por la directiva desde el backoffice.
          </p>
        </div>

        <div className="bg-white p-8 rounded-lg shadow border">
          <p className="text-sm text-gray-700">
            Si necesitas acceso, escribe a la directiva desde{' '}
            <Link href="/contacto" className="text-primary font-medium hover:underline">
              Contacto
            </Link>
            .
          </p>

          <p className="mt-6 text-center text-sm text-gray-600">
            ¿Ya tienes cuenta?{' '}
            <Link href="/auth/login" className="text-primary font-medium hover:underline">
              Inicia sesion
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
