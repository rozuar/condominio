package cl.vinapelvin.condominio.ui.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import cl.vinapelvin.condominio.data.api.ApiService
import cl.vinapelvin.condominio.data.local.TokenManager
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val tokenManager: TokenManager,
    private val apiService: ApiService
) : ViewModel() {

    val userName = tokenManager.userName

    private val _notificationCount = MutableStateFlow(0)
    val notificationCount: StateFlow<Int> = _notificationCount.asStateFlow()

    init {
        loadNotificationCount()
    }

    private fun loadNotificationCount() {
        viewModelScope.launch {
            try {
                val response = apiService.getNotificacionStats()
                if (response.isSuccessful) {
                    _notificationCount.value = response.body()?.unread ?: 0
                }
            } catch (e: Exception) {
                // Ignore errors for notification count
            }
        }
    }
}
