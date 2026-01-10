package cl.vinapelvin.condominio

import android.app.Application
import cl.vinapelvin.condominio.service.FCMTokenManager
import cl.vinapelvin.condominio.service.NotificationHelper
import dagger.hilt.android.HiltAndroidApp
import javax.inject.Inject

@HiltAndroidApp
class CondominioApp : Application() {

    @Inject
    lateinit var notificationHelper: NotificationHelper

    @Inject
    lateinit var fcmTokenManager: FCMTokenManager

    override fun onCreate() {
        super.onCreate()
        // Initialize notification channels
        notificationHelper
        // Subscribe to default topics
        fcmTokenManager.subscribeToDefaultTopics()
    }
}
