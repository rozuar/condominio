package cl.vinapelvin.condominio.ui.comunicados

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Comunicado
import cl.vinapelvin.condominio.data.repository.ComunicadoRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ComunicadoDetailUiState(
    val isLoading: Boolean = false,
    val comunicado: Comunicado? = null,
    val error: String? = null
)

@HiltViewModel
class ComunicadoDetailViewModel @Inject constructor(
    private val repository: ComunicadoRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ComunicadoDetailUiState())
    val uiState: StateFlow<ComunicadoDetailUiState> = _uiState.asStateFlow()

    fun loadComunicado(id: String) {
        viewModelScope.launch {
            _uiState.value = ComunicadoDetailUiState(isLoading = true)

            when (val result = repository.getComunicado(id)) {
                is Result.Success -> {
                    _uiState.value = ComunicadoDetailUiState(comunicado = result.data)
                }
                is Result.Error -> {
                    _uiState.value = ComunicadoDetailUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = ComunicadoDetailUiState(isLoading = true)
                }
            }
        }
    }
}
