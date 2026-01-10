package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Evento
import cl.vinapelvin.condominio.data.model.EventoListResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class EventoRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getEventos(page: Int = 1, perPage: Int = 20): Result<EventoListResponse> {
        return try {
            val response = apiService.getEventos(page, perPage)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar eventos")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getUpcomingEventos(limit: Int = 5): Result<List<Evento>> {
        return try {
            val response = apiService.getUpcomingEventos(limit)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar eventos")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getEvento(id: String): Result<Evento> {
        return try {
            val response = apiService.getEvento(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar evento")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
