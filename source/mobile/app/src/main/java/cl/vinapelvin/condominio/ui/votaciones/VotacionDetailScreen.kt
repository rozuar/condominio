package cl.vinapelvin.condominio.ui.votaciones

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
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
import cl.vinapelvin.condominio.data.model.VotacionOpcion
import cl.vinapelvin.condominio.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VotacionDetailScreen(
    votacionId: String,
    viewModel: VotacionDetailViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    val userParcelaId by viewModel.userParcelaId.collectAsState(initial = null)

    LaunchedEffect(votacionId) {
        viewModel.loadVotacion(votacionId)
    }

    Scaffold(
        containerColor = Gray50,
        topBar = {
            TopAppBar(
                title = { Text("Votacion") },
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
                    CircularProgressIndicator(modifier = Modifier.align(Alignment.Center))
                }
                uiState.error != null -> {
                    Text(
                        uiState.error!!,
                        modifier = Modifier.align(Alignment.Center),
                        color = MaterialTheme.colorScheme.error
                    )
                }
                uiState.votacion != null -> {
                    val votacion = uiState.votacion!!
                    val canVote = userParcelaId != null

                    Column(
                        modifier = Modifier
                            .fillMaxSize()
                            .verticalScroll(rememberScrollState())
                            .padding(16.dp)
                    ) {
                        Text(
                            text = votacion.title,
                            fontWeight = FontWeight.Bold,
                            fontSize = 24.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(8.dp))

                        val statusColor = when (votacion.status) {
                            "active" -> Green600
                            "closed" -> Gray500
                            else -> Amber500
                        }
                        AssistChip(
                            onClick = {},
                            label = {
                                Text(
                                    when (votacion.status) {
                                        "active" -> "Activa"
                                        "closed" -> "Cerrada"
                                        else -> votacion.status
                                    }
                                )
                            },
                            colors = AssistChipDefaults.assistChipColors(
                                containerColor = statusColor.copy(alpha = 0.1f),
                                labelColor = statusColor
                            )
                        )

                        Spacer(modifier = Modifier.height(16.dp))

                        Text(
                            text = votacion.description,
                            fontSize = 16.sp,
                            color = Gray500
                        )

                        Spacer(modifier = Modifier.height(24.dp))

                        Text(
                            text = "Opciones",
                            fontWeight = FontWeight.SemiBold,
                            fontSize = 18.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(12.dp))

                        votacion.options.forEach { option ->
                            OptionCard(
                                option = option,
                                isSelected = uiState.selectedOptionId == option.id,
                                hasVoted = votacion.userVoted,
                                isActive = votacion.status == "active",
                                canVote = canVote,
                                onClick = {
                                    if (canVote && !votacion.userVoted && votacion.status == "active") {
                                        viewModel.selectOption(option.id)
                                    }
                                }
                            )
                            Spacer(modifier = Modifier.height(8.dp))
                        }

                        if (!votacion.userVoted && votacion.status == "active") {
                            Spacer(modifier = Modifier.height(16.dp))

                            if (!canVote) {
                                Card(
                                    modifier = Modifier.fillMaxWidth(),
                                    colors = CardDefaults.cardColors(
                                        containerColor = Amber500.copy(alpha = 0.10f)
                                    )
                                ) {
                                    Text(
                                        text = "Solo puedes ver el estado de la votación. Para votar debes tener una parcela asociada.",
                                        modifier = Modifier.padding(16.dp),
                                        color = Gray900
                                    )
                                }
                                return@Column
                            }

                            Button(
                                onClick = { viewModel.submitVote(votacionId) },
                                modifier = Modifier.fillMaxWidth(),
                                enabled = uiState.selectedOptionId != null && !uiState.isVoting
                            ) {
                                if (uiState.isVoting) {
                                    CircularProgressIndicator(
                                        modifier = Modifier.size(24.dp),
                                        color = Color.White
                                    )
                                } else {
                                    Text("Emitir Voto")
                                }
                            }
                        }

                        if (votacion.userVoted) {
                            Spacer(modifier = Modifier.height(16.dp))
                            Card(
                                modifier = Modifier.fillMaxWidth(),
                                colors = CardDefaults.cardColors(
                                    containerColor = Green600.copy(alpha = 0.1f)
                                )
                            ) {
                                Row(
                                    modifier = Modifier.padding(16.dp),
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        Icons.Default.Check,
                                        contentDescription = null,
                                        tint = Green600
                                    )
                                    Spacer(modifier = Modifier.width(8.dp))
                                    Text(
                                        "Ya emitiste tu voto",
                                        color = Green600,
                                        fontWeight = FontWeight.Medium
                                    )
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun OptionCard(
    option: VotacionOpcion,
    isSelected: Boolean,
    hasVoted: Boolean,
    isActive: Boolean,
    canVote: Boolean,
    onClick: () -> Unit
) {
    val borderColor = when {
        isSelected -> Blue600
        else -> Gray100
    }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        border = BorderStroke(2.dp, borderColor),
        colors = CardDefaults.cardColors(
            containerColor = if (isSelected) borderColor.copy(alpha = 0.1f) else Color.White
        ),
        onClick = onClick,
        enabled = canVote && !hasVoted && isActive
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = option.text,
                fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Normal,
                color = Gray900
            )

            if (!isActive || hasVoted) {
                Text(
                    text = "${option.votes} votos",
                    fontSize = 12.sp,
                    color = Gray500
                )
            }
        }
    }
}
