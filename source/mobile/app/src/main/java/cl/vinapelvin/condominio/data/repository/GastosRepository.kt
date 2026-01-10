package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.EstadoCuenta
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GastosRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getMiEstadoCuenta(): Result<EstadoCuenta> {
        return try {
            val response = apiService.getMiEstadoCuenta()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar estado de cuenta")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
