package cl.vinapelvin.condominio.ui.home

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material.icons.outlined.ExitToApp
import androidx.compose.foundation.BorderStroke
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.ui.navigation.NavRoutes
import cl.vinapelvin.condominio.ui.theme.*

data class MenuItem(
    val title: String,
    val icon: ImageVector,
    val route: String,
    val color: Color,
    val badge: Int? = null
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(
    viewModel: HomeViewModel = hiltViewModel(),
    onNavigate: (String) -> Unit,
    onLogout: () -> Unit
) {
    val userName by viewModel.userName.collectAsState(initial = "Usuario")
    val userRole by viewModel.userRole.collectAsState(initial = "")
    val userParcelaId by viewModel.userParcelaId.collectAsState(initial = null)
    val notificationCount by viewModel.notificationCount.collectAsState()

    val role = userRole.orEmpty()
    val isVecinoODirectiva = role == "vecino" || role == "directiva" || role == "admin"

    val publicItems = listOf(
        MenuItem("Comunicados", Icons.Default.Campaign, NavRoutes.Comunicados.route, Blue600),
        MenuItem("Eventos", Icons.Default.Event, NavRoutes.Eventos.route, Green600),
        MenuItem("Emergencias", Icons.Default.Warning, NavRoutes.Emergencias.route, Red600),
        MenuItem("Galerías", Icons.Default.PhotoLibrary, NavRoutes.Galerias.route, Amber500),
        MenuItem("Mapa", Icons.Default.Map, NavRoutes.Mapa.route, Green500),
    )

    val protectedItems = listOf(
        MenuItem("Votaciones", Icons.Default.HowToVote, NavRoutes.Votaciones.route, Amber500),
        // Gastos requiere parcela asociada
        MenuItem("Gastos", Icons.Default.Receipt, NavRoutes.Gastos.route, Blue700),
        MenuItem("Tesorería", Icons.Default.AccountBalance, NavRoutes.Tesoreria.route, Green600),
        MenuItem("Actas", Icons.Default.Description, NavRoutes.Actas.route, Blue800),
        MenuItem("Documentos", Icons.Default.Folder, NavRoutes.Documentos.route, Gray500),
        MenuItem("Notificaciones", Icons.Default.Notifications, NavRoutes.Notificaciones.route, Gray500, notificationCount),
        MenuItem("Contacto Directiva", Icons.Default.Mail, NavRoutes.Contacto.route, Tierra),
    )

    val filteredProtectedItems =
        if (userParcelaId == null) protectedItems.filterNot { it.route == NavRoutes.Gastos.route }
        else protectedItems

    val menuItems = if (isVecinoODirectiva) publicItems + filteredProtectedItems else publicItems

    Scaffold(
        containerColor = MaterialTheme.colorScheme.background,
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text("Vina Pelvin", fontWeight = FontWeight.Bold)
                        Text(
                            text = "Hola, $userName",
                            fontSize = 14.sp,
                            color = MaterialTheme.colorScheme.onPrimaryContainer.copy(alpha = 0.7f)
                        )
                    }
                },
                actions = {
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Outlined.ExitToApp, contentDescription = "Cerrar sesion")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primaryContainer,
                    titleContentColor = MaterialTheme.colorScheme.onPrimaryContainer
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            if (!isVecinoODirectiva && role.isNotBlank()) {
                Card(
                    modifier = Modifier
                        .padding(horizontal = 16.dp)
                        .padding(top = 16.dp),
                    colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
                    border = BorderStroke(1.dp, Gray100),
                    elevation = CardDefaults.cardElevation(defaultElevation = 0.dp)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text(
                            text = "Acceso limitado",
                            fontWeight = FontWeight.SemiBold,
                            color = Gray900
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(
                            text = "Esta app móvil está pensada para vecinos y directiva. Tu rol actual: $role.",
                            color = Gray500,
                            fontSize = 13.sp
                        )
                    }
                }
                Spacer(modifier = Modifier.height(8.dp))
            }

            LazyVerticalGrid(
                columns = GridCells.Fixed(2),
                modifier = Modifier
                    .fillMaxSize()
                    .padding(16.dp),
                horizontalArrangement = Arrangement.spacedBy(16.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                items(menuItems) { item ->
                    MenuCard(
                        item = item,
                        onClick = { onNavigate(item.route) }
                    )
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MenuCard(
    item: MenuItem,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .aspectRatio(1f)
            .clickable(onClick = onClick),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        border = BorderStroke(1.dp, Gray100),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Box(modifier = Modifier.fillMaxSize()) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                Icon(
                    imageVector = item.icon,
                    contentDescription = item.title,
                    modifier = Modifier.size(44.dp),
                    tint = item.color.copy(alpha = 0.9f)
                )
                Spacer(modifier = Modifier.height(12.dp))
                Text(
                    text = item.title,
                    fontWeight = FontWeight.Medium,
                    fontSize = 14.sp,
                    color = Gray900
                )
            }

            // Badge
            item.badge?.let { count ->
                if (count > 0) {
                    Badge(
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .padding(8.dp),
                        containerColor = Red600
                    ) {
                        Text(
                            text = if (count > 99) "99+" else count.toString(),
                            color = Color.White
                        )
                    }
                }
            }
        }
    }
}
