package cl.vinapelvin.condominio.ui.emergencias

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.Emergencia
import cl.vinapelvin.condominio.ui.theme.*
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EmergenciasScreen(
    viewModel: EmergenciasViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Emergencias") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Volver")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Red600,
                    titleContentColor = Color.White,
                    navigationIconContentColor = Color.White
                )
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
                    CircularProgressIndicator(
                        modifier = Modifier.align(Alignment.Center)
                    )
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(uiState.error!!, color = MaterialTheme.colorScheme.error)
                        Spacer(modifier = Modifier.height(16.dp))
                        Button(onClick = { viewModel.loadEmergencias() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.emergencias.isEmpty() -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Icon(
                            Icons.Default.CheckCircle,
                            contentDescription = null,
                            modifier = Modifier.size(64.dp),
                            tint = Green600
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "No hay emergencias activas",
                            fontWeight = FontWeight.Medium,
                            fontSize = 18.sp,
                            color = Gray900
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(
                            "Todo esta en orden",
                            color = Gray500
                        )
                    }
                }
                else -> {
                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        items(uiState.emergencias) { emergencia ->
                            EmergenciaCard(emergencia = emergencia)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun EmergenciaCard(emergencia: Emergencia) {
    val priorityColor = when (emergencia.priority) {
        "critical" -> Red600
        "high" -> Color(0xFFEA580C)
        "medium" -> Amber500
        else -> Green600
    }

    val priorityLabel = when (emergencia.priority) {
        "critical" -> "Critica"
        "high" -> "Alta"
        "medium" -> "Media"
        else -> "Baja"
    }

    val statusColor = when (emergencia.status) {
        "active" -> Red600
        "in_progress" -> Amber500
        "resolved" -> Green600
        else -> Gray500
    }

    val statusLabel = when (emergencia.status) {
        "active" -> "Activa"
        "in_progress" -> "En progreso"
        "resolved" -> "Resuelta"
        else -> emergencia.status
    }

    val isActive = emergencia.status == "active" || emergencia.status == "in_progress"

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(
            containerColor = if (isActive) priorityColor.copy(alpha = 0.05f) else Color.White
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = if (isActive) 4.dp else 2.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    AssistChip(
                        onClick = {},
                        leadingIcon = {
                            Icon(
                                Icons.Default.Warning,
                                contentDescription = null,
                                modifier = Modifier.size(16.dp),
                                tint = priorityColor
                            )
                        },
                        label = { Text(priorityLabel, fontSize = 12.sp) },
                        colors = AssistChipDefaults.assistChipColors(
                            containerColor = priorityColor.copy(alpha = 0.1f),
                            labelColor = priorityColor
                        )
                    )
                    AssistChip(
                        onClick = {},
                        label = { Text(statusLabel, fontSize = 12.sp) },
                        colors = AssistChipDefaults.assistChipColors(
                            containerColor = statusColor.copy(alpha = 0.1f),
                            labelColor = statusColor
                        )
                    )
                }
            }

            Spacer(modifier = Modifier.height(12.dp))

            Text(
                text = emergencia.title,
                fontWeight = FontWeight.SemiBold,
                fontSize = 16.sp,
                color = Gray900
            )

            Spacer(modifier = Modifier.height(4.dp))

            Text(
                text = emergencia.description,
                fontSize = 14.sp,
                color = Gray500,
                maxLines = 3,
                overflow = TextOverflow.Ellipsis
            )

            Spacer(modifier = Modifier.height(12.dp))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.Person,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp),
                        tint = Gray500
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = emergencia.reporterName ?: "Anonimo",
                        fontSize = 12.sp,
                        color = Gray500
                    )
                }

                Text(
                    text = formatDate(emergencia.createdAt),
                    fontSize = 12.sp,
                    color = Gray500
                )
            }

            if (emergencia.status == "resolved" && emergencia.resolvedAt != null) {
                Spacer(modifier = Modifier.height(8.dp))
                Divider(color = Gray100)
                Spacer(modifier = Modifier.height(8.dp))
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        Icons.Default.CheckCircle,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp),
                        tint = Green600
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "Resuelta el ${formatDate(emergencia.resolvedAt)}",
                        fontSize = 12.sp,
                        color = Green600
                    )
                }
            }
        }
    }
}

private fun formatDate(dateString: String): String {
    return try {
        val date = ZonedDateTime.parse(dateString)
        date.format(DateTimeFormatter.ofPattern("dd MMM yyyy, HH:mm"))
    } catch (e: Exception) {
        dateString.take(16)
    }
}
