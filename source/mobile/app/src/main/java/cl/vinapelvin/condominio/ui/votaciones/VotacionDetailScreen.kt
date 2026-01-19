package cl.vinapelvin.condominio.ui.votaciones

import androidx.compose.animation.animateContentSize
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.HowToVote
import androidx.compose.material.icons.filled.Groups
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.Votacion
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

                        val showResults = votacion.userVoted || votacion.status == "closed"

                        // Results summary when showing results
                        if (showResults && votacion.totalVotes > 0) {
                            ResultsSummaryCard(votacion = votacion)
                            Spacer(modifier = Modifier.height(24.dp))
                        }

                        Text(
                            text = if (showResults) "Resultados" else "Opciones",
                            fontWeight = FontWeight.SemiBold,
                            fontSize = 18.sp,
                            color = Gray900
                        )

                        Spacer(modifier = Modifier.height(12.dp))

                        votacion.options.forEach { option ->
                            if (showResults) {
                                ResultOptionCard(
                                    option = option,
                                    totalVotes = votacion.totalVotes
                                )
                            } else {
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
                            }
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

@Composable
private fun ResultsSummaryCard(votacion: Votacion) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = Blue600.copy(alpha = 0.05f)
        ),
        border = BorderStroke(1.dp, Blue600.copy(alpha = 0.2f))
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        imageVector = Icons.Default.HowToVote,
                        contentDescription = null,
                        tint = Blue600,
                        modifier = Modifier.size(20.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(
                        text = "Total de votos",
                        fontSize = 14.sp,
                        color = Gray700
                    )
                }
                Text(
                    text = "${votacion.totalVotes}",
                    fontWeight = FontWeight.Bold,
                    fontSize = 18.sp,
                    color = Blue600
                )
            }

            if (votacion.requiresQuorum) {
                Spacer(modifier = Modifier.height(12.dp))
                Divider(color = Gray200)
                Spacer(modifier = Modifier.height(12.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            imageVector = Icons.Default.Groups,
                            contentDescription = null,
                            tint = Gray600,
                            modifier = Modifier.size(20.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(
                            text = "Quórum requerido",
                            fontSize = 14.sp,
                            color = Gray700
                        )
                    }
                    Text(
                        text = "${votacion.quorumPercentage}%",
                        fontWeight = FontWeight.SemiBold,
                        fontSize = 14.sp,
                        color = Gray700
                    )
                }
            }
        }
    }
}

@Composable
private fun ResultOptionCard(
    option: VotacionOpcion,
    totalVotes: Int
) {
    val percentage = if (totalVotes > 0) {
        (option.votes.toFloat() / totalVotes.toFloat()) * 100
    } else {
        0f
    }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = BorderStroke(1.dp, Gray100)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
                .animateContentSize()
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = option.text,
                    fontWeight = FontWeight.SemiBold,
                    color = Gray900,
                    modifier = Modifier.weight(1f)
                )
                Text(
                    text = "${option.votes} votos",
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Medium,
                    color = Blue600
                )
            }

            Spacer(modifier = Modifier.height(12.dp))

            // Progress bar
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(8.dp)
                    .clip(RoundedCornerShape(4.dp))
                    .background(Gray100)
            ) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth(fraction = (percentage / 100f).coerceIn(0f, 1f))
                        .fillMaxHeight()
                        .clip(RoundedCornerShape(4.dp))
                        .background(
                            if (percentage >= 50) Green600 else Blue600
                        )
                )
            }

            Spacer(modifier = Modifier.height(8.dp))

            Text(
                text = String.format("%.1f%%", percentage),
                fontSize = 13.sp,
                fontWeight = FontWeight.Medium,
                color = if (percentage >= 50) Green600 else Gray600
            )
        }
    }
}
