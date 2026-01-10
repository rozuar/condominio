package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Comunicado
import cl.vinapelvin.condominio.data.model.ComunicadoListResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ComunicadoRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getComunicados(page: Int = 1, perPage: Int = 20): Result<ComunicadoListResponse> {
        return try {
            val response = apiService.getComunicados(page, perPage)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar comunicados")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getLatestComunicados(limit: Int = 5): Result<List<Comunicado>> {
        return try {
            val response = apiService.getLatestComunicados(limit)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar comunicados")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getComunicado(id: String): Result<Comunicado> {
        return try {
            val response = apiService.getComunicado(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar comunicado")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
