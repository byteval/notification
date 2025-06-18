package ports

import (
	"context"

	"notification/internal/domain/notification"
)

// Интерфейс для работы с хранилищем уведомлений
type NotificationRepository interface {
	Create(ctx context.Context, n *notification.Notification) (*notification.Notification, error)

	FindByID(ctx context.Context, id string) (*notification.Notification, error)

	UpdateStatus(ctx context.Context, id string, status notification.Status) error
}
