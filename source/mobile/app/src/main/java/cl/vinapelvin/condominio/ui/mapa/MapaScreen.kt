package cl.vinapelvin.condominio.ui.mapa

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.ExperimentalMaterialApi
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material.pullrefresh.PullRefreshIndicator
import androidx.compose.material.pullrefresh.pullRefresh
import androidx.compose.material.pullrefresh.rememberPullRefreshState
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.MapaArea
import cl.vinapelvin.condominio.data.model.MapaPunto
import cl.vinapelvin.condominio.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class, ExperimentalMaterialApi::class)
@Composable
fun MapaScreen(
    onBack: () -> Unit,
    viewModel: MapaViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val pullRefreshState = rememberPullRefreshState(
        refreshing = uiState.isRefreshing,
        onRefresh = { viewModel.refresh() }
    )

    var selectedTab by remember { mutableIntStateOf(0) }
    val tabs = listOf("Puntos de Interés", "Áreas")

    Scaffold(
        containerColor = MaterialTheme.colorScheme.background,
        topBar = {
            TopAppBar(
                title = { Text("Mapa del Condominio") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Volver")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primaryContainer,
                    titleContentColor = MaterialTheme.colorScheme.onPrimaryContainer,
                    navigationIconContentColor = MaterialTheme.colorScheme.onPrimaryContainer
                )
            )
        }
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .pullRefresh(pullRefreshState)
        ) {
            when {
                uiState.isLoading -> {
                    CircularProgressIndicator(
                        modifier = Modifier.align(Alignment.Center)
                    )
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(16.dp),
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.Center
                    ) {
                        Text(
                            text = uiState.error ?: "Error desconocido",
                            color = MaterialTheme.colorScheme.error
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Button(onClick = { viewModel.loadMapaData() }) {
                            Text("Reintentar")
                        }
                    }
                }
                else -> {
                    Column(modifier = Modifier.fillMaxSize()) {
                        // Tab Row
                        TabRow(
                            selectedTabIndex = selectedTab,
                            containerColor = MaterialTheme.colorScheme.surface
                        ) {
                            tabs.forEachIndexed { index, title ->
                                Tab(
                                    selected = selectedTab == index,
                                    onClick = { selectedTab = index },
                                    text = { Text(title) }
                                )
                            }
                        }

                        // Content
                        when (selectedTab) {
                            0 -> PuntosTab(
                                puntos = uiState.puntos,
                                onPuntoClick = { viewModel.selectPunto(it) }
                            )
                            1 -> AreasTab(
                                areas = uiState.areas,
                                onAreaClick = { viewModel.selectArea(it) }
                            )
                        }
                    }

                    // Detail dialogs
                    uiState.selectedPunto?.let { punto ->
                        PuntoDetailDialog(
                            punto = punto,
                            onDismiss = { viewModel.clearSelection() }
                        )
                    }

                    uiState.selectedArea?.let { area ->
                        AreaDetailDialog(
                            area = area,
                            onDismiss = { viewModel.clearSelection() }
                        )
                    }
                }
            }

            PullRefreshIndicator(
                refreshing = uiState.isRefreshing,
                state = pullRefreshState,
                modifier = Modifier.align(Alignment.TopCenter)
            )
        }
    }
}

@Composable
private fun PuntosTab(
    puntos: List<MapaPunto>,
    onPuntoClick: (MapaPunto) -> Unit
) {
    if (puntos.isEmpty()) {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(16.dp),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(
                    imageVector = Icons.Default.LocationOff,
                    contentDescription = null,
                    modifier = Modifier.size(64.dp),
                    tint = Gray400
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "No hay puntos de interés",
                    color = Gray500
                )
            }
        }
    } else {
        LazyColumn(
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            items(puntos, key = { it.id }) { punto ->
                PuntoCard(
                    punto = punto,
                    onClick = { onPuntoClick(punto) }
                )
            }
        }
    }
}

@Composable
private fun AreasTab(
    areas: List<MapaArea>,
    onAreaClick: (MapaArea) -> Unit
) {
    if (areas.isEmpty()) {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(16.dp),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(
                    imageVector = Icons.Default.Layers,
                    contentDescription = null,
                    modifier = Modifier.size(64.dp),
                    tint = Gray400
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = "No hay áreas definidas",
                    color = Gray500
                )
            }
        }
    } else {
        // Group areas by type
        val groupedAreas = areas.groupBy { it.type }

        LazyColumn(
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            groupedAreas.forEach { (type, areasInGroup) ->
                item {
                    Text(
                        text = getAreaTypeLabel(type),
                        fontWeight = FontWeight.Bold,
                        fontSize = 16.sp,
                        color = Gray900,
                        modifier = Modifier.padding(vertical = 8.dp)
                    )
                }
                items(areasInGroup, key = { it.id }) { area ->
                    AreaCard(
                        area = area,
                        onClick = { onAreaClick(area) }
                    )
                }
            }
        }
    }
}

@Composable
private fun PuntoCard(
    punto: MapaPunto,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .clip(CircleShape)
                    .background(getIconColorForType(punto.type).copy(alpha = 0.1f)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = getIconForType(punto.type),
                    contentDescription = null,
                    tint = getIconColorForType(punto.type),
                    modifier = Modifier.size(24.dp)
                )
            }

            Spacer(modifier = Modifier.width(16.dp))

            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = punto.name,
                    fontWeight = FontWeight.SemiBold,
                    color = Gray900
                )
                punto.description?.takeIf { it.isNotBlank() }?.let { desc ->
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        text = desc,
                        color = Gray500,
                        fontSize = 13.sp
                    )
                }
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = punto.type.replaceFirstChar { it.uppercase() },
                    color = Gray400,
                    fontSize = 12.sp
                )
            }

            Icon(
                imageVector = Icons.Default.ChevronRight,
                contentDescription = null,
                tint = Gray400
            )
        }
    }
}

