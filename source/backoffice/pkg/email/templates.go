package email

// Template names
const (
	TemplateContactoRecibido    = "contacto_recibido"
	TemplateContactoRespuesta   = "contacto_respuesta"
	TemplateNotificacion        = "notificacion"
	TemplateEmergencia          = "emergencia"
	TemplateVotacionNueva       = "votacion_nueva"
	TemplateGastoComun          = "gasto_comun"
	TemplatePagoRecibido        = "pago_recibido"
	TemplateBienvenida          = "bienvenida"
	TemplateRecuperarPassword   = "recuperar_password"
)

// BaseTemplate is the base HTML template for all emails
const BaseTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            padding: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #2D5016;
        }
        .header h1 {
            color: #2D5016;
            margin: 0;
            font-size: 24px;
        }
        .header p {
            color: #666;
            margin: 5px 0 0 0;
            font-size: 14px;
        }
        .content {
            margin-bottom: 30px;
        }
        .content h2 {
            color: #2D5016;
            font-size: 20px;
            margin-top: 0;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #2D5016;
            color: #ffffff !important;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 600;
            margin: 10px 0;
        }
        .button:hover {
            background-color: #4A7C23;
        }
        .info-box {
            background-color: #f0f7eb;
            border-left: 4px solid #2D5016;
            padding: 15px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .warning-box {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .danger-box {
            background-color: #f8d7da;
            border-left: 4px solid #dc3545;
            padding: 15px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .footer {
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
            font-size: 12px;
        }
        .footer a {
            color: #2D5016;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Comunidad Viña Pelvin</h1>
            <p>Parcelas de Agrado</p>
        </div>
        <div class="content">
            {{.Content}}
        </div>
        <div class="footer">
            <p>Este es un correo automatico de Comunidad Viña Pelvin.</p>
            <p>Por favor no responda directamente a este correo.</p>
            <p><a href="https://vinapelvin.cl">vinapelvin.cl</a></p>
        </div>
    </div>
</body>
</html>
`

// ContactoRecibidoTemplate is sent to the user when they submit a contact form
const ContactoRecibidoTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mensaje Recibido</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #2D5016; }
        .header h1 { color: #2D5016; margin: 0; font-size: 24px; }
        .info-box { background-color: #f0f7eb; border-left: 4px solid #2D5016; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <h2>Hemos recibido tu mensaje</h2>
            <p>Hola <strong>{{.Nombre}}</strong>,</p>
            <p>Gracias por contactarnos. Hemos recibido tu mensaje y la directiva lo revisara a la brevedad.</p>
            <div class="info-box">
                <p><strong>Asunto:</strong> {{.Asunto}}</p>
                <p><strong>Mensaje:</strong></p>
                <p>{{.Mensaje}}</p>
            </div>
            <p>Te responderemos a este correo electronico lo antes posible.</p>
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// ContactoRespuestaTemplate is sent when the directiva responds to a contact message
const ContactoRespuestaTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Respuesta a tu mensaje</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #2D5016; }
        .header h1 { color: #2D5016; margin: 0; font-size: 24px; }
        .info-box { background-color: #f0f7eb; border-left: 4px solid #2D5016; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .response-box { background-color: #e8f4fd; border-left: 4px solid #0d6efd; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <h2>Respuesta de la Directiva</h2>
            <p>Hola <strong>{{.Nombre}}</strong>,</p>
            <p>La directiva ha respondido a tu mensaje:</p>
            <div class="info-box">
                <p><strong>Tu mensaje original:</strong></p>
                <p><em>{{.Asunto}}</em></p>
                <p>{{.MensajeOriginal}}</p>
            </div>
            <div class="response-box">
                <p><strong>Respuesta:</strong></p>
                <p>{{.Respuesta}}</p>
            </div>
            <p>Si tienes mas consultas, no dudes en contactarnos nuevamente.</p>
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// EmergenciaTemplate is sent for emergency notifications
const EmergenciaTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Priority}} - {{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #dc3545; }
        .header h1 { color: #dc3545; margin: 0; font-size: 24px; }
        .priority-critical { background-color: #f8d7da; border-left: 4px solid #dc3545; }
        .priority-high { background-color: #fff3cd; border-left: 4px solid #ffc107; }
        .priority-medium { background-color: #cfe2ff; border-left: 4px solid #0d6efd; }
        .priority-low { background-color: #d1e7dd; border-left: 4px solid #198754; }
        .alert-box { padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .button { display: inline-block; padding: 12px 24px; background-color: #dc3545; color: #ffffff !important; text-decoration: none; border-radius: 6px; font-weight: 600; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚠️ ALERTA - Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <div class="alert-box priority-{{.PriorityClass}}">
                <h2 style="margin-top: 0;">{{.Title}}</h2>
                <p><strong>Prioridad:</strong> {{.Priority}}</p>
            </div>
            <p>{{.Content}}</p>
            {{if .URL}}
            <p style="text-align: center;">
                <a href="{{.URL}}" class="button">Ver detalles</a>
            </p>
            {{end}}
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// NotificacionTemplate is a generic notification template
const NotificacionTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #2D5016; }
        .header h1 { color: #2D5016; margin: 0; font-size: 24px; }
        .info-box { background-color: #f0f7eb; border-left: 4px solid #2D5016; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .button { display: inline-block; padding: 12px 24px; background-color: #2D5016; color: #ffffff !important; text-decoration: none; border-radius: 6px; font-weight: 600; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <h2>{{.Title}}</h2>
            <div class="info-box">
                <p>{{.Body}}</p>
            </div>
            {{if .URL}}
            <p style="text-align: center;">
                <a href="{{.URL}}" class="button">Ver mas</a>
            </p>
            {{end}}
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// GastoComunTemplate is for expense notifications
const GastoComunTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gasto Comun - {{.Periodo}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #2D5016; }
        .header h1 { color: #2D5016; margin: 0; font-size: 24px; }
        .amount { font-size: 32px; font-weight: bold; color: #2D5016; text-align: center; margin: 20px 0; }
        .info-box { background-color: #f0f7eb; border-left: 4px solid #2D5016; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .warning-box { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .button { display: inline-block; padding: 12px 24px; background-color: #2D5016; color: #ffffff !important; text-decoration: none; border-radius: 6px; font-weight: 600; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <h2>Gasto Comun - {{.Periodo}}</h2>
            <p>Hola <strong>{{.Nombre}}</strong>,</p>
            <p>Se ha generado el gasto comun correspondiente al periodo {{.Periodo}}.</p>
            <div class="amount">${{.Monto}}</div>
            <div class="info-box">
                <p><strong>Parcela:</strong> {{.Parcela}}</p>
                <p><strong>Fecha de vencimiento:</strong> {{.FechaVencimiento}}</p>
            </div>
            {{if .Descripcion}}
            <p>{{.Descripcion}}</p>
            {{end}}
            <div class="warning-box">
                <p>Recuerde realizar el pago antes de la fecha de vencimiento para evitar recargos.</p>
            </div>
            {{if .URL}}
            <p style="text-align: center;">
                <a href="{{.URL}}" class="button">Ver estado de cuenta</a>
            </p>
            {{end}}
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// BienvenidaTemplate is sent to new users
const BienvenidaTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bienvenido a Comunidad Viña Pelvin</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding-bottom: 20px; border-bottom: 2px solid #2D5016; }
        .header h1 { color: #2D5016; margin: 0; font-size: 24px; }
        .info-box { background-color: #f0f7eb; border-left: 4px solid #2D5016; padding: 15px; margin: 20px 0; border-radius: 0 4px 4px 0; }
        .button { display: inline-block; padding: 12px 24px; background-color: #2D5016; color: #ffffff !important; text-decoration: none; border-radius: 6px; font-weight: 600; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏡 Comunidad Viña Pelvin</h1>
        </div>
        <div class="content">
            <h2>¡Bienvenido/a {{.Nombre}}!</h2>
            <p>Tu cuenta ha sido creada exitosamente en el portal de la Comunidad Viña Pelvin.</p>
            <div class="info-box">
                <p><strong>Email:</strong> {{.Email}}</p>
                {{if .Parcela}}
                <p><strong>Parcela:</strong> {{.Parcela}}</p>
                {{end}}
            </div>
            <p>Desde el portal podras:</p>
            <ul>
                <li>Ver comunicados y eventos de la comunidad</li>
                <li>Consultar tu estado de cuenta de gastos comunes</li>
                <li>Participar en votaciones</li>
                <li>Acceder a documentos y actas</li>
                <li>Contactar a la directiva</li>
            </ul>
            <p style="text-align: center;">
                <a href="{{.URL}}" class="button">Ingresar al portal</a>
            </p>
        </div>
        <div class="footer">
            <p>Comunidad Viña Pelvin - Parcelas de Agrado</p>
        </div>
    </div>
</body>
</html>
`

// RegisterAllTemplates registers all email templates with the service
func RegisterAllTemplates(s *Service) error {
	templates := map[string]string{
		TemplateContactoRecibido:  ContactoRecibidoTemplate,
		TemplateContactoRespuesta: ContactoRespuestaTemplate,
		TemplateNotificacion:      NotificacionTemplate,
		TemplateEmergencia:        EmergenciaTemplate,
		TemplateGastoComun:        GastoComunTemplate,
		TemplateBienvenida:        BienvenidaTemplate,
	}

	for name, content := range templates {
		if err := s.RegisterTemplateString(name, content); err != nil {
			return err
		}
	}

	return nil
}
