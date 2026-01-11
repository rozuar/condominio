package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Comunicado(
    val id: String,
    val title: String,
    val content: String,
    val type: String,
    // Backend doesn't always send priority; keep default for backwards compatibility.
    val priority: String = "",
    @SerialName("author_id") val authorId: String? = null,
    @SerialName("author_name") val authorName: String? = null,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String
)

/**
 * App-side fallback for older API responses that don't include `priority`.
 * Maps known comunicado types to a priority bucket used by UI coloring.
 */
fun Comunicado.effectivePriority(): String {
    if (priority.isNotBlank()) return priority
    return when (type) {
        "urgente", "seguridad" -> "high"
        "mantenimiento", "tesoreria" -> "medium"
        else -> "low"
    }
}

@Serializable
data class ComunicadoListResponse(
    val comunicados: List<Comunicado>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)
