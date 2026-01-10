package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Notificacion(
    val id: String,
    @SerialName("user_id") val userId: String,
    val title: String,
    val body: String,
    val type: String,
    @SerialName("is_read") val isRead: Boolean,
    @SerialName("read_at") val readAt: String? = null,
    @SerialName("created_at") val createdAt: String
)

@Serializable
data class NotificacionListResponse(
    val notificaciones: List<Notificacion>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

@Serializable
data class NotificacionStats(
    val total: Int,
    val unread: Int
)
