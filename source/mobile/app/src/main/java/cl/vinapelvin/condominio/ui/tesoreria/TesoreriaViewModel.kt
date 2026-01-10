package cl.vinapelvin.condominio.ui.tesoreria

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Movimiento
import cl.vinapelvin.condominio.data.model.TesoreriaResumen
import cl.vinapelvin.condominio.data.repository.Result
import cl.vinapelvin.condominio.data.repository.TesoreriaRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class TesoreriaUiState(
    val isLoading: Boolean = false,
    val resumen: TesoreriaResumen? = null,
    val movimientos: List<Movimiento> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class TesoreriaViewModel @Inject constructor(
    private val repository: TesoreriaRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(TesoreriaUiState(isLoading = true))
    val uiState: StateFlow<TesoreriaUiState> = _uiState.asStateFlow()

    init {
        load()
    }

    fun load() {
        viewModelScope.launch {
            _uiState.value = TesoreriaUiState(isLoading = true)

            val resumenResult = repository.getResumen()
            val movimientosResult = repository.getMovimientos()

            val resumen = (resumenResult as? Result.Success)?.data
            val movimientos = (movimientosResult as? Result.Success)?.data?.movimientos.orEmpty()

            val error = when {
                resumenResult is Result.Error -> resumenResult.message
                movimientosResult is Result.Error -> movimientosResult.message
                else -> null
            }

            _uiState.value = TesoreriaUiState(
                isLoading = false,
                resumen = resumen,
                movimientos = movimientos,
                error = error
            )
        }
    }
}

