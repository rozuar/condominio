package cl.vinapelvin.condominio.service

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.os.Build
import androidx.core.app.NotificationCompat
import cl.vinapelvin.condominio.MainActivity
import cl.vinapelvin.condominio.R
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class NotificationHelper @Inject constructor(
    @ApplicationContext private val context: Context
) {
    companion object {
        const val CHANNEL_GENERAL = "general"
        const val CHANNEL_EMERGENCIAS = "emergencias"
        const val CHANNEL_COMUNICADOS = "comunicados"
        const val CHANNEL_VOTACIONES = "votaciones"
    }

    private val notificationManager = context.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager

    init {
        createNotificationChannels()
    }

    private fun createNotificationChannels() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channels = listOf(
                NotificationChannel(
                    CHANNEL_GENERAL,
                    "General",
                    NotificationManager.IMPORTANCE_DEFAULT
                ).apply {
                    description = "Notificaciones generales del condominio"
                },
                NotificationChannel(
                    CHANNEL_EMERGENCIAS,
                    "Emergencias",
                    NotificationManager.IMPORTANCE_HIGH
                ).apply {
                    description = "Avisos de emergencia urgentes"
                    enableVibration(true)
                    enableLights(true)
                },
                NotificationChannel(
                    CHANNEL_COMUNICADOS,
                    "Comunicados",
                    NotificationManager.IMPORTANCE_DEFAULT
                ).apply {
                    description = "Comunicados de la comunidad"
                },
                NotificationChannel(
                    CHANNEL_VOTACIONES,
                    "Votaciones",
                    NotificationManager.IMPORTANCE_DEFAULT
                ).apply {
                    description = "Notificaciones de votaciones activas"
                }
            )

            channels.forEach { channel ->
                notificationManager.createNotificationChannel(channel)
            }
        }
    }

    fun showNotification(
        title: String,
        body: String,
        type: String = "general",
        data: Map<String, String> = emptyMap()
    ) {
        val channelId = when (type) {
            "emergencia" -> CHANNEL_EMERGENCIAS
            "comunicado" -> CHANNEL_COMUNICADOS
            "votacion" -> CHANNEL_VOTACIONES
            else -> CHANNEL_GENERAL
        }

        val intent = Intent(context, MainActivity::class.java).apply {
            flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
            data.forEach { (key, value) ->
                putExtra(key, value)
            }
        }

        val pendingIntent = PendingIntent.getActivity(
            context,
            System.currentTimeMillis().toInt(),
            intent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )

        val priority = if (type == "emergencia") {
            NotificationCompat.PRIORITY_HIGH
        } else {
            NotificationCompat.PRIORITY_DEFAULT
        }

        val notification = NotificationCompat.Builder(context, channelId)
            .setSmallIcon(R.drawable.ic_notification)
            .setContentTitle(title)
            .setContentText(body)
            .setStyle(NotificationCompat.BigTextStyle().bigText(body))
            .setPriority(priority)
            .setAutoCancel(true)
            .setContentIntent(pendingIntent)
            .build()

        notificationManager.notify(System.currentTimeMillis().toInt(), notification)
    }
}
