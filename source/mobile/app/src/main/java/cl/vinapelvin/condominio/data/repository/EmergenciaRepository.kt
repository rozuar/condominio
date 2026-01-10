package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Emergencia
import cl.vinapelvin.condominio.data.model.EmergenciaListResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class EmergenciaRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getEmergencias(page: Int = 1, perPage: Int = 20): Result<EmergenciaListResponse> {
        return try {
            val response = apiService.getEmergencias(page, perPage)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar emergencias")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getActiveEmergencias(): Result<List<Emergencia>> {
        return try {
            val response = apiService.getActiveEmergencias()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar emergencias activas")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getEmergencia(id: String): Result<Emergencia> {
        return try {
            val response = apiService.getEmergencia(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar emergencia")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
