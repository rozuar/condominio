package cl.vinapelvin.condominio.data.local

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.intPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "auth_prefs")

@Singleton
class TokenManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    companion object {
        private val ACCESS_TOKEN_KEY = stringPreferencesKey("access_token")
        private val REFRESH_TOKEN_KEY = stringPreferencesKey("refresh_token")
        private val USER_ID_KEY = stringPreferencesKey("user_id")
        private val USER_EMAIL_KEY = stringPreferencesKey("user_email")
        private val USER_NAME_KEY = stringPreferencesKey("user_name")
        private val USER_ROLE_KEY = stringPreferencesKey("user_role")
        private val USER_PARCELA_ID_KEY = intPreferencesKey("user_parcela_id")
    }

    val accessToken: Flow<String?> = context.dataStore.data.map { it[ACCESS_TOKEN_KEY] }
    val refreshToken: Flow<String?> = context.dataStore.data.map { it[REFRESH_TOKEN_KEY] }
    val isLoggedIn: Flow<Boolean> = context.dataStore.data.map { it[ACCESS_TOKEN_KEY] != null }

    val userEmail: Flow<String?> = context.dataStore.data.map { it[USER_EMAIL_KEY] }
    val userName: Flow<String?> = context.dataStore.data.map { it[USER_NAME_KEY] }
    val userRole: Flow<String?> = context.dataStore.data.map { it[USER_ROLE_KEY] }
    val userParcelaId: Flow<Int?> = context.dataStore.data.map { it[USER_PARCELA_ID_KEY] }

    suspend fun getAccessToken(): String? {
        return context.dataStore.data.first()[ACCESS_TOKEN_KEY]
    }

    suspend fun getRefreshToken(): String? {
        return context.dataStore.data.first()[REFRESH_TOKEN_KEY]
    }

    suspend fun saveTokens(accessToken: String, refreshToken: String) {
        context.dataStore.edit { prefs ->
            prefs[ACCESS_TOKEN_KEY] = accessToken
            prefs[REFRESH_TOKEN_KEY] = refreshToken
        }
    }

    suspend fun saveUser(id: String, email: String, name: String, role: String, parcelaId: Int?) {
        context.dataStore.edit { prefs ->
            prefs[USER_ID_KEY] = id
            prefs[USER_EMAIL_KEY] = email
            prefs[USER_NAME_KEY] = name
            prefs[USER_ROLE_KEY] = role
            if (parcelaId != null) {
                prefs[USER_PARCELA_ID_KEY] = parcelaId
            } else {
                prefs.remove(USER_PARCELA_ID_KEY)
            }
        }
    }

    suspend fun clearAll() {
        context.dataStore.edit { it.clear() }
    }
}
