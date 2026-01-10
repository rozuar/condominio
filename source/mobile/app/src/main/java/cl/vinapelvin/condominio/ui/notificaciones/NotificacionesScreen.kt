package cl.vinapelvin.condominio.ui.notificaciones

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Circle
import androidx.compose.material.icons.filled.DoneAll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.Notificacion
import cl.vinapelvin.condominio.ui.theme.*
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificacionesScreen(
    viewModel: NotificacionesViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Notificaciones") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Volver")
                    }
                },
                actions = {
                    if (uiState.notificaciones.any { !it.isRead }) {
                        IconButton(onClick = { viewModel.markAllAsRead() }) {
                            Icon(Icons.Default.DoneAll, contentDescription = "Marcar todas como leidas")
                        }
                    }
                }
            )
        }
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            when {
                uiState.isLoading -> {
                    CircularProgressIndicator(modifier = Modifier.align(Alignment.Center))
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(uiState.error!!, color = MaterialTheme.colorScheme.error)
                        Spacer(modifier = Modifier.height(16.dp))
                        Button(onClick = { viewModel.loadNotificaciones() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.notificaciones.isEmpty() -> {
                    Text(
                        "No tienes notificaciones",
                        modifier = Modifier.align(Alignment.Center),
                        color = Gray500
                    )
                }
                else -> {
                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        items(uiState.notificaciones) { notificacion ->
                            NotificacionCard(
                                notificacion = notificacion,
                                onClick = {
                                    if (!notificacion.isRead) {
                                        viewModel.markAsRead(notificacion.id)
                                    }
                                }
                            )
                        }
                    }
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificacionCard(
    notificacion: Notificacion,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (notificacion.isRead) Color.White else Blue600.copy(alpha = 0.05f)
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = if (notificacion.isRead) 1.dp else 2.dp),
        onClick = onClick
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            if (!notificacion.isRead) {
                Icon(
                    Icons.Default.Circle,
                    contentDescription = "No leida",
                    tint = Blue600,
                    modifier = Modifier
                        .size(8.dp)
                        .padding(top = 6.dp)
                )
                Spacer(modifier = Modifier.width(8.dp))
            }

            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = notificacion.title,
                    fontWeight = if (notificacion.isRead) FontWeight.Normal else FontWeight.SemiBold,
                    fontSize = 16.sp,
                    color = Gray900
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = notificacion.body,
                    fontSize = 14.sp,
                    color = Gray500
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = formatDate(notificacion.createdAt),
                    fontSize = 12.sp,
                    color = Gray500
                )
            }
        }
    }
}

private fun formatDate(dateString: String): String {
    return try {
        val date = ZonedDateTime.parse(dateString)
        date.format(DateTimeFormatter.ofPattern("dd MMM, HH:mm"))
    } catch (e: Exception) {
        dateString.take(16)
    }
}
