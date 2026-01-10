package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Movimiento(
    val id: String,
    val description: String,
    val amount: Double,
    val type: String,
    val category: String? = null,
    val date: String,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class MovimientoListResponse(
    val movimientos: List<Movimiento>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

@Serializable
data class TesoreriaResumen(
    @SerialName("total_ingresos") val totalIngresos: Double,
    @SerialName("total_egresos") val totalEgresos: Double,
    val balance: Double,
    @SerialName("movimientos_count") val movimientosCount: Int? = null
)

