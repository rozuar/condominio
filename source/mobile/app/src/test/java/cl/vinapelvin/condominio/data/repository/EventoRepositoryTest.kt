package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Evento
import cl.vinapelvin.condominio.data.model.EventoListResponse
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Before
import org.junit.Test
import retrofit2.Response

class EventoRepositoryTest {

    private lateinit var apiService: ApiService
    private lateinit var repository: EventoRepository

    @Before
    fun setup() {
        apiService = mockk()
        repository = EventoRepository(apiService)
    }

    @Test
    fun `getEventos returns success when API call is successful`() = runTest {
        // Given
        val eventos = listOf(
            createTestEvento("1", "Asamblea General"),
            createTestEvento("2", "Reunion Directiva")
        )
        val response = EventoListResponse(
            eventos = eventos,
            total = 2,
            page = 1,
            perPage = 20
        )
        coEvery { apiService.getEventos(any(), any()) } returns Response.success(response)

        // When
        val result = repository.getEventos()

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.eventos).hasSize(2)
    }

    @Test
    fun `getEventos returns error when API call fails`() = runTest {
        // Given
        coEvery { apiService.getEventos(any(), any()) } returns Response.error(
            500,
            "Server error".toResponseBody()
        )

        // When
        val result = repository.getEventos()

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
    }

    @Test
    fun `getUpcomingEventos returns success`() = runTest {
        // Given
        val eventos = listOf(createTestEvento("1", "Proximo Evento"))
        coEvery { apiService.getUpcomingEventos(any()) } returns Response.success(eventos)

        // When
        val result = repository.getUpcomingEventos(5)

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data).hasSize(1)
    }

    @Test
    fun `getEvento returns success when API call is successful`() = runTest {
        // Given
        val evento = createTestEvento("1", "Evento Test")
        coEvery { apiService.getEvento("1") } returns Response.success(evento)

        // When
        val result = repository.getEvento("1")

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.id).isEqualTo("1")
    }

    @Test
    fun `getEvento returns error when exception is thrown`() = runTest {
        // Given
        coEvery { apiService.getEvento(any()) } throws Exception("Connection failed")

        // When
        val result = repository.getEvento("1")

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
        val errorResult = result as Result.Error
        assertThat(errorResult.message).isEqualTo("Connection failed")
    }

    private fun createTestEvento(id: String, title: String) = Evento(
        id = id,
        title = title,
        description = "Test description",
        location = "Sede Social",
        startDate = "2024-01-15T10:00:00Z",
        endDate = "2024-01-15T12:00:00Z",
        type = "asamblea",
        isMandatory = false,
        createdAt = "2024-01-01T00:00:00Z"
    )
}
