package smtp

import (
	"net/smtp"
	"notification/config"
	"notification/pkg/logger"
	"strconv"
)

type SMTPSender struct {
	cfg    config.Config
	logger logger.Logger
}

func NewSender(cfg config.Config, logger logger.Logger) *SMTPSender {
	return &SMTPSender{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *SMTPSender) Send(to string, subject string, message string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTP.Username, s.cfg.SMTP.Password, s.cfg.SMTP.Host)

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message

	err := smtp.SendMail(
		s.cfg.SMTP.Host+":"+strconv.Itoa(s.cfg.SMTP.Port),
		auth,
		s.cfg.SMTP.From,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		s.logger.Error("Failed to send email", "error", err)
		return err
	}

	return nil
}
