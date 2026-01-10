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
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class VotacionDetailUiState(
    val isLoading: Boolean = false,
    val votacion: Votacion? = null,
    val selectedOptionId: String? = null,
    val isVoting: Boolean = false,
    val voteSuccess: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class VotacionDetailViewModel @Inject constructor(
    private val repository: VotacionRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(VotacionDetailUiState())
    val uiState: StateFlow<VotacionDetailUiState> = _uiState.asStateFlow()

    fun loadVotacion(id: String) {
        viewModelScope.launch {
            _uiState.value = VotacionDetailUiState(isLoading = true)

            when (val result = repository.getVotacion(id)) {
                is Result.Success -> {
                    _uiState.value = VotacionDetailUiState(votacion = result.data)
                }
                is Result.Error -> {
                    _uiState.value = VotacionDetailUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = VotacionDetailUiState(isLoading = true)
                }
            }
        }
    }

    fun selectOption(optionId: String) {
        _uiState.update { it.copy(selectedOptionId = optionId) }
    }

    fun submitVote(votacionId: String) {
        val optionId = _uiState.value.selectedOptionId ?: return

        viewModelScope.launch {
            _uiState.update { it.copy(isVoting = true) }

            when (val result = repository.votar(votacionId, optionId)) {
                is Result.Success -> {
                    _uiState.update {
                        it.copy(
                            votacion = result.data,
                            isVoting = false,
                            voteSuccess = true,
                            selectedOptionId = null
                        )
                    }
                }
                is Result.Error -> {
                    _uiState.update {
                        it.copy(isVoting = false, error = result.message)
                    }
                }
                is Result.Loading -> {}
            }
        }
    }
}
