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

data class ActaDetailUiState(
    val isLoading: Boolean = false,
    val acta: Acta? = null,
    val error: String? = null
)

@HiltViewModel
class ActaDetailViewModel @Inject constructor(
    private val repository: ActasRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ActaDetailUiState(isLoading = true))
    val uiState: StateFlow<ActaDetailUiState> = _uiState.asStateFlow()

    fun loadActa(id: String) {
        viewModelScope.launch {
            _uiState.value = ActaDetailUiState(isLoading = true)
            when (val result = repository.getActa(id)) {
                is Result.Success -> _uiState.value = ActaDetailUiState(acta = result.data)
                is Result.Error -> _uiState.value = ActaDetailUiState(error = result.message)
                is Result.Loading -> _uiState.value = ActaDetailUiState(isLoading = true)
            }
        }
    }
}

