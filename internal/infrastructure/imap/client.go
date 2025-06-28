package imap

import (
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"notification/config"
	"notification/internal/domain/notification"
	"notification/pkg/logger"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type IMAPClient struct {
	config config.IMAP
	client *client.Client
	logger logger.Logger
}

func NewIMAPClient(config config.IMAP, logger logger.Logger) *IMAPClient {
	return &IMAPClient{
		config: config,
		logger: logger,
	}
}

// Connect устанавливает соединение с IMAP-сервером
func (ic *IMAPClient) Connect() error {
	address := fmt.Sprintf("%s:%d", ic.config.Host, ic.config.Port)

	var c *client.Client
	var err error

	if ic.config.UseTLS {
		tlsConfig := &tls.Config{ServerName: ic.config.Host}
		c, err = client.DialTLS(address, tlsConfig)
	} else {
		c, err = client.Dial(address)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	ic.client = c

	if ic.config.Timeout > 0 {
		ic.client.Timeout = ic.config.Timeout
	}

	if err := ic.client.Login(ic.config.Username, ic.config.Password); err != nil {
		return fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	ic.logger.Info("Successfully connected to IMAP server")
	return nil
}

// Disconnect завершает соединение с IMAP-сервером
func (ic *IMAPClient) Disconnect() {
	if ic.client != nil {
		if err := ic.client.Logout(); err != nil {
			ic.logger.Error("Failed to logout from IMAP server", "error", err)
		}
		ic.logger.Info("Disconnected from IMAP server")
	}
}

func (ic *IMAPClient) FetchUnseenEmails() ([]*notification.EmailMessage, error) {
	if ic.client == nil {
		return nil, fmt.Errorf("IMAP client is not connected")
	}

	// Выбираем INBOX в режиме не только для чтения
	mbox, err := ic.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}

	if mbox.Messages == 0 {
		return []*notification.EmailMessage{}, nil
	}

	// Ищем непрочитанные сообщения
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	ids, err := ic.client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	if len(ids) == 0 {
		return []*notification.EmailMessage{}, nil
	}

	ic.logger.Info("New messages", "count", len(ids))

	// Ограничиваем количество сообщений
	if len(ids) > 10 {
		ids = ids[len(ids)-10:]
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	// Запрашиваем полное сообщение и envelope
	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchRFC822,
		imap.FetchInternalDate,
	}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- ic.client.Fetch(seqset, items, messages)
	}()

	var result []*notification.EmailMessage
	for msg := range messages {
		email, err := ic.parseMessage(msg)
		if err != nil {
			ic.logger.Error("Failed to parse message",
				"error", err,
				"message_id", msg.SeqNum)
			continue
		}
		result = append(result, email)

		// Помечаем сообщение как прочитанное
		if err := ic.markAsRead(msg.SeqNum); err != nil {
			ic.logger.Error("Failed to mark message as read",
				"error", err,
				"message_id", msg.SeqNum)
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	return result, nil
}

// parseMessage преобразует imap.Message в notification.EmailMessage
func (ic *IMAPClient) parseMessage(msg *imap.Message) (*notification.EmailMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}

	email := &notification.EmailMessage{
		UID:     msg.Uid,
		Subject: msg.Envelope.Subject,
		From:    ic.formatAddress(msg.Envelope.From),
		To:      ic.formatAddresses(msg.Envelope.To),
		Date:    msg.InternalDate.Format(time.RFC3339),
	}

	// Получаем тело сообщения
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		return nil, fmt.Errorf("failed to get message body")
	}

	// Создаем reader для MIME-парсинга прямо из потока данных
	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail reader: %w", err)
	}

	// Обрабатываем части сообщения
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read part: %w", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			body, err := io.ReadAll(p.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read body: %w", err)
			}

			switch contentType {
			case "text/plain":
				email.TextBody = string(body)
			case "text/html":
				email.HTMLBody = string(body)
			}

		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			contentType, _, _ := h.ContentType()
			body, err := io.ReadAll(p.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read attachment: %w", err)
			}

			email.Attachments = append(email.Attachments, notification.EmailAttachment{
				Name:        filename,
				ContentType: contentType,
				Content:     body,
			})
		}
	}

	return email, nil
}

// markAsRead помечает сообщение как прочитанное
func (ic *IMAPClient) markAsRead(seqNum uint32) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	return ic.client.Store(seqset, item, flags, nil)
}

// formatAddress форматирует адрес email
func (ic *IMAPClient) formatAddress(addrs []*imap.Address) string {
	if len(addrs) == 0 {
		return ""
	}
	return fmt.Sprintf("%s@%s", addrs[0].MailboxName, addrs[0].HostName)
}

// formatAddresses форматирует список адресов email
func (ic *IMAPClient) formatAddresses(addrs []*imap.Address) []string {
	var result []string
	for _, addr := range addrs {
		result = append(result, ic.formatAddress([]*imap.Address{addr}))
	}
	return result
}
