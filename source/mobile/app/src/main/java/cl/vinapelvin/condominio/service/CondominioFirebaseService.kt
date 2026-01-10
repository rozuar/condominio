package cl.vinapelvin.condominio.service

import android.util.Log
import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class CondominioFirebaseService : FirebaseMessagingService() {

    @Inject
    lateinit var notificationHelper: NotificationHelper

    @Inject
    lateinit var fcmTokenManager: FCMTokenManager

    override fun onMessageReceived(remoteMessage: RemoteMessage) {
        super.onMessageReceived(remoteMessage)

        Log.d(TAG, "From: ${remoteMessage.from}")

        val data = remoteMessage.data
        val notification = remoteMessage.notification

        val title = notification?.title ?: data["title"] ?: "Condominio Vina Pelvin"
        val body = notification?.body ?: data["body"] ?: ""
        val type = data["type"] ?: "general"

        if (title.isNotEmpty() && body.isNotEmpty()) {
            notificationHelper.showNotification(
                title = title,
                body = body,
                type = type,
                data = data
            )
        }
    }

    override fun onNewToken(token: String) {
        super.onNewToken(token)
        Log.d(TAG, "New FCM token: $token")
        fcmTokenManager.saveToken(token)
    }

    companion object {
        private const val TAG = "FCMService"
    }
}
