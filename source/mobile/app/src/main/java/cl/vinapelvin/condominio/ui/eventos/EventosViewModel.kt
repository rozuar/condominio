package cl.vinapelvin.condominio.ui.eventos

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

data class EventosUiState(
    val isLoading: Boolean = false,
    val eventos: List<Evento> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class EventosViewModel @Inject constructor(
    private val repository: EventoRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(EventosUiState(isLoading = true))
    val uiState: StateFlow<EventosUiState> = _uiState.asStateFlow()

    init {
        loadEventos()
    }

    fun loadEventos() {
        viewModelScope.launch {
            _uiState.value = EventosUiState(isLoading = true)

            when (val result = repository.getEventos()) {
                is Result.Success -> {
                    _uiState.value = EventosUiState(eventos = result.data.eventos)
                }
                is Result.Error -> {
                    _uiState.value = EventosUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = EventosUiState(isLoading = true)
                }
            }
        }
    }
}
