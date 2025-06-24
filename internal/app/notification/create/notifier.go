package create

import (
	"context"
	"sync"
	"time"

	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
	"notification/pkg/logger"
)

type Notifier struct {
	repo       ports.NotificationRepository
	smtpSender ports.SMTPSender
	logger     logger.Logger
	wg         sync.WaitGroup
	mu         sync.Mutex
}

// NewNotifier создает новый Notifier и возвращает указатель на него
func NewNotifier(
	repo ports.NotificationRepository,
	smtpSender ports.SMTPSender,
	logger logger.Logger,
) *Notifier {
	return &Notifier{
		repo:   repo,
		logger: logger,
		// wg и mu инициализируются автоматически
	}
}

// Асинхронная отправка
func (s *Notifier) SendNotificationAsync(n *notification.Notification) {
	s.mu.Lock()
	s.wg.Add(1)
	s.mu.Unlock()

	go func(s *Notifier, n *notification.Notification) {
		defer func() {
			s.mu.Lock()
			s.wg.Done()
			s.mu.Unlock()

			if r := recover(); r != nil {
				s.logger.Error("Recovered from panic in sender",
					"panic", r,
					"notification_id", n.ID)
			}
		}()

		s.sendNotification(n)
	}(s, n)
}

// Отправка уведомления
func (s *Notifier) sendNotification(n *notification.Notification) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.logger.Info("Отправка уведомления")
	s.smtpSender.Send(n.Receiver, n.Title, n.Content)

	s.updateStatus(ctx, n, notification.StatusSent)
}

// Обновление статуса уведомления
func (s *Notifier) updateStatus(ctx context.Context, n *notification.Notification, status notification.Status) {
	s.logger.Info("Изменение статуса уведомления",
		"notification_id", n.ID,
		"status", status)

	n.Status = status
	if err := s.repo.UpdateStatus(ctx, n.ID, status); err != nil {
		s.logger.Error("Failed to update status",
			"error", err,
			"notification_id", n.ID,
			"status", status)
	}
}

// WaitForCompletion ожидает завершения всех отправок
func (s *Notifier) WaitForCompletion() {
	s.wg.Wait()
}
