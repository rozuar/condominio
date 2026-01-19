package cl.vinapelvin.condominio.ui.theme

import android.app.Activity
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat

// Muted accent palette (avoid overly saturated "kids app" look)
val Blue400 = Color(0xFF3B82F6)
val Blue500 = Color(0xFF2563EB)
val Blue600 = Color(0xFF1F4B99)
val Blue700 = Color(0xFF1B3F80)
val Blue800 = Color(0xFF163462)
val Green500 = Color(0xFF22C55E)
val Green600 = Color(0xFF1E6B3A)
val Red500 = Color(0xFFEF4444)
val Red600 = Color(0xFFB3261E)
val Amber500 = Color(0xFFC58F00)
// Brand palette (docs)
val VerdePrincipal = Color(0xFF2D5016)
val VerdeClaro = Color(0xFF4A7C23)
val Tierra = Color(0xFF8B7355)
val Agua = Color(0xFF3B82A0)
val Gray50 = Color(0xFFF9FAFB)
val Gray100 = Color(0xFFF3F4F6)
val Gray200 = Color(0xFFE5E7EB)
val Gray400 = Color(0xFF9CA3AF)
val Gray500 = Color(0xFF6B7280)
val Gray600 = Color(0xFF4B5563)
val Gray700 = Color(0xFF374151)
val Gray900 = Color(0xFF111827)

private val LightColorScheme = lightColorScheme(
    primary = VerdePrincipal,
    onPrimary = Color.White,
    primaryContainer = VerdeClaro,
    secondary = Gray500,
    onSecondary = Color.White,
    background = Gray50,
    onBackground = Gray900,
    surface = Color.White,
    onSurface = Gray900,
    error = Red600,
    onError = Color.White,
)

private val DarkColorScheme = darkColorScheme(
    primary = VerdePrincipal,
    onPrimary = Color.White,
    primaryContainer = VerdeClaro,
    secondary = Gray500,
    onSecondary = Color.White,
    background = Gray900,
    onBackground = Color.White,
    surface = Color(0xFF1F2937),
    onSurface = Color.White,
    error = Red600,
    onError = Color.White,
)

@Composable
fun CondominioTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) DarkColorScheme else LightColorScheme

    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.primaryContainer.toArgb()
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = false
        }
    }

    MaterialTheme(
        colorScheme = colorScheme,
        content = content
    )
}
