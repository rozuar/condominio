package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Notificacion
import cl.vinapelvin.condominio.data.model.NotificacionListResponse
import cl.vinapelvin.condominio.data.model.NotificacionStats
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class NotificacionRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getNotificaciones(page: Int = 1, perPage: Int = 20): Result<NotificacionListResponse> {
        return try {
            val response = apiService.getNotificaciones(page, perPage)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar notificaciones")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getStats(): Result<NotificacionStats> {
        return try {
            val response = apiService.getNotificacionStats()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar estadisticas")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun markAsRead(id: String): Result<Notificacion> {
        return try {
            val response = apiService.markNotificacionAsRead(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al marcar como leida")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun markAllAsRead(): Result<Unit> {
        return try {
            val response = apiService.markAllNotificacionesAsRead()
            if (response.isSuccessful) {
                Result.Success(Unit)
            } else {
                Result.Error("Error al marcar todas como leidas")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
