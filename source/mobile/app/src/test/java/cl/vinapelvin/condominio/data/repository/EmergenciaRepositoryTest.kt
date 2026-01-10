package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Emergencia
import cl.vinapelvin.condominio.data.model.EmergenciaListResponse
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Before
import org.junit.Test
import retrofit2.Response

class EmergenciaRepositoryTest {

    private lateinit var apiService: ApiService
    private lateinit var repository: EmergenciaRepository

    @Before
    fun setup() {
        apiService = mockk()
        repository = EmergenciaRepository(apiService)
    }

    @Test
    fun `getEmergencias returns success when API call is successful`() = runTest {
        // Given
        val emergencias = listOf(
            createTestEmergencia("1", "Corte de agua", "active"),
            createTestEmergencia("2", "Robo reportado", "resolved")
        )
        val response = EmergenciaListResponse(
            emergencias = emergencias,
            total = 2,
            page = 1,
            perPage = 20
        )
        coEvery { apiService.getEmergencias(any(), any()) } returns Response.success(response)

        // When
        val result = repository.getEmergencias()

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.emergencias).hasSize(2)
    }

    @Test
    fun `getEmergencias returns error when API call fails`() = runTest {
        // Given
        coEvery { apiService.getEmergencias(any(), any()) } returns Response.error(
            500,
            "Server error".toResponseBody()
        )

        // When
        val result = repository.getEmergencias()

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
        val errorResult = result as Result.Error
        assertThat(errorResult.message).isEqualTo("Error al cargar emergencias")
    }

    @Test
    fun `getActiveEmergencias returns only active emergencias`() = runTest {
        // Given
        val activeEmergencias = listOf(
            createTestEmergencia("1", "Emergencia Activa", "active")
        )
        coEvery { apiService.getActiveEmergencias() } returns Response.success(activeEmergencias)

        // When
        val result = repository.getActiveEmergencias()

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data).hasSize(1)
        assertThat(successResult.data[0].status).isEqualTo("active")
    }

    @Test
    fun `getEmergencia returns success when API call is successful`() = runTest {
        // Given
        val emergencia = createTestEmergencia("1", "Test Emergencia", "active")
        coEvery { apiService.getEmergencia("1") } returns Response.success(emergencia)

        // When
        val result = repository.getEmergencia("1")

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.id).isEqualTo("1")
    }

    @Test
    fun `getEmergencia returns error when exception is thrown`() = runTest {
        // Given
        coEvery { apiService.getEmergencia(any()) } throws Exception("Network error")

        // When
        val result = repository.getEmergencia("1")

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
        val errorResult = result as Result.Error
        assertThat(errorResult.message).isEqualTo("Network error")
    }

    private fun createTestEmergencia(id: String, title: String, status: String) = Emergencia(
        id = id,
        title = title,
        description = "Test description",
        priority = "high",
        status = status,
        reportedBy = "user1",
        reporterName = "Test User",
        resolvedAt = null,
        resolvedBy = null,
        createdAt = "2024-01-01T00:00:00Z",
        updatedAt = "2024-01-01T00:00:00Z"
    )
}
