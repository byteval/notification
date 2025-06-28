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
	for _, nr := range n.NotificationReceivers {
		s.pool.Submit(func() {
			// Добавляем recover для отлова паник
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("Recovered from panic in sender",
						"panic", r,
						"notification_id", n.ID,
						"notification_receiver_id", nr.ID,
					)

					// Пытаемся сохранить статус "failed" при панике
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					s.updateStatus(ctx, nr, notification.StatusFailed)
				}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Рендер заголовка и текста уведомления
			title, content, err := s.templateRenderer.Render(ctx, n)
			if err != nil {
				s.logFailed(ctx, err, nr)
				return
			}

			sendErr := s.smtpSender.Send(nr.ID, nr.Email, title, content, s.getAttachments(n))
			if sendErr == nil {
				s.updateStatus(ctx, nr, notification.StatusSent)
				return
			}

			s.logFailed(ctx, sendErr, nr)
		})
	}
}

func (s *Notifier) getAttachments(n *notification.Notification) map[string]string {
	resultMap := make(map[string]string, len(n.Attachments))

	for _, attachment := range n.Attachments {
		resultMap[attachment.FileName] = attachment.FilePath
	}

	return resultMap
}

// Логирование ошибок
func (s *Notifier) logFailed(ctx context.Context, sendErr error, nr notification.NotificationReceiver) {
	s.logger.Error("Send failed", "error", sendErr, "notification_receiver_id", nr.ID)
	s.updateStatus(ctx, nr, notification.StatusFailed)
}

// Обновление статуса уведомления
func (s *Notifier) updateStatus(ctx context.Context, nr notification.NotificationReceiver, status notification.Status) {
	s.logger.Info("Изменение статуса уведомления",
		"notification_receiver_id", nr.ID,
		"status", status)

	if err := s.repo.UpdateStatus(ctx, nr.ID, status); err != nil {
		s.logger.Error("Failed to update status",
			"error", err,
			"notification_receiver_id", nr.ID,
			"status", status)
	}
}

// WaitForCompletion ожидает завершения всех отправок
func (n *Notifier) WaitForCompletion() {
	n.pool.StopAndWait()
}
