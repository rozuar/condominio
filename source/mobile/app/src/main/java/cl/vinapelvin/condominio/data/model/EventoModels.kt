package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Evento(
    val id: String,
    val title: String,
    val description: String = "",
    val location: String = "",
    @SerialName("event_date") val eventDate: String,
    @SerialName("event_end_date") val eventEndDate: String? = null,
    val type: String,
    @SerialName("is_public") val isPublic: Boolean = true,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String? = null
)

@Serializable
data class EventoListResponse(
    val eventos: List<Evento>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
