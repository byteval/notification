package get

import (
	"context"
	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
)

type Getter struct {
	repo ports.NotificationRepository
}

func NewGetter(repo ports.NotificationRepository) *Getter {
	return &Getter{repo: repo}
}

// GetByID возвращает уведомление по ID
func (g *Getter) GetByID(ctx context.Context, id string) (*Response, error) {
	if id == "" {
		return nil, notification.ErrEmptyID
	}

	notif, err := g.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ToResponse(notif), nil
}
