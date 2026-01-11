package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class PeriodoGasto(
    val year: Int = 0,
    val month: Int = 0
)

@Serializable
data class GastoComun(
    val id: String,
    @SerialName("periodo_id") val periodoId: String,
    @SerialName("parcela_id") val parcelaId: Int,
    @SerialName("parcela_numero") val parcelaNumero: String? = null,
    @SerialName("user_id") val userId: String? = null,
    @SerialName("user_name") val userName: String? = null,
    val monto: Double = 0.0,
    @SerialName("monto_pagado") val montoPagado: Double = 0.0,
    val status: String = "pending",
    @SerialName("fecha_pago") val fechaPago: String? = null,
    @SerialName("metodo_pago") val metodoPago: String? = null,
    @SerialName("referencia_pago") val referenciaPago: String? = null,
    @SerialName("created_at") val createdAt: String? = null,
    @SerialName("updated_at") val updatedAt: String? = null,
    val periodo: PeriodoGasto? = null
)

@Serializable
data class MiEstadoCuenta(
    @SerialName("has_parcela") val hasParcela: Boolean = true,
    val message: String? = null,
    @SerialName("parcela_id") val parcelaId: Int,
    @SerialName("parcela_numero") val parcelaNumero: String,
    @SerialName("gastos_pendientes") val gastosPendientes: List<GastoComun> = emptyList(),
    @SerialName("gastos_pagados") val gastosPagados: List<GastoComun> = emptyList(),
    @SerialName("total_pendiente") val totalPendiente: Double = 0.0,
    @SerialName("total_pagado") val totalPagado: Double = 0.0
)
