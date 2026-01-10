package cl.vinapelvin.condominio.ui.notificaciones

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Notificacion
import cl.vinapelvin.condominio.data.repository.NotificacionRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class NotificacionesUiState(
    val isLoading: Boolean = false,
    val notificaciones: List<Notificacion> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class NotificacionesViewModel @Inject constructor(
    private val repository: NotificacionRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(NotificacionesUiState())
    val uiState: StateFlow<NotificacionesUiState> = _uiState.asStateFlow()

    init {
        loadNotificaciones()
    }

    fun loadNotificaciones() {
        viewModelScope.launch {
            _uiState.value = NotificacionesUiState(isLoading = true)

            when (val result = repository.getNotificaciones()) {
                is Result.Success -> {
                    _uiState.value = NotificacionesUiState(
                        notificaciones = result.data.notificaciones
                    )
                }
                is Result.Error -> {
                    _uiState.value = NotificacionesUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = NotificacionesUiState(isLoading = true)
                }
            }
        }
    }

    fun markAsRead(id: String) {
        viewModelScope.launch {
            when (repository.markAsRead(id)) {
                is Result.Success -> {
                    _uiState.update { state ->
                        state.copy(
                            notificaciones = state.notificaciones.map {
                                if (it.id == id) it.copy(isRead = true) else it
                            }
                        )
                    }
                }
                else -> {}
            }
        }
    }

    fun markAllAsRead() {
        viewModelScope.launch {
            when (repository.markAllAsRead()) {
                is Result.Success -> {
                    _uiState.update { state ->
                        state.copy(
                            notificaciones = state.notificaciones.map { it.copy(isRead = true) }
                        )
                    }
                }
                else -> {}
            }
        }
    }
}
