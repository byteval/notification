package create

import (
	"notification/internal/domain/notification"
	"time"
)

// ToDomain конвертирует Request в доменную модель Notification и список ReceiverEmail
func ToDomain(r Request) (*notification.Notification, error) {
	// Создаем основное уведомление
	n := &notification.Notification{
		LayoutID:  r.LayoutID,
		Title:     r.Title,
		Data:      notification.JSONB(r.Data),
		CreatedAt: time.Now().UTC(),
	}

	// Создаем получателей
	var receivers []notification.NotificationReceiver
	for _, email := range r.Emails {
		receivers = append(receivers, notification.NotificationReceiver{
			Email:  email,
			Status: notification.StatusPending,
		})
	}

	n.NotificationReceivers = receivers

	// Создаем вложения
	var attachments []notification.Attachment
	for _, attachment := range r.Attachments {
		attachments = append(attachments, notification.Attachment{
			FileName:     attachment.Filename,
			OriginalName: attachment.Filename,
			ContentType:  attachment.ContentType,
			Size:         attachment.Size,
			FilePath:     attachment.FilePath,
		})
	}

	n.Attachments = attachments

	return n, nil
}
