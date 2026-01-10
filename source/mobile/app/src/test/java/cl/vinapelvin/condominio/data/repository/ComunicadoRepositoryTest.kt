package cl.vinapelvin.condominio.data.repository

import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.model.Comunicado
import cl.vinapelvin.condominio.data.model.ComunicadoListResponse
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Before
import org.junit.Test
import retrofit2.Response

class ComunicadoRepositoryTest {

    private lateinit var apiService: ApiService
    private lateinit var repository: ComunicadoRepository

    @Before
    fun setup() {
        apiService = mockk()
        repository = ComunicadoRepository(apiService)
    }

    @Test
    fun `getComunicados returns success when API call is successful`() = runTest {
        // Given
        val comunicados = listOf(
            createTestComunicado("1", "Titulo 1"),
            createTestComunicado("2", "Titulo 2")
        )
        val response = ComunicadoListResponse(
            comunicados = comunicados,
            total = 2,
            page = 1,
            perPage = 20
        )
        coEvery { apiService.getComunicados(any(), any()) } returns Response.success(response)

        // When
        val result = repository.getComunicados()

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.comunicados).hasSize(2)
        assertThat(successResult.data.comunicados[0].title).isEqualTo("Titulo 1")
    }

    @Test
    fun `getComunicados returns error when API call fails`() = runTest {
        // Given
        coEvery { apiService.getComunicados(any(), any()) } returns Response.error(
            500,
            "Server error".toResponseBody()
        )

        // When
        val result = repository.getComunicados()

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
        val errorResult = result as Result.Error
        assertThat(errorResult.message).isEqualTo("Error al cargar comunicados")
    }

    @Test
    fun `getComunicados returns error when exception is thrown`() = runTest {
        // Given
        coEvery { apiService.getComunicados(any(), any()) } throws Exception("Network error")

        // When
        val result = repository.getComunicados()

        // Then
        assertThat(result).isInstanceOf(Result.Error::class.java)
        val errorResult = result as Result.Error
        assertThat(errorResult.message).isEqualTo("Network error")
    }

    @Test
    fun `getComunicado returns success when API call is successful`() = runTest {
        // Given
        val comunicado = createTestComunicado("1", "Test Comunicado")
        coEvery { apiService.getComunicado("1") } returns Response.success(comunicado)

        // When
        val result = repository.getComunicado("1")

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data.id).isEqualTo("1")
        assertThat(successResult.data.title).isEqualTo("Test Comunicado")
    }

    @Test
    fun `getLatestComunicados returns success with limited results`() = runTest {
        // Given
        val comunicados = listOf(
            createTestComunicado("1", "Latest 1"),
            createTestComunicado("2", "Latest 2"),
            createTestComunicado("3", "Latest 3")
        )
        coEvery { apiService.getLatestComunicados(3) } returns Response.success(comunicados)

        // When
        val result = repository.getLatestComunicados(3)

        // Then
        assertThat(result).isInstanceOf(Result.Success::class.java)
        val successResult = result as Result.Success
        assertThat(successResult.data).hasSize(3)
    }

    private fun createTestComunicado(id: String, title: String) = Comunicado(
        id = id,
        title = title,
        content = "Test content",
        type = "general",
        priority = "normal",
        authorId = "author1",
        authorName = "Test Author",
        createdAt = "2024-01-01T00:00:00Z",
        updatedAt = "2024-01-01T00:00:00Z"
    )
}
