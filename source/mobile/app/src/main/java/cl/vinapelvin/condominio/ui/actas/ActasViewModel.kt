package cl.vinapelvin.condominio.ui.actas

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Acta
import cl.vinapelvin.condominio.data.repository.ActasRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ActasUiState(
    val isLoading: Boolean = false,
    val actas: List<Acta> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class ActasViewModel @Inject constructor(
    private val repository: ActasRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ActasUiState(isLoading = true))
    val uiState: StateFlow<ActasUiState> = _uiState.asStateFlow()

    init {
        loadActas()
    }

    fun loadActas() {
        viewModelScope.launch {
            _uiState.value = ActasUiState(isLoading = true)
            when (val result = repository.getActas()) {
                is Result.Success -> _uiState.value = ActasUiState(actas = result.data.actas)
                is Result.Error -> _uiState.value = ActasUiState(error = result.message)
                is Result.Loading -> _uiState.value = ActasUiState(isLoading = true)
            }
        }
    }
}

