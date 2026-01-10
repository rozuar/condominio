package cl.vinapelvin.condominio.ui.comunicados

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
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
import cl.vinapelvin.condominio.data.model.Comunicado
import cl.vinapelvin.condominio.ui.theme.*
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ComunicadosScreen(
    viewModel: ComunicadosViewModel = hiltViewModel(),
    onBack: () -> Unit,
    onComunicadoClick: (String) -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Comunicados") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Volver")
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
                        Button(onClick = { viewModel.loadComunicados() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.comunicados.isEmpty() -> {
                    Text(
                        "No hay comunicados",
                        modifier = Modifier.align(Alignment.Center),
                        color = Gray500
                    )
                }
                else -> {
                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        items(uiState.comunicados) { comunicado ->
                            ComunicadoCard(
                                comunicado = comunicado,
                                onClick = { onComunicadoClick(comunicado.id) }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun ComunicadoCard(
    comunicado: Comunicado,
    onClick: () -> Unit
) {
    val priorityColor = when (comunicado.priority) {
        "high" -> Red600
        "medium" -> Amber500
        else -> Green600
    }

    val typeLabel = when (comunicado.type) {
        "urgente" -> "Urgente"
        "mantenimiento" -> "Mantenimiento"
        "seguridad" -> "Seguridad"
        "evento" -> "Evento"
        else -> "General"
    }

    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                AssistChip(
                    onClick = {},
                    label = { Text(typeLabel, fontSize = 12.sp) },
                    colors = AssistChipDefaults.assistChipColors(
                        containerColor = priorityColor.copy(alpha = 0.1f),
                        labelColor = priorityColor
                    )
                )
                Text(
                    text = formatDate(comunicado.createdAt),
                    fontSize = 12.sp,
                    color = Gray500
                )
            }

            Spacer(modifier = Modifier.height(8.dp))

            Text(
                text = comunicado.title,
                fontWeight = FontWeight.SemiBold,
                fontSize = 16.sp,
                color = Gray900
            )

            Spacer(modifier = Modifier.height(4.dp))

            Text(
                text = comunicado.content,
                fontSize = 14.sp,
                color = Gray500,
                maxLines = 2,
                overflow = TextOverflow.Ellipsis
            )

            comunicado.authorName?.let { author ->
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = "Por: $author",
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
        date.format(DateTimeFormatter.ofPattern("dd MMM yyyy"))
    } catch (e: Exception) {
        dateString.take(10)
    }
}