@Composable
private fun AreaCard(
    area: MapaArea,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(16.dp)
                    .clip(RoundedCornerShape(4.dp))
                    .background(parseColor(area.fillColor))
            )

            Spacer(modifier = Modifier.width(16.dp))

            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = area.name,
                    fontWeight = FontWeight.SemiBold,
                    color = Gray900
                )
                area.description?.takeIf { it.isNotBlank() }?.let { desc ->
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        text = desc,
                        color = Gray500,
                        fontSize = 13.sp
                    )
                }
            }

            Icon(
                imageVector = Icons.Default.ChevronRight,
                contentDescription = null,
                tint = Gray400
            )
        }
    }
}

@Composable
private fun PuntoDetailDialog(
    punto: MapaPunto,
    onDismiss: () -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        icon = {
            Icon(
                imageVector = getIconForType(punto.type),
                contentDescription = null,
                tint = getIconColorForType(punto.type),
                modifier = Modifier.size(32.dp)
            )
        },
        title = { Text(punto.name) },
        text = {
            Column {
                punto.description?.takeIf { it.isNotBlank() }?.let { desc ->
                    Text(desc)
                    Spacer(modifier = Modifier.height(16.dp))
                }
                Row {
                    Icon(
                        imageVector = Icons.Default.LocationOn,
                        contentDescription = null,
                        tint = Gray500,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(
                        text = "Lat: ${String.format("%.6f", punto.lat)}, Lng: ${String.format("%.6f", punto.lng)}",
                        color = Gray500,
                        fontSize = 13.sp
                    )
                }
                Spacer(modifier = Modifier.height(8.dp))
                Row {
                    Icon(
                        imageVector = Icons.Default.Category,
                        contentDescription = null,
                        tint = Gray500,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(
                        text = "Tipo: ${punto.type.replaceFirstChar { it.uppercase() }}",
                        color = Gray500,
                        fontSize = 13.sp
                    )
                }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) {
                Text("Cerrar")
            }
        }
    )
}

@Composable
private fun AreaDetailDialog(
    area: MapaArea,
    onDismiss: () -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        icon = {
            Box(
                modifier = Modifier
                    .size(32.dp)
                    .clip(RoundedCornerShape(8.dp))
                    .background(parseColor(area.fillColor))
            )
        },
        title = { Text(area.name) },
        text = {
            Column {
                area.description?.takeIf { it.isNotBlank() }?.let { desc ->
                    Text(desc)
                    Spacer(modifier = Modifier.height(16.dp))
                }
                Row {
                    Icon(
                        imageVector = Icons.Default.Layers,
                        contentDescription = null,
                        tint = Gray500,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(
                        text = "Tipo: ${getAreaTypeLabel(area.type)}",
                        color = Gray500,
                        fontSize = 13.sp
                    )
                }
                area.centerLat?.let { lat ->
                    area.centerLng?.let { lng ->
                        Spacer(modifier = Modifier.height(8.dp))
                        Row {
                            Icon(
                                imageVector = Icons.Default.LocationOn,
                                contentDescription = null,
                                tint = Gray500,
                                modifier = Modifier.size(16.dp)
                            )
                            Spacer(modifier = Modifier.width(8.dp))
                            Text(
                                text = "Centro: ${String.format("%.6f", lat)}, ${String.format("%.6f", lng)}",
                                color = Gray500,
                                fontSize = 13.sp
                            )
                        }
                    }
                }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) {
                Text("Cerrar")
            }
        }
    )
}

// Helper functions

private fun getIconForType(type: String): ImageVector {
    return when (type.lowercase()) {
        "entrada", "acceso" -> Icons.Default.MeetingRoom
        "estacionamiento", "parking" -> Icons.Default.LocalParking
        "piscina" -> Icons.Default.Pool
        "cancha", "deporte" -> Icons.Default.SportsSoccer
        "area_verde", "parque" -> Icons.Default.Park
        "seguridad", "guardia" -> Icons.Default.Security
        "porteria" -> Icons.Default.Home
        "basura", "reciclaje" -> Icons.Default.Delete
        "agua", "riego" -> Icons.Default.Water
        "emergencia" -> Icons.Default.Warning
        else -> Icons.Default.Place
    }
}

private fun getIconColorForType(type: String): Color {
    return when (type.lowercase()) {
        "entrada", "acceso" -> Blue600
        "estacionamiento", "parking" -> Gray600
        "piscina" -> Blue500
        "cancha", "deporte" -> Green600
        "area_verde", "parque" -> Green500
        "seguridad", "guardia" -> Red600
        "porteria" -> Amber500
        "basura", "reciclaje" -> Gray500
        "agua", "riego" -> Blue400
        "emergencia" -> Red500
        else -> Gray600
    }
}

private fun getAreaTypeLabel(type: String): String {
    return when (type.lowercase()) {
        "parcela" -> "Parcelas"
        "area_comun" -> "Áreas Comunes"
        "acceso" -> "Accesos"
        "canal" -> "Canales"
        "camino" -> "Caminos"
        else -> type.replaceFirstChar { it.uppercase() }
    }
}

private fun parseColor(colorString: String): Color {
    return try {
        Color(android.graphics.Color.parseColor(colorString))
    } catch (e: Exception) {
        Gray400
    }
}
