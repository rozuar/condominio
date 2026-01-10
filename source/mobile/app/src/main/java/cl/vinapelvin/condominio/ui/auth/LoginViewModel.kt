package cl.vinapelvin.condominio.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.repository.AuthRepository
import cl.vinapelvin.condominio.data.repository.Result
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class LoginUiState(
    val isLoading: Boolean = false,
    val error: String? = null,
    val isSuccess: Boolean = false
)

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(LoginUiState())
    val uiState: StateFlow<LoginUiState> = _uiState.asStateFlow()

    val isLoggedIn: Flow<Boolean> = authRepository.isLoggedIn

    fun login(email: String, password: String) {
        if (email.isBlank() || password.isBlank()) {
            _uiState.value = LoginUiState(error = "Email y contrasena son requeridos")
            return
        }

        viewModelScope.launch {
            _uiState.value = LoginUiState(isLoading = true)

            when (val result = authRepository.login(email, password)) {
                is Result.Success -> {
                    _uiState.value = LoginUiState(isSuccess = true)
                }
                is Result.Error -> {
                    _uiState.value = LoginUiState(error = result.message)
                }
                is Result.Loading -> {
                    _uiState.value = LoginUiState(isLoading = true)
                }
            }
        }
    }

    fun logout() {
        viewModelScope.launch {
            authRepository.logout()
            _uiState.value = LoginUiState()
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}
