package get

import (
	"context"
	"notification/internal/app/notification/common"
	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
)

type NotificationGetter struct {
	repo ports.NotificationRepository
}

func NewGetter(repo ports.NotificationRepository) *NotificationGetter {
	return &NotificationGetter{repo: repo}
}

// Execute возвращает уведомление по ID
func (g *NotificationGetter) Execute(ctx context.Context, id string) (*common.Response, error) {
	if id == "" {
		return nil, notification.ErrEmptyID
	}

	n, err := g.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return common.ToResponse(n), nil
}
