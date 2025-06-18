package ports

// Интерфейс для отправки уведомлений по SMTP
type SMTPSender interface {
	Send(to string, subject string, message string) error
}
