package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Galeria(
    val id: String,
    val title: String,
    val description: String? = null,
    @SerialName("event_date") val eventDate: String? = null,
    @SerialName("is_public") val isPublic: Boolean,
    @SerialName("cover_image_url") val coverImageUrl: String? = null,
    @SerialName("items_count") val itemsCount: Int = 0,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class GaleriaItem(
    val id: String,
    @SerialName("galeria_id") val galeriaId: String,
    @SerialName("file_url") val fileUrl: String,
    @SerialName("thumbnail_url") val thumbnailUrl: String? = null,
    @SerialName("file_type") val fileType: String, // "image" or "video"
    val caption: String? = null,
    @SerialName("order_index") val orderIndex: Int = 0,
    @SerialName("created_at") val createdAt: String
)

@Serializable
data class GaleriaWithItems(
    val id: String,
    val title: String,
    val description: String? = null,
    @SerialName("event_date") val eventDate: String? = null,
    @SerialName("is_public") val isPublic: Boolean,
    @SerialName("cover_image_url") val coverImageUrl: String? = null,
    @SerialName("items_count") val itemsCount: Int = 0,
    @SerialName("created_by") val createdBy: String? = null,
    @SerialName("creator_name") val creatorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String,
    val items: List<GaleriaItem> = emptyList()
)

@Serializable
data class GaleriaListResponse(
    val galerias: List<Galeria>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
