package create

import (
	"testing"

	"notification/internal/domain/layout"
	"notification/internal/domain/notification"
	"notification/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotifier_SendNotificationAsync(t *testing.T) {
	t.Run("успешная отправка уведомления", func(t *testing.T) {
		mockRepo := new(MockNotificationRepository)
		mockLayoutRepo := new(MockLayoutRepository)
		mockSMTP := new(MockSMTPSender)
		mockLogger := logger.New(nil)
		notifier := NewNotifier(mockRepo, mockLayoutRepo, mockSMTP, mockLogger, 2)

		notif := &notification.Notification{
			ID:       "test-notification-id",
			LayoutID: "test-layout-id",
			Title:    "Test Notification",
			Data:     notification.JSONB{"name": "John"},
			NotificationReceivers: []notification.NotificationReceiver{
				{ID: "receiver-1", Email: "test1@example.com"},
				{ID: "receiver-2", Email: "test2@example.com"},
			},
		}
		l := &layout.Layout{
			ID:       "test-layout-id",
			Subject:  "Hello {{.name}}",
			Body:     "Welcome {{.name}}!",
			IsActive: true,
		}
		mockLayoutRepo.On("GetByID", mock.Anything, "test-layout-id").Return(l, nil)
		mockSMTP.On("Send", "receiver-1", "test1@example.com", "Hello John", "Welcome John!", mock.Anything).Return(nil)
		mockSMTP.On("Send", "receiver-2", "test2@example.com", "Hello John", "Welcome John!", mock.Anything).Return(nil)
		mockRepo.On("UpdateStatus", mock.Anything, "receiver-1", notification.StatusSent).Return(nil)
		mockRepo.On("UpdateStatus", mock.Anything, "receiver-2", notification.StatusSent).Return(nil)

		notifier.SendNotificationAsync(notif)
		notifier.WaitForCompletion()

		mockLayoutRepo.AssertExpectations(t)
		mockSMTP.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка рендеринга шаблона", func(t *testing.T) {
		mockRepo := new(MockNotificationRepository)
		mockLayoutRepo := new(MockLayoutRepository)
		mockSMTP := new(MockSMTPSender)
		mockLogger := logger.New(nil)
		notifier := NewNotifier(mockRepo, mockLayoutRepo, mockSMTP, mockLogger, 1)

		notif := &notification.Notification{
			ID:       "test-notification-id",
			LayoutID: "invalid-layout-id",
			Title:    "Test Notification",
			Data:     notification.JSONB{"name": "John"},
			NotificationReceivers: []notification.NotificationReceiver{
				{ID: "receiver-1", Email: "test@example.com"},
			},
		}
		mockLayoutRepo.On("GetByID", mock.Anything, "invalid-layout-id").Return((*layout.Layout)(nil), assert.AnError)
		mockRepo.On("UpdateStatus", mock.Anything, "receiver-1", notification.StatusFailed).Return(nil)

		notifier.SendNotificationAsync(notif)
		notifier.WaitForCompletion()

		mockLayoutRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка отправки SMTP", func(t *testing.T) {
		mockRepo := new(MockNotificationRepository)
		mockLayoutRepo := new(MockLayoutRepository)
		mockSMTP := new(MockSMTPSender)
		mockLogger := logger.New(nil)
		notifier := NewNotifier(mockRepo, mockLayoutRepo, mockSMTP, mockLogger, 1)

		notif := &notification.Notification{
			ID:       "test-notification-id",
			LayoutID: "test-layout-id",
			Title:    "Test Notification",
			Data:     notification.JSONB{"name": "John"},
			NotificationReceivers: []notification.NotificationReceiver{
				{ID: "receiver-1", Email: "test@example.com"},
			},
		}
		l := &layout.Layout{
			ID:       "test-layout-id",
			Subject:  "Hello {{.name}}",
			Body:     "Welcome {{.name}}!",
			IsActive: true,
		}
		mockLayoutRepo.On("GetByID", mock.Anything, "test-layout-id").Return(l, nil)
		mockSMTP.On("Send", "receiver-1", "test@example.com", "Hello John", "Welcome John!", mock.Anything).Return(assert.AnError)
		mockRepo.On("UpdateStatus", mock.Anything, "receiver-1", notification.StatusFailed).Return(nil)

		notifier.SendNotificationAsync(notif)
		notifier.WaitForCompletion()

		mockLayoutRepo.AssertExpectations(t)
		mockSMTP.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestNotifier_getAttachments(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	mockLayoutRepo := new(MockLayoutRepository)
	mockSMTP := new(MockSMTPSender)
	mockLogger := logger.New(nil)

	notifier := NewNotifier(mockRepo, mockLayoutRepo, mockSMTP, mockLogger, 1)

	t.Run("получение вложений", func(t *testing.T) {
		notif := &notification.Notification{
			ID: "test-notification-id",
			Attachments: []notification.Attachment{
				{
					FileName: "test1.pdf",
					FilePath: "/tmp/test1.pdf",
				},
				{
					FileName: "test2.jpg",
					FilePath: "/tmp/test2.jpg",
				},
			},
		}

		attachments := notifier.getAttachments(notif)

		expected := map[string]string{
			"test1.pdf": "/tmp/test1.pdf",
			"test2.jpg": "/tmp/test2.jpg",
		}

		assert.Equal(t, expected, attachments)
	})

	t.Run("пустой список вложений", func(t *testing.T) {
		notif := &notification.Notification{
			ID:          "test-notification-id",
			Attachments: []notification.Attachment{},
		}

		attachments := notifier.getAttachments(notif)

		assert.Empty(t, attachments)
	})
}
