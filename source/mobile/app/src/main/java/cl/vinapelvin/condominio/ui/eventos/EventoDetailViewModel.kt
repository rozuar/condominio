package cl.vinapelvin.condominio.ui.eventos

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Evento
import cl.vinapelvin.condominio.data.repository.EventoRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class EventoDetailUiState(
    val isLoading: Boolean = false,
    val evento: Evento? = null,
    val error: String? = null
)

@HiltViewModel
class EventoDetailViewModel @Inject constructor(
    private val repository: EventoRepository,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val eventoId: String = savedStateHandle.get<String>("id") ?: ""

    private val _uiState = MutableStateFlow(EventoDetailUiState())
    val uiState: StateFlow<EventoDetailUiState> = _uiState.asStateFlow()

    init {
        loadEvento()
    }

    fun loadEvento() {
        if (eventoId.isBlank()) {
            _uiState.value = EventoDetailUiState(error = "ID de evento invalido")
            return
        }

        viewModelScope.launch {
            _uiState.value = EventoDetailUiState(isLoading = true)

            when (val result = repository.getEvento(eventoId)) {
                is Result.Success -> {
                    _uiState.value = EventoDetailUiState(evento = result.data)
                }
                is Result.Error -> {
                    _uiState.value = EventoDetailUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = EventoDetailUiState(isLoading = true)
                }
            }
        }
    }
}
