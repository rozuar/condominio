package cl.vinapelvin.condominio.ui.navigation

sealed class NavRoutes(val route: String) {
    object Login : NavRoutes("login")
    object Home : NavRoutes("home")
    object Comunicados : NavRoutes("comunicados")
    object ComunicadoDetail : NavRoutes("comunicado/{id}") {
        fun createRoute(id: String) = "comunicado/$id"
    }
    object Eventos : NavRoutes("eventos")
    object EventoDetail : NavRoutes("evento/{id}") {
        fun createRoute(id: String) = "evento/$id"
    }
    object Emergencias : NavRoutes("emergencias")
    object Votaciones : NavRoutes("votaciones")
    object VotacionDetail : NavRoutes("votacion/{id}") {
        fun createRoute(id: String) = "votacion/$id"
    }
    object Gastos : NavRoutes("gastos")
    object Tesoreria : NavRoutes("tesoreria")
    object Actas : NavRoutes("actas")
    object ActaDetail : NavRoutes("acta/{id}") {
        fun createRoute(id: String) = "acta/$id"
    }
    object Documentos : NavRoutes("documentos")
    object Notificaciones : NavRoutes("notificaciones")
    object Contacto : NavRoutes("contacto")
    object Galerias : NavRoutes("galerias")
    object GaleriaDetail : NavRoutes("galeria/{id}") {
        fun createRoute(id: String) = "galeria/$id"
    }
    object Mapa : NavRoutes("mapa")
}
