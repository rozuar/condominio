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
        containerColor = Gray50,
        topBar = {
            TopAppBar(
                title = { Text("Mis Gastos") },
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
                        Button(onClick = { viewModel.loadEstadoCuenta() }) {
                            Text("Reintentar")
                        }
                    }
                }
                uiState.estadoCuenta != null -> {
                    val estado = uiState.estadoCuenta!!
                    if (!estado.hasParcela) {
                        Column(
                            modifier = Modifier
                                .fillMaxSize()
                                .padding(24.dp),
                            horizontalAlignment = Alignment.CenterHorizontally,
                            verticalArrangement = Arrangement.Center
                        ) {
                            Text(
                                text = estado.message ?: "Ud no posee asociada una parcela que genere gastos",
                                color = Gray900,
                                fontWeight = FontWeight.Medium
                            )
                            Spacer(modifier = Modifier.height(16.dp))
                            Button(onClick = { viewModel.loadEstadoCuenta() }) {
                                Text("Reintentar")
                            }
                        }
                        return@Box
                    }
                    val gastosOverdue = estado.gastosPendientes.filter { it.status == "overdue" }
                    val gastosPending = estado.gastosPendientes.filter { it.status != "overdue" }
                    val totalOverdue = gastosOverdue.sumOf { (it.monto - it.montoPagado).coerceAtLeast(0.0) }
                    val totalPending = gastosPending.sumOf { (it.monto - it.montoPagado).coerceAtLeast(0.0) }

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
                                    amount = totalPending,
                                    color = Amber500,
                                    modifier = Modifier.weight(1f)
                                )
                                SummaryCard(
                                    title = "Vencido",
                                    amount = totalOverdue,
                                    color = Red600,
                                    modifier = Modifier.weight(1f)
                                )
                            }
                        }

                        item {
                            SummaryCard(
                                title = "Pagado",
                                amount = estado.totalPagado,
                                color = Green600,
                                modifier = Modifier.fillMaxWidth()
                            )
                        }

                        // Pending section
                        if (gastosPending.isNotEmpty()) {
                            item {
                                Text(
                                    "Pendientes",
                                    fontWeight = FontWeight.SemiBold,
                                    fontSize = 18.sp,
                                    color = Gray900
                                )
                            }
                            items(gastosPending) { gasto ->
                                GastoCard(gasto, estado.parcelaNumero)
                            }
                        }

                        // Overdue section
                        if (gastosOverdue.isNotEmpty()) {
                            item {
                                Text(
                                    "Vencidos",
                                    fontWeight = FontWeight.SemiBold,
                                    fontSize = 18.sp,
                                    color = Red600
                                )
                            }
                            items(gastosOverdue) { gasto ->
                                GastoCard(gasto, estado.parcelaNumero)
                            }
                        }

                        if (estado.gastosPendientes.isEmpty()) {
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
fun GastoCard(gasto: GastoComun, parcelaNumeroFallback: String) {
    val statusColor = when (gasto.status) {
        "pending" -> Amber500
        "overdue" -> Red600
        "paid" -> Green600
        else -> Gray500
    }
    val montoPendiente = (gasto.monto - gasto.montoPagado).coerceAtLeast(0.0)
    val periodoLabel = gasto.periodo?.let { p ->
        if (p.year > 0 && p.month > 0) "%04d-%02d".format(p.year, p.month) else null
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
                    "Parcela ${gasto.parcelaNumero ?: parcelaNumeroFallback}",
                    fontWeight = FontWeight.Medium,
                    color = Gray900
                )
                Text(
                    periodoLabel?.let { "Periodo: $it" } ?: "Periodo: -",
                    fontSize = 12.sp,
                    color = Gray500
                )
            }
            Column(horizontalAlignment = Alignment.End) {
                Text(
                    formatCurrency(montoPendiente),
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
