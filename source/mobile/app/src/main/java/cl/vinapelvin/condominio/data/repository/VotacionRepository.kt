package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.*
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class VotacionRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getVotaciones(page: Int = 1, perPage: Int = 20, status: String? = null): Result<VotacionListResponse> {
        return try {
            val response = apiService.getVotaciones(page, perPage, status)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar votaciones")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getActiveVotaciones(): Result<List<Votacion>> {
        return try {
            val response = apiService.getActiveVotaciones()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar votaciones activas")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getVotacion(id: String): Result<Votacion> {
        return try {
            val response = apiService.getVotacion(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar votacion")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun votar(votacionId: String, optionId: String): Result<Votacion> {
        return try {
            val response = apiService.votar(votacionId, VoteRequest(optionId))
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al emitir voto")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getResultados(id: String): Result<VotacionResultado> {
        return try {
            val response = apiService.getVotacionResultados(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar resultados")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
