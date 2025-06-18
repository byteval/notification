package create

import (
	"context"
	"fmt"
	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
	"notification/pkg/logger"
	"time"
)

type Creator struct {
	repo     ports.NotificationRepository
	notifier *Notifier
	logger   logger.Logger
}

func NewNotificationService(
	repo ports.NotificationRepository,
	notifier *Notifier,
	logger logger.Logger,
) *Creator {
	return &Creator{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

func (s *Creator) CreateNotification(ctx context.Context, req Request) (*Response, error) {
	n, err := req.ToDomain()
	if err != nil {
		s.logger.Error("Failed to convert request to domain", "error", err)
		return nil, fmt.Errorf("failed to convert request to domain: %w", err)
	}

	n.CreatedAt = time.Now()
	n.Status = notification.StatusPending

	created, err := s.repo.Create(ctx, n)

	if err != nil {
		s.logger.Error("Failed to save notification", "error", err)
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	s.notifier.SendNotificationAsync(n)

	return ToResponse(created), nil
}
