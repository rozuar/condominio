package cl.vinapelvin.condominio.service

import android.content.Context
import android.util.Log
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import com.google.firebase.messaging.FirebaseMessaging
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.launch
import kotlinx.coroutines.tasks.await
import javax.inject.Inject
import javax.inject.Singleton

private val Context.fcmDataStore by preferencesDataStore(name = "fcm_prefs")

@Singleton
class FCMTokenManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val scope = CoroutineScope(Dispatchers.IO)

    companion object {
        private const val TAG = "FCMTokenManager"
        private val FCM_TOKEN_KEY = stringPreferencesKey("fcm_token")
    }

    fun saveToken(token: String) {
        scope.launch {
            context.fcmDataStore.edit { preferences ->
                preferences[FCM_TOKEN_KEY] = token
            }
            Log.d(TAG, "FCM token saved")
        }
    }

    suspend fun getToken(): String? {
        return try {
            val storedToken = context.fcmDataStore.data.map { preferences ->
                preferences[FCM_TOKEN_KEY]
            }.first()

            if (storedToken != null) {
                storedToken
            } else {
                val newToken = FirebaseMessaging.getInstance().token.await()
                saveToken(newToken)
                newToken
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error getting FCM token", e)
            null
        }
    }

    fun subscribeToTopic(topic: String) {
        FirebaseMessaging.getInstance().subscribeToTopic(topic)
            .addOnCompleteListener { task ->
                if (task.isSuccessful) {
                    Log.d(TAG, "Subscribed to topic: $topic")
                } else {
                    Log.e(TAG, "Failed to subscribe to topic: $topic")
                }
            }
    }

    fun unsubscribeFromTopic(topic: String) {
        FirebaseMessaging.getInstance().unsubscribeFromTopic(topic)
            .addOnCompleteListener { task ->
                if (task.isSuccessful) {
                    Log.d(TAG, "Unsubscribed from topic: $topic")
                } else {
                    Log.e(TAG, "Failed to unsubscribe from topic: $topic")
                }
            }
    }

    fun subscribeToDefaultTopics() {
        subscribeToTopic("all_users")
        subscribeToTopic("comunicados")
        subscribeToTopic("emergencias")
    }
}
