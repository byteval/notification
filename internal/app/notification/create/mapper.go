package create

import (
	"notification/internal/domain/notification"
	"time"
)

// Request в доменную модель
func (r Request) ToDomain() (*notification.Notification, error) {
	return &notification.Notification{
		LayoutID:  r.LayoutID,
		Status:    notification.StatusPending,
		Title:     r.Title,
		Content:   r.Content,
		Data:      notification.MapStringInterface(r.Data),
		Channels:  notification.StringArray(r.Channels),
		Receiver:  r.Receiver,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// Response из доменной модели
func ToResponse(n *notification.Notification) *Response {
	return &Response{
		ID:        n.ID,
		Status:    string(n.Status),
		CreatedAt: n.CreatedAt,
	}
}
