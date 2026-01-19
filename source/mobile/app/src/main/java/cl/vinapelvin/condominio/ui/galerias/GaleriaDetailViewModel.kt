package cl.vinapelvin.condominio.ui.galerias

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.GaleriaWithItems
import cl.vinapelvin.condominio.data.repository.GaleriaRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class GaleriaDetailUiState(
    val isLoading: Boolean = false,
    val galeria: GaleriaWithItems? = null,
    val error: String? = null,
    val selectedImageIndex: Int = -1
)

@HiltViewModel
class GaleriaDetailViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val repository: GaleriaRepository
) : ViewModel() {

    private val galeriaId: String = savedStateHandle.get<String>("id") ?: ""

    private val _uiState = MutableStateFlow(GaleriaDetailUiState(isLoading = true))
    val uiState: StateFlow<GaleriaDetailUiState> = _uiState.asStateFlow()

    init {
        loadGaleria()
    }

    fun loadGaleria() {
        viewModelScope.launch {
            _uiState.value = GaleriaDetailUiState(isLoading = true)
            when (val result = repository.getGaleria(galeriaId)) {
                is Result.Success -> _uiState.value = GaleriaDetailUiState(galeria = result.data)
                is Result.Error -> _uiState.value = GaleriaDetailUiState(error = result.message)
                is Result.Loading -> _uiState.value = GaleriaDetailUiState(isLoading = true)
            }
        }
    }

    fun selectImage(index: Int) {
        _uiState.value = _uiState.value.copy(selectedImageIndex = index)
    }

    fun clearSelection() {
        _uiState.value = _uiState.value.copy(selectedImageIndex = -1)
    }
}
