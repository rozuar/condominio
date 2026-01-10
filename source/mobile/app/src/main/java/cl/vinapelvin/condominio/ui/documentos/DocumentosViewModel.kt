package cl.vinapelvin.condominio.ui.documentos

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.model.Documento
import cl.vinapelvin.condominio.data.repository.DocumentosRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class DocumentosUiState(
    val isLoading: Boolean = false,
    val documentos: List<Documento> = emptyList(),
    val error: String? = null
)

@HiltViewModel
class DocumentosViewModel @Inject constructor(
    private val repository: DocumentosRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(DocumentosUiState(isLoading = true))
    val uiState: StateFlow<DocumentosUiState> = _uiState.asStateFlow()

    init {
        loadDocumentos()
    }

    fun loadDocumentos() {
        viewModelScope.launch {
            _uiState.value = DocumentosUiState(isLoading = true)
            when (val result = repository.getDocumentos()) {
                is Result.Success -> _uiState.value = DocumentosUiState(documentos = result.data.documentos)
                is Result.Error -> _uiState.value = DocumentosUiState(error = result.message)
                is Result.Loading -> _uiState.value = DocumentosUiState(isLoading = true)
            }
        }
    }
}

