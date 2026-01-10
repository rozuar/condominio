package cl.vinapelvin.condominio.ui.viewmodel

import app.cash.turbine.test
import cl.vinapelvin.condominio.data.model.Evento
import cl.vinapelvin.condominio.data.model.EventoListResponse
import cl.vinapelvin.condominio.data.repository.EventoRepository
import cl.vinapelvin.condominio.data.repository.Result
import cl.vinapelvin.condominio.ui.eventos.EventosViewModel
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
class EventosViewModelTest {

    private lateinit var repository: EventoRepository
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
    fun `initial state shows loading then loads eventos`() = runTest {
        // Given
        val eventos = listOf(
            createTestEvento("1", "Asamblea"),
            createTestEvento("2", "Reunion")
        )
        val response = EventoListResponse(eventos, 2, 1, 20)
        coEvery { repository.getEventos(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = EventosViewModel(repository)

        // Then
        viewModel.uiState.test {
            val initialState = awaitItem()
            assertThat(initialState.isLoading).isTrue()

            testDispatcher.scheduler.advanceUntilIdle()

            val loadedState = awaitItem()
            assertThat(loadedState.isLoading).isFalse()
            assertThat(loadedState.eventos).hasSize(2)
            assertThat(loadedState.error).isNull()
        }
    }

    @Test
    fun `loadEventos updates state with error on failure`() = runTest {
        // Given
        coEvery { repository.getEventos(any(), any()) } returns Result.Error("Failed to load")

        // When
        val viewModel = EventosViewModel(repository)

        // Then
        viewModel.uiState.test {
            skipItems(1)
            testDispatcher.scheduler.advanceUntilIdle()

            val errorState = awaitItem()
            assertThat(errorState.isLoading).isFalse()
            assertThat(errorState.error).isEqualTo("Failed to load")
            assertThat(errorState.eventos).isEmpty()
        }
    }

    @Test
    fun `eventos are sorted and accessible`() = runTest {
        // Given
        val eventos = listOf(
            createTestEvento("1", "Evento A"),
            createTestEvento("2", "Evento B"),
            createTestEvento("3", "Evento C")
        )
        val response = EventoListResponse(eventos, 3, 1, 20)
        coEvery { repository.getEventos(any(), any()) } returns Result.Success(response)

        // When
        val viewModel = EventosViewModel(repository)
        testDispatcher.scheduler.advanceUntilIdle()

        // Then
        viewModel.uiState.test {
            val state = awaitItem()
            assertThat(state.eventos).hasSize(3)
            assertThat(state.eventos.map { it.title }).containsExactly("Evento A", "Evento B", "Evento C")
        }
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
