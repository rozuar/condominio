package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
data class MapaArea(
    val id: String,
    @SerialName("parcela_id") val parcelaId: Int? = null,
    val type: String, // parcela, area_comun, acceso, canal, camino
    val name: String,
    val description: String? = null,
    val coordinates: JsonElement, // GeoJSON polygon coordinates
    @SerialName("center_lat") val centerLat: Double? = null,
    @SerialName("center_lng") val centerLng: Double? = null,
    @SerialName("fill_color") val fillColor: String,
    @SerialName("stroke_color") val strokeColor: String,
    @SerialName("is_clickable") val isClickable: Boolean,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class MapaPunto(
    val id: String,
    val name: String,
    val description: String? = null,
    val lat: Double,
    val lng: Double,
    val icon: String,
    val type: String,
    @SerialName("is_public") val isPublic: Boolean,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

@Serializable
data class MapaData(
    val areas: List<MapaArea> = emptyList(),
    val puntos: List<MapaPunto> = emptyList()
)

@Serializable
data class MapaAreaListResponse(
    val areas: List<MapaArea>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

@Serializable
data class MapaPuntoListResponse(
    val puntos: List<MapaPunto>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
