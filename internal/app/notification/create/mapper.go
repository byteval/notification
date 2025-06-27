package create

import (
	"notification/internal/domain/notification"
	"time"
)

// ToDomain конвертирует Request в доменную модель Notification и список ReceiverEmail
func ToDomain(r Request) (*notification.Notification, []notification.NotificationReceiver, error) {
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

	return n, receivers, nil
}

// ToResponse конвертирует доменную модель в Response
func ToResponse(n *notification.Notification) *Response {
	response := &Response{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
	}

	// Добавляем информацию о получателях
	for _, receiver := range n.NotificationReceivers {
		response.Receivers = append(response.Receivers, ReceiverResponse{
			Email:  receiver.Email,
			Status: string(receiver.Status),
		})
	}

	return response
}
