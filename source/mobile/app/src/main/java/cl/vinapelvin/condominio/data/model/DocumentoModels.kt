package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Documento(
    val id: String,
    val title: String,
    val description: String? = null,
    @SerialName("file_url") val fileUrl: String? = null,
    val category: String,
    @SerialName("is_public") val isPublic: Boolean,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class DocumentoListResponse(
    val documentos: List<Documento>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

