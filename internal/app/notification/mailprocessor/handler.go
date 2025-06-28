package mailprocessor

import (
	"context"
	"fmt"
	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
)

type ProcessIncomingEmail struct {
	imapClient ports.IMAPClient
}

func NewProcessIncomingEmail(imapClient ports.IMAPClient) *ProcessIncomingEmail {
	return &ProcessIncomingEmail{imapClient: imapClient}
}

func (h *ProcessIncomingEmail) Execute(ctx context.Context) error {
	h.imapClient.Connect()

	messages, err := h.imapClient.FetchUnseenEmails()
	if err != nil {
		return fmt.Errorf("failed to read INBOX: %w", err)
	}

	h.imapClient.Disconnect()

	if len(messages) > 0 {
		for _, message := range messages {
			h.findNotification(message)
		}
	}

	return nil
}

func (h *ProcessIncomingEmail) findNotification(message *notification.EmailMessage) {
	// TODO findNotification

	return
}
