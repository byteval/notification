package create

import (
	"context"
	"fmt"
	"notification/internal/domain/notification/ports"
	"notification/pkg/logger"
	"time"
)

type Creator struct {
	repo     ports.NotificationRepository
	notifier *Notifier
	logger   logger.Logger
}

func NewCreator(
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
	req = *s.checkUniqueReceiver(&req)

	n, receivers, err := ToDomain(req)
	if err != nil {
		s.logger.Error("Failed to convert request to domain", "error", err)
		return nil, fmt.Errorf("failed to convert request to domain: %w", err)
	}

	n.CreatedAt = time.Now()

	created, err := s.repo.CreateWithReceivers(ctx, n, receivers)

	if err != nil {
		s.logger.Error("Failed to save notification", "error", err)
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	s.notifier.SendNotificationAsync(n)

	return ToResponse(created), nil
}

// Оставляем только уникальных получателей
func (s *Creator) checkUniqueReceiver(req *Request) *Request {
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
