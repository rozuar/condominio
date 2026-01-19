package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.GaleriaListResponse
import cl.vinapelvin.condominio.data.model.GaleriaWithItems
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GaleriaRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getGalerias(
        page: Int = 1,
        perPage: Int = 20,
        isPublic: Boolean? = null
    ): Result<GaleriaListResponse> {
        return try {
            val response = apiService.getGalerias(page, perPage, isPublic)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar galerías")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexión")
        }
    }

    suspend fun getGaleria(id: String): Result<GaleriaWithItems> {
        return try {
            val response = apiService.getGaleria(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar galería")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexión")
        }
    }
}
