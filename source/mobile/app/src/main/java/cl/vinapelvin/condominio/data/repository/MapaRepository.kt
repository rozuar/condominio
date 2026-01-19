package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.MapaData
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class MapaRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getMapaData(): Result<MapaData> {
        return try {
            val response = apiService.getMapaData()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar datos del mapa")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexión")
        }
    }
}
