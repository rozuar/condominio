package cl.vinapelvin.condominio.ui.votaciones

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Votacion
import cl.vinapelvin.condominio.data.repository.VotacionRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class VotacionesUiState(
    val isLoading: Boolean = false,
    val votaciones: List<Votacion> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class VotacionesViewModel @Inject constructor(
    private val repository: VotacionRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(VotacionesUiState())
    val uiState: StateFlow<VotacionesUiState> = _uiState.asStateFlow()

    init {
        loadVotaciones()
    }

    fun loadVotaciones() {
        viewModelScope.launch {
            _uiState.value = VotacionesUiState(isLoading = true)

            when (val result = repository.getVotaciones()) {
                is Result.Success -> {
                    _uiState.value = VotacionesUiState(votaciones = result.data.votaciones)
                }
                is Result.Error -> {
                    _uiState.value = VotacionesUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = VotacionesUiState(isLoading = true)
                }
            }
        }
    }
}
