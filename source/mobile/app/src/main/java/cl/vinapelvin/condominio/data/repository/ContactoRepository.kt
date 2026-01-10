package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.CreateMensajeContactoRequest
import cl.vinapelvin.condominio.data.model.MensajeContacto
import cl.vinapelvin.condominio.data.model.MisMensajesResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ContactoRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun enviarMensajeContacto(
        nombre: String,
        email: String,
        asunto: String,
        mensaje: String
    ): Result<MensajeContacto> {
        return try {
            val response = apiService.enviarMensajeContacto(
                CreateMensajeContactoRequest(
                    nombre = nombre,
                    email = email,
                    asunto = asunto,
                    mensaje = mensaje
                )
            )
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al enviar mensaje")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getMisMensajes(): Result<MisMensajesResponse> {
        return try {
            val response = apiService.getMisMensajes()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar mis mensajes")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}

