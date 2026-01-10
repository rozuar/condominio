package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class GastoComun(
    val id: String,
    @SerialName("periodo_id") val periodoId: String,
    @SerialName("user_id") val userId: String,
    @SerialName("parcela_id") val parcelaId: Int,
    @SerialName("monto_base") val montoBase: Double,
    @SerialName("monto_extra") val montoExtra: Double,
    @SerialName("monto_total") val montoTotal: Double,
    val status: String,
    @SerialName("fecha_vencimiento") val fechaVencimiento: String,
    @SerialName("fecha_pago") val fechaPago: String? = null,
    @SerialName("metodo_pago") val metodoPago: String? = null
)

@Serializable
data class EstadoCuenta(
    @SerialName("user_id") val userId: String,
    @SerialName("parcela_id") val parcelaId: Int,
    @SerialName("total_pending") val totalPending: Double,
    @SerialName("total_overdue") val totalOverdue: Double,
    @SerialName("gastos_pending") val gastosPending: List<GastoComun>,
    @SerialName("gastos_overdue") val gastosOverdue: List<GastoComun>,
    @SerialName("ultimo_pago") val ultimoPago: GastoComun? = null
)
