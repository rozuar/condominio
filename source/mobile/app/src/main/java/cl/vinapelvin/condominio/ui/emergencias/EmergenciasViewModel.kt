package cl.vinapelvin.condominio.ui.emergencias

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Emergencia
import cl.vinapelvin.condominio.data.repository.EmergenciaRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class EmergenciasUiState(
    val isLoading: Boolean = false,
    val emergencias: List<Emergencia> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class EmergenciasViewModel @Inject constructor(
    private val repository: EmergenciaRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(EmergenciasUiState(isLoading = true))
    val uiState: StateFlow<EmergenciasUiState> = _uiState.asStateFlow()

    init {
        loadEmergencias()
    }

    fun loadEmergencias() {
        viewModelScope.launch {
            _uiState.value = EmergenciasUiState(isLoading = true)

            when (val result = repository.getEmergencias()) {
                is Result.Success -> {
                    _uiState.value = EmergenciasUiState(emergencias = result.data.emergencias)
                }
                is Result.Error -> {
                    _uiState.value = EmergenciasUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = EmergenciasUiState(isLoading = true)
                }
            }
        }
    }
}
