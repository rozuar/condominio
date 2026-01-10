package cl.vinapelvin.condominio.ui.viewmodel

import app.cash.turbine.test
import cl.vinapelvin.condominio.data.model.Emergencia
import cl.vinapelvin.condominio.data.model.EmergenciaListResponse
import cl.vinapelvin.condominio.data.repository.EmergenciaRepository
import cl.vinapelvin.condominio.data.repository.Result
import cl.vinapelvin.condominio.ui.emergencias.EmergenciasViewModel
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
class EmergenciasViewModelTest {

    private lateinit var repository: EmergenciaRepository
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
    fun `initial state shows loading then loads emergencias`() = runTest {
        // Given
        val emergencias = listOf(
            createTestEmergencia("1", "Corte de agua", "active"),
            createTestEmergencia("2", "Incidente resuelto", "resolved")
        )
        val response = EmergenciaListResponse(emergencias, 2, 1, 20)
        coEvery { repository.getEmergencias(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = EmergenciasViewModel(repository)

        // Then
        viewModel.uiState.test {
            val initialState = awaitItem()
            assertThat(initialState.isLoading).isTrue()

            testDispatcher.scheduler.advanceUntilIdle()

            val loadedState = awaitItem()
            assertThat(loadedState.isLoading).isFalse()
            assertThat(loadedState.emergencias).hasSize(2)
            assertThat(loadedState.error).isNull()
        }
    }

    @Test
    fun `loadEmergencias updates state with error on failure`() = runTest {
        // Given
        coEvery { repository.getEmergencias(any(), any()) } returns Result.Error("Connection failed")

        // When
        val viewModel = EmergenciasViewModel(repository)

        // Then
        viewModel.uiState.test {
            skipItems(1)
            testDispatcher.scheduler.advanceUntilIdle()

            val errorState = awaitItem()
            assertThat(errorState.isLoading).isFalse()
            assertThat(errorState.error).isEqualTo("Connection failed")
            assertThat(errorState.emergencias).isEmpty()
        }
    }

    @Test
    fun `empty emergencias list shows empty state`() = runTest {
        // Given
        val response = EmergenciaListResponse(emptyList(), 0, 1, 20)
        coEvery { repository.getEmergencias(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = EmergenciasViewModel(repository)
        testDispatcher.scheduler.advanceUntilIdle()

        // Then
        viewModel.uiState.test {
            val state = awaitItem()
            assertThat(state.emergencias).isEmpty()
            assertThat(state.isLoading).isFalse()
            assertThat(state.error).isNull()
        }
    }

    @Test
    fun `emergencias with different priorities are loaded correctly`() = runTest {
        // Given
        val emergencias = listOf(
            createTestEmergencia("1", "Critica", "active", "critical"),
            createTestEmergencia("2", "Alta", "active", "high"),
            createTestEmergencia("3", "Media", "active", "medium")
        )
        val response = EmergenciaListResponse(emergencias, 3, 1, 20)
        coEvery { repository.getEmergencias(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = EmergenciasViewModel(repository)
        testDispatcher.scheduler.advanceUntilIdle()

        // Then
        viewModel.uiState.test {
            val state = awaitItem()
            assertThat(state.emergencias).hasSize(3)
            assertThat(state.emergencias.map { it.priority }).containsExactly("critical", "high", "medium")
        }
    }

    private fun createTestEmergencia(
        id: String,
        title: String,
        status: String,
        priority: String = "high"
    ) = Emergencia(
        id = id,
        title = title,
        description = "Test description",
        priority = priority,
        status = status,
        reportedBy = "user1",
        reporterName = "Test User",
        resolvedAt = null,
        resolvedBy = null,
        createdAt = "2024-01-01T00:00:00Z",
        updatedAt = "2024-01-01T00:00:00Z"
    )
}
