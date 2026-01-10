package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Emergencia(
    val id: String,
    val title: String,
    val description: String,
    val priority: String,
    val status: String,
    @SerialName("reported_by") val reportedBy: String,
    @SerialName("reporter_name") val reporterName: String? = null,
    @SerialName("resolved_at") val resolvedAt: String? = null,
    @SerialName("resolved_by") val resolvedBy: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class EmergenciaListResponse(
    val emergencias: List<Emergencia>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
