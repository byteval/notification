package get

import (
	"notification/internal/domain/notification"
)

func ToResponse(n *notification.Notification) *Response {
	if n == nil {
		return nil
	}

	return &Response{
		ID:        n.ID,
		Title:     n.Title,
		Status:    string(n.Status),
		CreatedAt: n.CreatedAt,
	}
}
