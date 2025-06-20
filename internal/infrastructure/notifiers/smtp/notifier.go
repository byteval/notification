package smtp

import (
	"crypto/tls"
	"notification/config"
	"notification/pkg/logger"

	"gopkg.in/gomail.v2"
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

// Отправка сообщения по SMTP
func (s *SMTPSender) Send(to string, subject string, message string) error {
	d := gomail.NewDialer(
		s.cfg.SMTP.Host,
		s.cfg.SMTP.Port,
		s.cfg.SMTP.Username,
		s.cfg.SMTP.Password,
	)

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
		/* VerifyConnection: func(cs tls.ConnectionState) error {
			opts := x509.VerifyOptions{
				DNSName:       cs.ServerName,
				Intermediates: x509.NewCertPool(),
			}
			for _, cert := range cs.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := cs.PeerCertificates[0].Verify(opts)
			return err
		}, */
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTP.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	m.SetHeader("Non-Delivery-Report", s.cfg.SMTP.Username)
	// m.SetHeader("Disposition-Notification-To", s.cfg.SMTP.Username)
	//m.SetHeader("Message-Id", strconv.FormatInt(ID, 10))

	// Добавление вложения:
	// m.Attach("/path/to/file.pdf")

	// Встроенные изображения
	// m.Embed("/tmp/image.png")

	if err := d.DialAndSend(m); err != nil {
		s.logger.Error("Failed to send email", "error", err)
		return err
	}

	return nil
}
