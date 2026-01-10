package cl.vinapelvin.condominio.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class VotacionOpcion(
    val id: String,
    val text: String,
    val votes: Int = 0
)

@Serializable
data class Votacion(
    val id: String,
    val title: String,
    val description: String,
    val options: List<VotacionOpcion>,
    @SerialName("start_date") val startDate: String,
    @SerialName("end_date") val endDate: String,
    val status: String,
    @SerialName("requires_quorum") val requiresQuorum: Boolean,
    @SerialName("quorum_percentage") val quorumPercentage: Int,
    @SerialName("user_voted") val userVoted: Boolean = false,
    @SerialName("user_vote") val userVote: String? = null,
    @SerialName("created_at") val createdAt: String
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
    @SerialName("option_id") val optionId: String
)

@Serializable
data class VotacionResultado(
    @SerialName("votacion_id") val votacionId: String,
    @SerialName("total_votes") val totalVotes: Int,
    @SerialName("total_eligible") val totalEligible: Int,
    val participation: Double,
    @SerialName("quorum_met") val quorumMet: Boolean,
    val results: List<VotacionOpcion>
)
