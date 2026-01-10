package cl.vinapelvin.condominio.data.api

import cl.vinapelvin.condominio.data.local.TokenManager
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject

class AuthInterceptor @Inject constructor(
    private val tokenManager: TokenManager
) : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()

        // Skip auth header for login and refresh endpoints
        if (originalRequest.url.encodedPath.contains("auth/login") ||
            originalRequest.url.encodedPath.contains("auth/refresh")) {
            return chain.proceed(originalRequest)
        }

        val token = runBlocking { tokenManager.getAccessToken() }

        return if (token != null) {
            val newRequest = originalRequest.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
            chain.proceed(newRequest)
        } else {
            chain.proceed(originalRequest)
        }
    }
}
