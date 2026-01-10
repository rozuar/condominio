package cl.vinapelvin.condominio.ui.navigation

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import cl.vinapelvin.condominio.ui.auth.LoginScreen
import cl.vinapelvin.condominio.ui.auth.LoginViewModel
import cl.vinapelvin.condominio.ui.home.HomeScreen
import cl.vinapelvin.condominio.ui.comunicados.ComunicadosScreen
import cl.vinapelvin.condominio.ui.comunicados.ComunicadoDetailScreen
import cl.vinapelvin.condominio.ui.eventos.EventosScreen
import cl.vinapelvin.condominio.ui.eventos.EventoDetailScreen
import cl.vinapelvin.condominio.ui.emergencias.EmergenciasScreen
import cl.vinapelvin.condominio.ui.votaciones.VotacionesScreen
import cl.vinapelvin.condominio.ui.votaciones.VotacionDetailScreen
import cl.vinapelvin.condominio.ui.gastos.GastosScreen
import cl.vinapelvin.condominio.ui.notificaciones.NotificacionesScreen
import cl.vinapelvin.condominio.ui.contacto.ContactoScreen
import cl.vinapelvin.condominio.ui.tesoreria.TesoreriaScreen
import cl.vinapelvin.condominio.ui.actas.ActasScreen
import cl.vinapelvin.condominio.ui.actas.ActaDetailScreen
import cl.vinapelvin.condominio.ui.documentos.DocumentosScreen

@Composable
fun AppNavigation() {
    val navController = rememberNavController()
    val loginViewModel: LoginViewModel = hiltViewModel()
    val isLoggedIn by loginViewModel.isLoggedIn.collectAsState(initial = false)

    val startDestination = if (isLoggedIn) NavRoutes.Home.route else NavRoutes.Login.route

    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable(NavRoutes.Login.route) {
            LoginScreen(
                viewModel = loginViewModel,
                onLoginSuccess = {
                    navController.navigate(NavRoutes.Home.route) {
                        popUpTo(NavRoutes.Login.route) { inclusive = true }
                    }
                }
            )
        }

        composable(NavRoutes.Home.route) {
            HomeScreen(
                onNavigate = { route -> navController.navigate(route) },
                onLogout = {
                    loginViewModel.logout()
                    navController.navigate(NavRoutes.Login.route) {
                        popUpTo(0) { inclusive = true }
                    }
                }
            )
        }

        composable(NavRoutes.Comunicados.route) {
            ComunicadosScreen(
                onBack = { navController.popBackStack() },
                onComunicadoClick = { id ->
                    navController.navigate(NavRoutes.ComunicadoDetail.createRoute(id))
                }
            )
        }

        composable(
            route = NavRoutes.ComunicadoDetail.route,
            arguments = listOf(navArgument("id") { type = NavType.StringType })
        ) { backStackEntry ->
            val id = backStackEntry.arguments?.getString("id") ?: return@composable
            ComunicadoDetailScreen(
                comunicadoId = id,
                onBack = { navController.popBackStack() }
            )
        }

        composable(NavRoutes.Eventos.route) {
            EventosScreen(
                onBack = { navController.popBackStack() },
                onEventoClick = { id ->
                    navController.navigate(NavRoutes.EventoDetail.createRoute(id))
                }
            )
        }

        composable(
            route = NavRoutes.EventoDetail.route,
            arguments = listOf(navArgument("id") { type = NavType.StringType })
        ) { backStackEntry ->
            val id = backStackEntry.arguments?.getString("id") ?: return@composable
            EventoDetailScreen(
                eventoId = id,
                onBack = { navController.popBackStack() }
            )
        }

        composable(NavRoutes.Emergencias.route) {
            EmergenciasScreen(onBack = { navController.popBackStack() })
        }

        composable(NavRoutes.Votaciones.route) {
            VotacionesScreen(
                onBack = { navController.popBackStack() },
                onVotacionClick = { id ->
                    navController.navigate(NavRoutes.VotacionDetail.createRoute(id))
                }
            )
        }

        composable(
            route = NavRoutes.VotacionDetail.route,
            arguments = listOf(navArgument("id") { type = NavType.StringType })
        ) { backStackEntry ->
            val id = backStackEntry.arguments?.getString("id") ?: return@composable
            VotacionDetailScreen(
                votacionId = id,
                onBack = { navController.popBackStack() }
            )
        }

        composable(NavRoutes.Gastos.route) {
            GastosScreen(onBack = { navController.popBackStack() })
        }

        composable(NavRoutes.Tesoreria.route) {
            TesoreriaScreen(onBack = { navController.popBackStack() })
        }

        composable(NavRoutes.Actas.route) {
            ActasScreen(
                onBack = { navController.popBackStack() },
                onActaClick = { id ->
                    navController.navigate(NavRoutes.ActaDetail.createRoute(id))
                }
            )
        }

        composable(
            route = NavRoutes.ActaDetail.route,
            arguments = listOf(navArgument("id") { type = NavType.StringType })
        ) { backStackEntry ->
            val id = backStackEntry.arguments?.getString("id") ?: return@composable
            ActaDetailScreen(
                actaId = id,
                onBack = { navController.popBackStack() }
            )
        }

        composable(NavRoutes.Documentos.route) {
            DocumentosScreen(onBack = { navController.popBackStack() })
        }

        composable(NavRoutes.Notificaciones.route) {
            NotificacionesScreen(onBack = { navController.popBackStack() })
        }

        composable(NavRoutes.Contacto.route) {
            ContactoScreen(onBack = { navController.popBackStack() })
        }
    }
}
