package cl.vinapelvin.condominio.ui.gastos

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.MiEstadoCuenta
import cl.vinapelvin.condominio.data.repository.GastosRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class GastosUiState(
    val isLoading: Boolean = false,
    val estadoCuenta: MiEstadoCuenta? = null,
    val error: String? = null
)

@HiltViewModel
class GastosViewModel @Inject constructor(
    private val repository: GastosRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(GastosUiState())
    val uiState: StateFlow<GastosUiState> = _uiState.asStateFlow()

    init {
        loadEstadoCuenta()
    }

    fun loadEstadoCuenta() {
        viewModelScope.launch {
            _uiState.value = GastosUiState(isLoading = true)

            when (val result = repository.getMiEstadoCuenta()) {
                is Result.Success -> {
                    _uiState.value = GastosUiState(estadoCuenta = result.data)
                }
                is Result.Error -> {
                    _uiState.value = GastosUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = GastosUiState(isLoading = true)
                }
            }
        }
    }
}
