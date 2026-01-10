package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class MensajeContacto(
    val id: String,
    @SerialName("user_id") val userId: String? = null,
    @SerialName("user_name") val userName: String? = null,
    val nombre: String,
    val email: String,
    val asunto: String,
    val mensaje: String,
    val status: String,
    @SerialName("read_at") val readAt: String? = null,
    @SerialName("read_by") val readBy: String? = null,
    @SerialName("read_by_name") val readByName: String? = null,
    @SerialName("replied_at") val repliedAt: String? = null,
    @SerialName("replied_by") val repliedBy: String? = null,
    @SerialName("replied_by_name") val repliedByName: String? = null,
    val respuesta: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class CreateMensajeContactoRequest(
    val nombre: String,
    val email: String,
    val asunto: String,
    val mensaje: String
)

@Serializable
data class MisMensajesResponse(
    val mensajes: List<MensajeContacto>,
    val total: Int
)

