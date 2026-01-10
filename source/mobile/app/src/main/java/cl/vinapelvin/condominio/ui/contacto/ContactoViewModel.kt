package cl.vinapelvin.condominio.ui.contacto

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.local.TokenManager
import cl.vinapelvin.condominio.data.model.MensajeContacto
import cl.vinapelvin.condominio.data.repository.ContactoRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ContactoUiState(
    val isLoading: Boolean = false,
    val isSending: Boolean = false,
    val misMensajes: List<MensajeContacto> = emptyList(),
    val error: String? = null,
    val lastSentMessageId: String? = null
)

@HiltViewModel
class ContactoViewModel @Inject constructor(
    private val repository: ContactoRepository,
    private val tokenManager: TokenManager
) : ViewModel() {

    private val _uiState = MutableStateFlow(ContactoUiState(isLoading = true))
    val uiState: StateFlow<ContactoUiState> = _uiState.asStateFlow()

    private val _prefillNombre = MutableStateFlow("")
    val prefillNombre: StateFlow<String> = _prefillNombre.asStateFlow()

    private val _prefillEmail = MutableStateFlow("")
    val prefillEmail: StateFlow<String> = _prefillEmail.asStateFlow()

    init {
        viewModelScope.launch {
            _prefillNombre.value = tokenManager.userName.first().orEmpty()
            _prefillEmail.value = tokenManager.userEmail.first().orEmpty()
        }
        loadMisMensajes()
    }

    fun loadMisMensajes() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            when (val result = repository.getMisMensajes()) {
                is Result.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        misMensajes = result.data.mensajes,
                        error = null
                    )
                }
                is Result.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.message
                    )
                }
                is Result.Loading -> {
                    _uiState.value = _uiState.value.copy(isLoading = true)
                }
            }
        }
    }

    fun enviarMensaje(nombre: String, email: String, asunto: String, mensaje: String) {
        val n = nombre.trim()
        val e = email.trim()
        val a = asunto.trim()
        val m = mensaje.trim()
        if (n.isEmpty() || e.isEmpty() || a.isEmpty() || m.isEmpty()) {
            _uiState.value = _uiState.value.copy(error = "Todos los campos son requeridos")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isSending = true, error = null, lastSentMessageId = null)

            when (val result = repository.enviarMensajeContacto(n, e, a, m)) {
                is Result.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isSending = false,
                        lastSentMessageId = result.data.id
                    )
                    loadMisMensajes()
                }
                is Result.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isSending = false,
                        error = result.message
                    )
                }
                is Result.Loading -> {
                    _uiState.value = _uiState.value.copy(isSending = true)
                }
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

