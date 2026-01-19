package cl.vinapelvin.condominio.ui.galerias

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Galeria
import cl.vinapelvin.condominio.data.repository.GaleriaRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class GaleriasUiState(
    val isLoading: Boolean = false,
    val galerias: List<Galeria> = emptyList(),
    val error: String? = null,
    val isRefreshing: Boolean = false
)

@HiltViewModel
class GaleriasViewModel @Inject constructor(
    private val repository: GaleriaRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(GaleriasUiState(isLoading = true))
    val uiState: StateFlow<GaleriasUiState> = _uiState.asStateFlow()

    init {
        loadGalerias()
    }

    fun loadGalerias() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            when (val result = repository.getGalerias()) {
                is Result.Success -> _uiState.value = GaleriasUiState(galerias = result.data.galerias)
                is Result.Error -> _uiState.value = GaleriasUiState(error = result.message)
                is Result.Loading -> _uiState.value = GaleriasUiState(isLoading = true)
            }
        }
    }

    fun refresh() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isRefreshing = true)
            when (val result = repository.getGalerias()) {
                is Result.Success -> _uiState.value = GaleriasUiState(galerias = result.data.galerias)
                is Result.Error -> _uiState.value = _uiState.value.copy(isRefreshing = false, error = result.message)
                is Result.Loading -> {}
            }
        }
    }
}
