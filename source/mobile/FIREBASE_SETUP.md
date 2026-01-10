# Configuracion de Firebase para Condominio App

## Pasos para configurar Firebase Cloud Messaging

### 1. Crear proyecto en Firebase Console

1. Ve a [Firebase Console](https://console.firebase.google.com/)
2. Haz clic en "Agregar proyecto"
3. Nombre del proyecto: `Condominio Vina Pelvin`
4. Sigue los pasos para crear el proyecto

### 2. Agregar aplicacion Android

1. En la consola de Firebase, haz clic en el icono de Android
2. Completa los datos:
   - **Package name**: `cl.vinapelvin.condominio`
   - **App nickname**: `Condominio App`
   - **SHA-1**: (opcional para FCM basico)

### 3. Descargar google-services.json

1. Descarga el archivo `google-services.json`
2. Colocalo en: `app/google-services.json`

### 4. Verificar la configuracion

El archivo debe tener esta estructura basica:

```json
{
  "project_info": {
    "project_number": "...",
    "project_id": "condominio-vina-pelvin",
    "storage_bucket": "..."
  },
  "client": [
    {
      "client_info": {
        "mobilesdk_app_id": "...",
        "android_client_info": {
          "package_name": "cl.vinapelvin.condominio"
        }
      },
      ...
    }
  ]
}
```

## Enviar notificaciones desde el servidor

### Usando Firebase Admin SDK (Backend)

```python
import firebase_admin
from firebase_admin import credentials, messaging

# Inicializar Firebase Admin
cred = credentials.Certificate("path/to/serviceAccountKey.json")
firebase_admin.initialize_app(cred)

# Enviar a un topic
message = messaging.Message(
    notification=messaging.Notification(
        title='Nueva Emergencia',
        body='Se ha reportado una emergencia en el condominio'
    ),
    data={
        'type': 'emergencia',
        'id': '123'
    },
    topic='emergencias'
)

response = messaging.send(message)
```

### Topics disponibles

- `all_users` - Todos los usuarios
- `comunicados` - Nuevos comunicados
- `emergencias` - Avisos de emergencia

## Probar notificaciones

1. Compila y ejecuta la app
2. Ve a Firebase Console > Cloud Messaging
3. Haz clic en "Send your first message"
4. Ingresa titulo y cuerpo
5. Selecciona la app Android
6. Envia la notificacion

## Notas importantes

- El archivo `google-services.json` contiene claves sensibles
- **NO** subir este archivo al repositorio publico
- Agregar a `.gitignore`: `app/google-services.json`
