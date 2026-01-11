package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class VotacionOpcion(
    val id: String,
    @SerialName("label") val text: String,
    @SerialName("votos_count") val votes: Int = 0
)

@Serializable
data class Votacion(
    val id: String,
    val title: String,
    val description: String = "",
    @SerialName("opciones") val options: List<VotacionOpcion> = emptyList(),
    @SerialName("start_date") val startDate: String? = null,
    @SerialName("end_date") val endDate: String? = null,
    val status: String,
    @SerialName("requires_quorum") val requiresQuorum: Boolean,
    @SerialName("quorum_percentage") val quorumPercentage: Int,
    @SerialName("allow_abstention") val allowAbstention: Boolean = true,
    @SerialName("total_votos") val totalVotes: Int = 0,
    @SerialName("has_voted") val userVoted: Boolean = false,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String? = null
)

@Serializable
data class VotacionListResponse(
    val votaciones: List<Votacion>,
    val total: Int,
    val page: Int,
    @SerialName("per_page") val perPage: Int
)

@Serializable
data class VoteRequest(
    // Backend expects "opcion_id" and "is_abstention"
    @SerialName("opcion_id") val optionId: String? = null,
    @SerialName("is_abstention") val isAbstention: Boolean = false
)

@Serializable
data class VoteResponse(
    val message: String
)

@Serializable
data class VotacionResultado(
    val votacion: Votacion,
    @SerialName("total_votos") val totalVotos: Int,
    @SerialName("total_abstenciones") val totalAbstenciones: Int,
    @SerialName("resultados") val resultados: List<OpcionResultado>,
    @SerialName("quorum_alcanzado") val quorumAlcanzado: Boolean,
    @SerialName("total_vecinos") val totalVecinos: Int,
    @SerialName("participacion") val participacion: Double
)

@Serializable
data class OpcionResultado(
    @SerialName("opcion_id") val opcionId: String,
    val label: String,
    val count: Int,
    val percentage: Double
)
