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

// ToResponse конвертирует доменную модель в Response
func ToResponse(n *notification.Notification) *Response {
	response := &Response{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
	}

	for _, receiver := range n.NotificationReceivers {
		response.Receivers = append(response.Receivers, ReceiverResponse{
			Email:  receiver.Email,
			Status: string(receiver.Status),
		})
	}

	for _, attachment := range n.Attachments {
		response.Attachments = append(response.Attachments, AttachmentResponse{
			FileName: attachment.FileName,
			Size:     attachment.Size,
		})
	}

	return response
}
