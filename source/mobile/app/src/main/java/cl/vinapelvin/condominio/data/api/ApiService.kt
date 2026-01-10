package cl.vinapelvin.condominio.data.api

import cl.vinapelvin.condominio.data.model.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {

    // Auth
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body request: RefreshRequest): Response<AuthResponse>

    @GET("auth/me")
    suspend fun getMe(): Response<User>

    // Comunicados
    @GET("api/v1/comunicados")
    suspend fun getComunicados(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20
    ): Response<ComunicadoListResponse>

    @GET("api/v1/comunicados/latest")
    suspend fun getLatestComunicados(
        @Query("limit") limit: Int = 5
    ): Response<List<Comunicado>>

    @GET("api/v1/comunicados/{id}")
    suspend fun getComunicado(@Path("id") id: String): Response<Comunicado>

    // Eventos
    @GET("api/v1/eventos")
    suspend fun getEventos(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20
    ): Response<EventoListResponse>

    @GET("api/v1/eventos/upcoming")
    suspend fun getUpcomingEventos(
        @Query("limit") limit: Int = 5
    ): Response<List<Evento>>

    @GET("api/v1/eventos/{id}")
    suspend fun getEvento(@Path("id") id: String): Response<Evento>

    // Emergencias
    @GET("api/v1/emergencias")
    suspend fun getEmergencias(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20
    ): Response<EmergenciaListResponse>

    @GET("api/v1/emergencias/active")
    suspend fun getActiveEmergencias(): Response<List<Emergencia>>

    @GET("api/v1/emergencias/{id}")
    suspend fun getEmergencia(@Path("id") id: String): Response<Emergencia>

    // Votaciones
    @GET("api/v1/votaciones")
    suspend fun getVotaciones(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20,
        @Query("status") status: String? = null
    ): Response<VotacionListResponse>

    @GET("api/v1/votaciones/active")
    suspend fun getActiveVotaciones(): Response<List<Votacion>>

    @GET("api/v1/votaciones/{id}")
    suspend fun getVotacion(@Path("id") id: String): Response<Votacion>

    @GET("api/v1/votaciones/{id}/resultados")
    suspend fun getVotacionResultados(@Path("id") id: String): Response<VotacionResultado>

    @POST("api/v1/votaciones/{id}/votar")
    suspend fun votar(
        @Path("id") id: String,
        @Body request: VoteRequest
    ): Response<Votacion>

    // Gastos
    @GET("api/v1/gastos/mi-cuenta")
    suspend fun getMiEstadoCuenta(): Response<EstadoCuenta>

    // Tesoreria
    @GET("api/v1/tesoreria/resumen")
    suspend fun getTesoreriaResumen(): Response<TesoreriaResumen>

    @GET("api/v1/tesoreria")
    suspend fun getMovimientos(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20,
        @Query("type") type: String? = null,
        @Query("year") year: Int? = null,
        @Query("month") month: Int? = null
    ): Response<MovimientoListResponse>

    // Actas
    @GET("api/v1/actas")
    suspend fun getActas(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20,
        @Query("type") type: String? = null,
        @Query("year") year: Int? = null
    ): Response<ActaListResponse>

    @GET("api/v1/actas/{id}")
    suspend fun getActa(@Path("id") id: String): Response<Acta>

    // Documentos
    @GET("api/v1/documentos")
    suspend fun getDocumentos(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20,
        @Query("category") category: String? = null
    ): Response<DocumentoListResponse>

    @GET("api/v1/documentos/{id}")
    suspend fun getDocumento(@Path("id") id: String): Response<Documento>

    // Notificaciones
    @GET("api/v1/notificaciones")
    suspend fun getNotificaciones(
        @Query("page") page: Int = 1,
        @Query("per_page") perPage: Int = 20
    ): Response<NotificacionListResponse>

    @GET("api/v1/notificaciones/stats")
    suspend fun getNotificacionStats(): Response<NotificacionStats>

    @POST("api/v1/notificaciones/{id}/read")
    suspend fun markNotificacionAsRead(@Path("id") id: String): Response<Notificacion>

    @POST("api/v1/notificaciones/read-all")
    suspend fun markAllNotificacionesAsRead(): Response<Unit>

    // Contacto
    @POST("api/v1/contacto")
    suspend fun enviarMensajeContacto(
        @Body request: CreateMensajeContactoRequest
    ): Response<MensajeContacto>

    @GET("api/v1/contacto/mis-mensajes")
    suspend fun getMisMensajes(): Response<MisMensajesResponse>
}
