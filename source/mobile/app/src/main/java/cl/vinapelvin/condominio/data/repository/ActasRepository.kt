package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Acta
import cl.vinapelvin.condominio.data.model.ActaListResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ActasRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getActas(
        page: Int = 1,
        perPage: Int = 20,
        type: String? = null,
        year: Int? = null
    ): Result<ActaListResponse> {
        return try {
            val response = apiService.getActas(page, perPage, type, year)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar actas")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getActa(id: String): Result<Acta> {
        return try {
            val response = apiService.getActa(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar acta")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}

