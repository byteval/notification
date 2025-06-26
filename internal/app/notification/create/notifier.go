package create

import (
	"context"
	"time"

	layoutPorts "notification/internal/domain/layout/ports"
	"notification/internal/domain/notification"
	"notification/internal/domain/notification/ports"
	"notification/pkg/logger"

	"github.com/alitto/pond/v2"
)

type Notifier struct {
	repo             ports.NotificationRepository
	layoutRepo       layoutPorts.LayoutRepository
	smtpSender       ports.SMTPSender
	logger           logger.Logger
	pool             pond.Pool
	templateRenderer *TemplateRenderer
}

// NewNotifier создает новый Notifier и возвращает указатель на него
func NewNotifier(
	repo ports.NotificationRepository,
	layoutRepo layoutPorts.LayoutRepository,
	smtpSender ports.SMTPSender,
	logger logger.Logger,
	workers int,
) *Notifier {
	return &Notifier{
		repo:             repo,
		layoutRepo:       layoutRepo,
		smtpSender:       smtpSender,
		logger:           logger,
		pool:             pond.NewPool(workers),
		templateRenderer: NewTemplateRenderer(layoutRepo, logger),
	}
}

// Отправляем уведомление через пул воркеров
func (s *Notifier) SendNotificationAsync(n *notification.Notification) {
	s.pool.Submit(func() {
		// Добавляем recover для отлова паник
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Recovered from panic in sender",
					"panic", r,
					"notification_id", n.ID)

				// Пытаемся сохранить статус "failed" даже при панике
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				s.updateStatus(ctx, n, notification.StatusFailed)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Рендер заголовка и текста уведомления
		title, content, err := s.templateRenderer.Render(ctx, n)
		if err != nil {
			s.logFailed(ctx, err, n)
			return
		}

		sendErr := s.smtpSender.Send(n.Receiver, title, content)
		if sendErr == nil {
			s.updateStatus(ctx, n, notification.StatusSent)
			return
		}

		s.logFailed(ctx, sendErr, n)
	})
}

// Логирование ошибок
func (s *Notifier) logFailed(ctx context.Context, sendErr error, n *notification.Notification) {
	s.logger.Error("Send failed", "error", sendErr, "notification_id", n.ID)
	s.updateStatus(ctx, n, notification.StatusFailed)
}

// Обновление статуса уведомления
func (s *Notifier) updateStatus(ctx context.Context, n *notification.Notification, status notification.Status) {
	s.logger.Info("Изменение статуса уведомления",
		"notification_id", n.ID,
		"status", status)

	if err := s.repo.UpdateStatus(ctx, n.ID, status); err != nil {
		s.logger.Error("Failed to update status",
			"error", err,
			"notification_id", n.ID,
			"status", status)
	}
}

// WaitForCompletion ожидает завершения всех отправок
func (n *Notifier) WaitForCompletion() {
	n.pool.StopAndWait()
}
