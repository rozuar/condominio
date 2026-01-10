package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Comunicado(
    val id: String,
    val title: String,
    val content: String,
    val type: String,
    val priority: String,
    @SerialName("author_id") val authorId: String,
    @SerialName("author_name") val authorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class ComunicadoListResponse(
    val comunicados: List<Comunicado>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
