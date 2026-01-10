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

data class ComunicadosUiState(
    val isLoading: Boolean = false,
    val comunicados: List<Comunicado> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class ComunicadosViewModel @Inject constructor(
    private val repository: ComunicadoRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ComunicadosUiState(isLoading = true))
    val uiState: StateFlow<ComunicadosUiState> = _uiState.asStateFlow()

    init {
        loadComunicados()
    }

    fun loadComunicados() {
        viewModelScope.launch {
            _uiState.value = ComunicadosUiState(isLoading = true)

            when (val result = repository.getComunicados()) {
                is Result.Success -> {
                    _uiState.value = ComunicadosUiState(comunicados = result.data.comunicados)
                }
                is Result.Error -> {
                    _uiState.value = ComunicadosUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = ComunicadosUiState(isLoading = true)
                }
            }
        }
    }
}
