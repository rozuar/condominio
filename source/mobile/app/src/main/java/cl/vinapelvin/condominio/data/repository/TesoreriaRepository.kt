package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.MovimientoListResponse
import cl.vinapelvin.condominio.data.model.TesoreriaResumen
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TesoreriaRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getResumen(): Result<TesoreriaResumen> {
        return try {
            val response = apiService.getTesoreriaResumen()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar resumen")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getMovimientos(
        page: Int = 1,
        perPage: Int = 20,
        type: String? = null,
        year: Int? = null,
        month: Int? = null
    ): Result<MovimientoListResponse> {
        return try {
            val response = apiService.getMovimientos(page, perPage, type, year, month)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar movimientos")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}

