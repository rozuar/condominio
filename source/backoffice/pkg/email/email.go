package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
)

// Config holds SMTP configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string
	Enabled  bool
}

// Service handles email sending
type Service struct {
	config    Config
	templates map[string]*template.Template
}

// NewService creates a new email service
func NewService(cfg Config) *Service {
	return &Service{
		config:    cfg,
		templates: make(map[string]*template.Template),
	}
}

// Email represents an email message
type Email struct {
	To       []string
	Subject  string
	Body     string
	HTML     bool
	Template string
	Data     interface{}
}

// Send sends an email
func (s *Service) Send(email Email) error {
	if !s.config.Enabled {
		log.Printf("[EMAIL] Disabled - Would send to %v: %s", email.To, email.Subject)
		return nil
	}

	if s.config.Host == "" {
		log.Printf("[EMAIL] No SMTP host configured - Would send to %v: %s", email.To, email.Subject)
		return nil
	}

	// Build message
	var body string
	if email.Template != "" && email.Data != nil {
		tmpl, ok := s.templates[email.Template]
		if !ok {
			return fmt.Errorf("template %s not found", email.Template)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, email.Data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
		body = buf.String()
		email.HTML = true
	} else {
		body = email.Body
	}

	// Build headers
	from := s.config.From
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(email.To, ", ")
	headers["Subject"] = email.Subject
	headers["MIME-Version"] = "1.0"

	if email.HTML {
		headers["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var auth smtp.Auth
	if s.config.User != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)
	}

	// Try TLS first (port 587), fall back to regular SMTP
	if s.config.Port == 587 || s.config.Port == 465 {
		return s.sendTLS(addr, auth, email.To, msg.Bytes())
	}

	return smtp.SendMail(addr, auth, s.config.From, email.To, msg.Bytes())
}

func (s *Service) sendTLS(addr string, auth smtp.Auth, to []string, msg []byte) error {
	// Connect
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.config.Host,
	})
	if err != nil {
		// Try STARTTLS instead
		return s.sendSTARTTLS(addr, auth, to, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(s.config.From); err != nil {
		return err
	}

	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	return w.Close()
}

func (s *Service) sendSTARTTLS(addr string, auth smtp.Auth, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	// STARTTLS
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		// Continue without TLS if STARTTLS fails
		log.Printf("[EMAIL] STARTTLS failed, continuing without TLS: %v", err)
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(s.config.From); err != nil {
		return err
	}

	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	return w.Close()
}

// RegisterTemplate registers an HTML template for emails
func (s *Service) RegisterTemplate(name string, tmpl *template.Template) {
	s.templates[name] = tmpl
}

// RegisterTemplateString registers a template from a string
func (s *Service) RegisterTemplateString(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}
	s.templates[name] = tmpl
	return nil
}

// IsEnabled returns whether email sending is enabled
func (s *Service) IsEnabled() bool {
	return s.config.Enabled && s.config.Host != ""
}
