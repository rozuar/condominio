package cl.vinapelvin.condominio.ui.comunicados

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.ui.theme.*
import cl.vinapelvin.condominio.data.model.effectivePriority
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ComunicadoDetailScreen(
    comunicadoId: String,
    viewModel: ComunicadoDetailViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(comunicadoId) {
        viewModel.loadComunicado(comunicadoId)
    }

    Scaffold(
        containerColor = Gray50,
        topBar = {
            TopAppBar(
                title = { Text("Comunicado") },
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
                    Text(
                        uiState.error!!,
                        modifier = Modifier.align(Alignment.Center),
                        color = MaterialTheme.colorScheme.error
                    )
                }
                uiState.comunicado != null -> {
                    val comunicado = uiState.comunicado!!

                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .verticalScroll(rememberScrollState())
                            .padding(16.dp)
                    ) {
                        Text(
                            text = comunicado.title,
                            fontWeight = FontWeight.Bold,
                            fontSize = 24.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(8.dp))

                        Row(
                            horizontalArrangement = Arrangement.spacedBy(8.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            val priorityColor = when (comunicado.effectivePriority()) {
                                "high" -> Red600
                                "medium" -> Amber500
                                else -> Green600
                            }
                            AssistChip(
                                onClick = {},
                                label = { Text(comunicado.type.replaceFirstChar { it.uppercase() }) },
                                colors = AssistChipDefaults.assistChipColors(
                                    containerColor = priorityColor.copy(alpha = 0.1f),
                                    labelColor = priorityColor
                                )
                            )

                            Text(
                                text = formatDate(comunicado.createdAt),
                                fontSize = 14.sp,
                                color = Gray500
                            )
                        }

                        Spacer(modifier = Modifier.height(16.dp))

                        Divider()

                        Spacer(modifier = Modifier.height(16.dp))

                        Text(
                            text = comunicado.content,
                            fontSize = 16.sp,
                            color = Gray900,
                            lineHeight = 24.sp
                        )

                        comunicado.authorName?.let { author ->
                            Spacer(modifier = Modifier.height(24.dp))
                            Text(
                                text = "Publicado por: $author",
                                fontSize = 14.sp,
                                color = Gray500
                            )
                        }
                    }
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
