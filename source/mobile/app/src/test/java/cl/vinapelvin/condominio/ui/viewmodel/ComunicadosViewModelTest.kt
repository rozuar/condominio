package cl.vinapelvin.condominio.ui.viewmodel

import app.cash.turbine.test
import cl.vinapelvin.condominio.data.model.Comunicado
import cl.vinapelvin.condominio.data.model.ComunicadoListResponse
import cl.vinapelvin.condominio.data.repository.ComunicadoRepository
import cl.vinapelvin.condominio.data.repository.Result
import cl.vinapelvin.condominio.ui.comunicados.ComunicadosViewModel
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ComunicadosViewModelTest {

    private lateinit var repository: ComunicadoRepository
    private val testDispatcher = StandardTestDispatcher()

    @Before
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        repository = mockk()
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `initial state shows loading then loads comunicados`() = runTest {
        // Given
        val comunicados = listOf(createTestComunicado("1", "Test"))
        val response = ComunicadoListResponse(comunicados, 1, 1, 20)
        coEvery { repository.getComunicados(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = ComunicadosViewModel(repository)

        // Then
        viewModel.uiState.test {
            val initialState = awaitItem()
            assertThat(initialState.isLoading).isTrue()

            testDispatcher.scheduler.advanceUntilIdle()

            val loadedState = awaitItem()
            assertThat(loadedState.isLoading).isFalse()
            assertThat(loadedState.comunicados).hasSize(1)
            assertThat(loadedState.error).isNull()
        }
    }

    @Test
    fun `loadComunicados updates state with error on failure`() = runTest {
        // Given
        coEvery { repository.getComunicados(any(), any()) } returns Result.Error("Network error")

        // When
        val viewModel = ComunicadosViewModel(repository)

        // Then
        viewModel.uiState.test {
            skipItems(1) // Skip loading state
            testDispatcher.scheduler.advanceUntilIdle()

            val errorState = awaitItem()
            assertThat(errorState.isLoading).isFalse()
            assertThat(errorState.error).isEqualTo("Network error")
            assertThat(errorState.comunicados).isEmpty()
        }
    }

    @Test
    fun `loadComunicados can be called to refresh`() = runTest {
        // Given
        val comunicados = listOf(createTestComunicado("1", "Test"))
        val response = ComunicadoListResponse(comunicados, 1, 1, 20)
        coEvery { repository.getComunicados(any(), any()) } returns Result.Success(response)

        val viewModel = ComunicadosViewModel(repository)
        testDispatcher.scheduler.advanceUntilIdle()

        // When
        viewModel.loadComunicados()
        testDispatcher.scheduler.advanceUntilIdle()

        // Then
        viewModel.uiState.test {
            val state = awaitItem()
            assertThat(state.comunicados).hasSize(1)
        }
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
