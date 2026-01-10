package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Acta(
    val id: String,
    val title: String,
    val content: String,
    @SerialName("meeting_date") val meetingDate: String,
    val type: String,
    @SerialName("attendees_count") val attendeesCount: Int? = null,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class ActaListResponse(
    val actas: List<Acta>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

