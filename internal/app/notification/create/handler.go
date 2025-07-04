package create

import (
	"context"
	"fmt"
	"notification/internal/app/notification/common"
	"notification/internal/domain/notification/ports"
	"notification/pkg/logger"
	"time"
)

type NotificationCreator struct {
	repo     ports.NotificationRepository
	notifier *Notifier
	logger   logger.Logger
}

func NewCreator(
	repo ports.NotificationRepository,
	notifier *Notifier,
	logger logger.Logger,
) *NotificationCreator {
	return &NotificationCreator{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

func (s *NotificationCreator) Execute(ctx context.Context, req Request) (*common.Response, error) {
	req = *s.checkUniqueReceiver(&req)

	n, err := ToDomain(req)
	if err != nil {
		s.logger.Error("Failed to convert request to domain", "error", err)
		return nil, fmt.Errorf("failed to convert request to domain: %w", err)
	}

	n.CreatedAt = time.Now()

	created, err := s.repo.Create(ctx, n)

	if err != nil {
		s.logger.Error("Failed to save notification", "error", err)
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	s.notifier.SendNotificationAsync(n)

	return common.ToResponse(created), nil
}

// Оставляем только уникальных получателей
func (s *NotificationCreator) checkUniqueReceiver(req *Request) *Request {
	uniqueMap := make(map[string]bool, len(req.Emails))
	uniqueSlice := make([]string, 0, len(req.Emails))

	for _, email := range req.Emails {
		if _, exists := uniqueMap[email]; !exists {
			uniqueMap[email] = true
			uniqueSlice = append(uniqueSlice, email)
		}
	}

	req.Emails = uniqueSlice
	return req
}
