package ports

import (
	"notification/internal/domain/notification"
)

type IMAPClient interface {
	Connect() error
	Disconnect()
	FetchUnseenEmails() ([]*notification.EmailMessage, error)
}
