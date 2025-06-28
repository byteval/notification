package ports

// Интерфейс для отправки уведомлений по SMTP
type SMTPSender interface {
	Send(id string, to string, subject string, message string, attachments map[string]string) error
}
