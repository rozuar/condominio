package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.MiEstadoCuenta
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GastosRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getMiEstadoCuenta(): Result<MiEstadoCuenta> {
        return try {
            val response = apiService.getMiEstadoCuenta()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                val code = response.code()
                val body = response.errorBody()?.string().orEmpty()
                when (code) {
                    401 -> Result.Error("Sesión expirada. Vuelve a iniciar sesión.")
                    403 -> Result.Error("No tienes permisos para ver Gastos (solo vecino/directiva).")
                    else -> {
                        val msg = when {
                            body.contains("user has no associated parcela", ignoreCase = true) ->
                                "Ud no posee asociada una parcela que genere gastos"
                            body.contains("Ud no posee asociada una parcela", ignoreCase = true) ->
                                "Ud no posee asociada una parcela que genere gastos"
                            else -> body
                        }.ifBlank { "Error al cargar estado de cuenta" }
                        Result.Error(msg)
                    }
                }
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}
