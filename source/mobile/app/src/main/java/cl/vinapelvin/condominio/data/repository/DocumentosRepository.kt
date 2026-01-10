package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Documento
import cl.vinapelvin.condominio.data.model.DocumentoListResponse
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class DocumentosRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getDocumentos(
        page: Int = 1,
        perPage: Int = 20,
        category: String? = null
    ): Result<DocumentoListResponse> {
        return try {
            val response = apiService.getDocumentos(page, perPage, category)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar documentos")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getDocumento(id: String): Result<Documento> {
        return try {
            val response = apiService.getDocumento(id)
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al cargar documento")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }
}

