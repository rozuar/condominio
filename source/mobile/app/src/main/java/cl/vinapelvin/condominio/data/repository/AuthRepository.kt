package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.local.TokenManager
import cl.vinapelvin.condominio.data.model.LoginRequest
import cl.vinapelvin.condominio.data.model.RefreshRequest
import cl.vinapelvin.condominio.data.model.User
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject
import javax.inject.Singleton

sealed class Result<out T> {
    data class Success<T>(val data: T) : Result<T>()
    data class Error(val message: String) : Result<Nothing>()
    object Loading : Result<Nothing>()
}

@Singleton
class AuthRepository @Inject constructor(
    private val apiService: ApiService,
    private val tokenManager: TokenManager
) {
    val isLoggedIn: Flow<Boolean> = tokenManager.isLoggedIn
    val userName: Flow<String?> = tokenManager.userName
    val userEmail: Flow<String?> = tokenManager.userEmail
    val userRole: Flow<String?> = tokenManager.userRole

    suspend fun login(email: String, password: String): Result<User> {
        return try {
            val response = apiService.login(LoginRequest(email, password))
            if (response.isSuccessful) {
                val authResponse = response.body()!!
                tokenManager.saveTokens(authResponse.accessToken, authResponse.refreshToken)
                tokenManager.saveUser(
                    authResponse.user.id,
                    authResponse.user.email,
                    authResponse.user.name,
                    authResponse.user.role,
                    authResponse.user.parcelaId
                )
                Result.Success(authResponse.user)
            } else {
                val errorBody = response.errorBody()?.string()
                Result.Error(errorBody ?: "Error de autenticacion")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun refreshToken(): Result<Unit> {
        return try {
            val refreshToken = tokenManager.getRefreshToken()
                ?: return Result.Error("No hay token de refresco")

            val response = apiService.refreshToken(RefreshRequest(refreshToken))
            if (response.isSuccessful) {
                val authResponse = response.body()!!
                tokenManager.saveTokens(authResponse.accessToken, authResponse.refreshToken)
                Result.Success(Unit)
            } else {
                Result.Error("Error al refrescar token")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun getMe(): Result<User> {
        return try {
            val response = apiService.getMe()
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Error al obtener usuario")
            }
        } catch (e: Exception) {
            Result.Error(e.message ?: "Error de conexion")
        }
    }

    suspend fun logout() {
        tokenManager.clearAll()
    }
}
