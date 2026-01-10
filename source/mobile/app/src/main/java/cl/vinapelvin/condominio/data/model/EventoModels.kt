package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Evento(
    val id: String,
    val title: String,
    val description: String,
    val location: String,
    @SerialName("start_date") val startDate: String,
    @SerialName("end_date") val endDate: String,
    val type: String,
    @SerialName("is_mandatory") val isMandatory: Boolean,
    @SerialName("created_at") val createdAt: String
)

@Serializable
data class EventoListResponse(
    val eventos: List<Evento>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
