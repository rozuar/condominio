package cl.vinapelvin.condominio.ui.eventos

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.ui.theme.*
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EventoDetailScreen(
    eventoId: String,
    viewModel: EventoDetailViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        containerColor = Gray50,
        topBar = {
            TopAppBar(
                title = { Text("Detalle del Evento") },
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
                        Button(onClick = { viewModel.loadEvento() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.evento != null -> {
                    val evento = uiState.evento!!

                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .verticalScroll(rememberScrollState())
                            .padding(16.dp)
                    ) {
                        val typeColor = when (evento.type) {
                            "asamblea" -> Blue600
                            "reunion" -> Amber500
                            "mantenimiento" -> Green600
                            "social" -> Color(0xFF8B5CF6)
                            else -> Gray500
                        }

                        val typeLabel = when (evento.type) {
                            "asamblea" -> "Asamblea"
                            "reunion" -> "Reunion"
                            "mantenimiento" -> "Mantenimiento"
                            "social" -> "Social"
                            else -> "Evento"
                        }

                        Row(
                            horizontalArrangement = Arrangement.spacedBy(8.dp)
                        ) {
                            AssistChip(
                                onClick = {},
                                label = { Text(typeLabel) },
                                colors = AssistChipDefaults.assistChipColors(
                                    containerColor = typeColor.copy(alpha = 0.1f),
                                    labelColor = typeColor
                                )
                            )
                        }

                        Spacer(modifier = Modifier.height(16.dp))

                        Text(
                            text = evento.title,
                            fontWeight = FontWeight.Bold,
                            fontSize = 24.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(24.dp))

                        // Info cards
                        Card(
                            modifier = Modifier.fillMaxWidth(),
                            shape = RoundedCornerShape(12.dp),
                            colors = CardDefaults.cardColors(containerColor = Gray100)
                        ) {
                            Column(modifier = Modifier.padding(16.dp)) {
                                Row(
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        Icons.Default.CalendarToday,
                                        contentDescription = null,
                                        tint = Blue600
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Column {
                                        Text(
                                            "Fecha de inicio",
                                            fontSize = 12.sp,
                                            color = Gray500
                                        )
                                        Text(
                                            formatDateTime(evento.eventDate),
                                            fontWeight = FontWeight.Medium,
                                            color = Gray900
                                        )
                                    }
                                }

                                Spacer(modifier = Modifier.height(16.dp))

                                Row(
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        Icons.Default.Schedule,
                                        contentDescription = null,
                                        tint = Blue600
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Column {
                                        Text(
                                            "Fecha de termino",
                                            fontSize = 12.sp,
                                            color = Gray500
                                        )
                                        Text(
                                            evento.eventEndDate?.let { formatDateTime(it) } ?: "Sin término",
                                            fontWeight = FontWeight.Medium,
                                            color = Gray900
                                        )
                                    }
                                }

                                Spacer(modifier = Modifier.height(16.dp))

                                Row(
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        Icons.Default.LocationOn,
                                        contentDescription = null,
                                        tint = Blue600
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Column {
                                        Text(
                                            "Ubicacion",
                                            fontSize = 12.sp,
                                            color = Gray500
                                        )
                                        Text(
                                            evento.location.ifBlank { "Sin ubicación" },
                                            fontWeight = FontWeight.Medium,
                                            color = Gray900
                                        )
                                    }
                                }
                            }
                        }

                        Spacer(modifier = Modifier.height(24.dp))

                        Text(
                            "Descripcion",
                            fontWeight = FontWeight.SemiBold,
                            fontSize = 18.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(8.dp))

                        Text(
                            text = evento.description,
                            fontSize = 15.sp,
                            color = Gray500,
                            lineHeight = 24.sp
                        )
                    }
                }
            }
        }
    }
}

private fun formatDateTime(dateString: String): String {
    return try {
        val date = ZonedDateTime.parse(dateString)
        date.format(DateTimeFormatter.ofPattern("dd MMMM yyyy, HH:mm 'hrs'"))
    } catch (e: Exception) {
        dateString.take(16)
    }
}
