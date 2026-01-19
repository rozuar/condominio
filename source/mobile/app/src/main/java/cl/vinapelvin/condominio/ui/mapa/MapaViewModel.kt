package cl.vinapelvin.condominio.ui.mapa

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.MapaArea
import cl.vinapelvin.condominio.data.model.MapaPunto
import cl.vinapelvin.condominio.data.repository.MapaRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class MapaUiState(
    val isLoading: Boolean = false,
    val areas: List<MapaArea> = emptyList(),
    val puntos: List<MapaPunto> = emptyList(),
    val error: String? = null,
    val selectedPunto: MapaPunto? = null,
    val selectedArea: MapaArea? = null,
    val isRefreshing: Boolean = false
)

@HiltViewModel
class MapaViewModel @Inject constructor(
    private val repository: MapaRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(MapaUiState(isLoading = true))
    val uiState: StateFlow<MapaUiState> = _uiState.asStateFlow()

    init {
        loadMapaData()
    }

    fun loadMapaData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            when (val result = repository.getMapaData()) {
                is Result.Success -> {
                    _uiState.value = MapaUiState(
                        areas = result.data.areas,
                        puntos = result.data.puntos
                    )
                }
                is Result.Error -> _uiState.value = MapaUiState(error = result.message)
                is Result.Loading -> _uiState.value = MapaUiState(isLoading = true)
            }
        }
    }

    fun refresh() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isRefreshing = true)
            when (val result = repository.getMapaData()) {
                is Result.Success -> {
                    _uiState.value = MapaUiState(
                        areas = result.data.areas,
                        puntos = result.data.puntos
                    )
                }
                is Result.Error -> _uiState.value = _uiState.value.copy(
                    isRefreshing = false,
                    error = result.message
                )
                is Result.Loading -> {}
            }
        }
    }

    fun selectPunto(punto: MapaPunto) {
        _uiState.value = _uiState.value.copy(selectedPunto = punto, selectedArea = null)
    }

    fun selectArea(area: MapaArea) {
        _uiState.value = _uiState.value.copy(selectedArea = area, selectedPunto = null)
    }

    fun clearSelection() {
        _uiState.value = _uiState.value.copy(selectedPunto = null, selectedArea = null)
    }
}
