package cl.vinapelvin.condominio.ui.gastos

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
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import cl.vinapelvin.condominio.data.model.GastoComun
import cl.vinapelvin.condominio.ui.theme.*
import java.text.NumberFormat
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GastosScreen(
    viewModel: GastosViewModel = hiltViewModel(),
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Mis Gastos") },
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
                    CircularProgressIndicator(modifier = Modifier.align(Alignment.Center))
                }
                uiState.error != null -> {
                    Column(
                        modifier = Modifier.align(Alignment.Center),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(uiState.error!!, color = MaterialTheme.colorScheme.error)
                        Spacer(modifier = Modifier.height(16.dp))
                        Button(onClick = { viewModel.loadEstadoCuenta() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.estadoCuenta != null -> {
                    val estado = uiState.estadoCuenta!!

                    LazyColumn(
                        modifier = Modifier.fillMaxSize(),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        // Summary cards
                        item {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.spacedBy(12.dp)
                            ) {
                                SummaryCard(
                                    title = "Pendiente",
                                    amount = estado.totalPending,
                                    color = Amber500,
                                    modifier = Modifier.weight(1f)
                                )
                                SummaryCard(
                                    title = "Vencido",
                                    amount = estado.totalOverdue,
                                    color = Red600,
                                    modifier = Modifier.weight(1f)
                                )
                            }
                        }

                        // Pending section
                        if (estado.gastosPending.isNotEmpty()) {
                            item {
                                Text(
                                    "Pendientes",
                                    fontWeight = FontWeight.SemiBold,
                                    fontSize = 18.sp,
                                    color = Gray900
                                )
                            }
                            items(estado.gastosPending) { gasto ->
                                GastoCard(gasto)
                            }
                        }

                        // Overdue section
                        if (estado.gastosOverdue.isNotEmpty()) {
                            item {
                                Text(
                                    "Vencidos",
                                    fontWeight = FontWeight.SemiBold,
                                    fontSize = 18.sp,
                                    color = Red600
                                )
                            }
                            items(estado.gastosOverdue) { gasto ->
                                GastoCard(gasto)
                            }
                        }

                        if (estado.gastosPending.isEmpty() && estado.gastosOverdue.isEmpty()) {
                            item {
                                Card(
                                    modifier = Modifier.fillMaxWidth(),
                                    colors = CardDefaults.cardColors(
                                        containerColor = Green600.copy(alpha = 0.1f)
                                    )
                                ) {
                                    Text(
                                        "No tienes gastos pendientes",
                                        modifier = Modifier.padding(24.dp),
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

@Composable
fun SummaryCard(
    title: String,
    amount: Double,
    color: Color,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = color.copy(alpha = 0.1f))
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(title, fontSize = 14.sp, color = color)
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                text = formatCurrency(amount),
                fontWeight = FontWeight.Bold,
                fontSize = 20.sp,
                color = color
            )
        }
    }
}

@Composable
fun GastoCard(gasto: GastoComun) {
    val statusColor = when (gasto.status) {
        "pending" -> Amber500
        "overdue" -> Red600
        "paid" -> Green600
        else -> Gray500
    }

    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(
                    "Parcela ${gasto.parcelaId}",
                    fontWeight = FontWeight.Medium,
                    color = Gray900
                )
                Text(
                    "Vence: ${gasto.fechaVencimiento.take(10)}",
                    fontSize = 12.sp,
                    color = Gray500
                )
            }
            Column(horizontalAlignment = Alignment.End) {
                Text(
                    formatCurrency(gasto.montoTotal),
                    fontWeight = FontWeight.Bold,
                    color = statusColor
                )
                AssistChip(
                    onClick = {},
                    label = {
                        Text(
                            when (gasto.status) {
                                "pending" -> "Pendiente"
                                "overdue" -> "Vencido"
                                "paid" -> "Pagado"
                                else -> gasto.status
                            },
                            fontSize = 10.sp
                        )
                    },
                    colors = AssistChipDefaults.assistChipColors(
                        containerColor = statusColor.copy(alpha = 0.1f),
                        labelColor = statusColor
                    )
                )
            }
        }
    }
}

private fun formatCurrency(amount: Double): String {
    val format = NumberFormat.getCurrencyInstance(Locale("es", "CL"))
    return format.format(amount)
}
