package cl.vinapelvin.condominio.ui.votaciones

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.Votacion
import cl.vinapelvin.condominio.ui.theme.*
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VotacionesScreen(
    viewModel: VotacionesViewModel = hiltViewModel(),
    onBack: () -> Unit,
    onVotacionClick: (String) -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        containerColor = Gray50,
        topBar = {
            TopAppBar(
                title = { Text("Votaciones") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Volver")
                    }
                }
                ,
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
                    CircularProgressIndicator(modifier = Modifier.align(Alignment.Center))
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(uiState.error!!, color = MaterialTheme.colorScheme.error)
                        Spacer(modifier = Modifier.height(16.dp))
                        Button(onClick = { viewModel.loadVotaciones() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.votaciones.isEmpty() -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(
                            "No hay votaciones",
                            fontWeight = FontWeight.Medium,
                            fontSize = 18.sp,
                            color = Gray900
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text("Cuando haya una nueva, aparecerá aquí.", color = Gray500)
                    }
                }
                else -> {
                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        items(uiState.votaciones) { votacion ->
                            VotacionCard(
                                votacion = votacion,
                                onClick = { onVotacionClick(votacion.id) }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun VotacionCard(
    votacion: Votacion,
    onClick: () -> Unit
) {
    val statusColor = when (votacion.status) {
        "active" -> Green600
        "closed" -> Gray500
        "draft" -> Amber500
        else -> Red600
    }

    val statusLabel = when (votacion.status) {
        "active" -> "Activa"
        "closed" -> "Cerrada"
        "draft" -> "Borrador"
        else -> "Cancelada"
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
                    label = { Text(statusLabel, fontSize = 12.sp) },
                    colors = AssistChipDefaults.assistChipColors(
                        containerColor = statusColor.copy(alpha = 0.1f),
                        labelColor = statusColor
                    )
                )

                if (votacion.userVoted) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            Icons.Default.Check,
                            contentDescription = null,
                            tint = Green600,
                            modifier = Modifier.size(16.dp)
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("Votaste", fontSize = 12.sp, color = Green600)
                    }
                }
            }

            Spacer(modifier = Modifier.height(8.dp))

            Text(
                text = votacion.title,
                fontWeight = FontWeight.SemiBold,
                fontSize = 16.sp,
                color = Gray900
            )

            Spacer(modifier = Modifier.height(4.dp))

            Text(
                text = votacion.description,
                fontSize = 14.sp,
                color = Gray500,
                maxLines = 2
            )

            Spacer(modifier = Modifier.height(8.dp))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = if (votacion.options.isNotEmpty()) {
                        "${votacion.options.size} opciones"
                    } else {
                        "${votacion.totalVotes} votos"
                    },
                    fontSize = 12.sp,
                    color = Gray500
                )
                Text(
                    text = votacion.endDate?.let { "Hasta: ${formatDate(it)}" } ?: "Sin fecha límite",
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
        date.format(DateTimeFormatter.ofPattern("dd MMM"))
    } catch (e: Exception) {
        dateString.take(10)
    }
}
