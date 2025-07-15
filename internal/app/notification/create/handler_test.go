package create

import (
	"context"
	"errors"
	"testing"

	"notification/internal/domain/notification"
	"notification/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNotificationCreator_Execute(t *testing.T) {
	tests := []struct {
		name           string
		request        Request
		prepareMocks   func(repo *MockNotificationRepository, notifier *Notifier)
		expectedError  string
		expectedEmails []string
	}{
		{
			name: "успешное создание уведомления",
			request: Request{
				LayoutID: "test-layout-id",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"name": "John"},
				Emails:   []string{"test1@example.com", "test2@example.com"},
			},
			prepareMocks: func(repo *MockNotificationRepository, notifier *Notifier) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*notification.Notification")).Return(createMockNotification(), nil)
				notifier.pool.StopAndWait() // чтобы не было гонок
			},
			expectedEmails: []string{"test1@example.com", "test2@example.com"},
		},
		{
			name: "ошибка сохранения уведомления",
			request: Request{
				LayoutID: "test-layout-id",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"name": "John"},
				Emails:   []string{"test@example.com"},
			},
			prepareMocks: func(repo *MockNotificationRepository, notifier *Notifier) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*notification.Notification")).Return((*notification.Notification)(nil), errors.New("db error"))
				notifier.pool.StopAndWait()
			},
			expectedError: "failed to save notification: db error",
		},
		{
			name: "удаление дубликатов email",
			request: Request{
				LayoutID: "test-layout-id",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"name": "John"},
				Emails:   []string{"a@example.com", "b@example.com", "a@example.com"},
			},
			prepareMocks: func(repo *MockNotificationRepository, notifier *Notifier) {
				resp := &notification.Notification{
					ID:       "test-notification-id",
					LayoutID: "test-layout-id",
					Title:    "Test Notification",
					Data:     notification.JSONB{"name": "John"},
					NotificationReceivers: []notification.NotificationReceiver{
						{ID: "r1", Email: "a@example.com", Status: notification.StatusPending},
						{ID: "r2", Email: "b@example.com", Status: notification.StatusPending},
					},
				}
				repo.On("Create", mock.Anything, mock.AnythingOfType("*notification.Notification")).Return(resp, nil)
				notifier.pool.StopAndWait()
			},
			expectedEmails: []string{"a@example.com", "b@example.com"},
		},
	}

	// Создаём моки и объекты один раз
	mockRepo := new(MockNotificationRepository)
	mockLayoutRepo := new(MockLayoutRepository)
	mockSMTP := new(MockSMTPSender)
	mockLogger := logger.New(nil)
	notifier := NewNotifier(mockRepo, mockLayoutRepo, mockSMTP, mockLogger, 1)
	creator := NewCreator(mockRepo, notifier, mockLogger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем ожидания моков
			mockRepo.ExpectedCalls = nil
			mockLayoutRepo.ExpectedCalls = nil
			mockSMTP.ExpectedCalls = nil

			if tt.prepareMocks != nil {
				tt.prepareMocks(mockRepo, notifier)
			}

			resp, err := creator.Execute(context.Background(), tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				if len(tt.expectedEmails) > 0 {
					// Проверяем, что дубликаты удалены
					actualEmails := make([]string, len(resp.Receivers))
					for i, receiver := range resp.Receivers {
						actualEmails[i] = receiver.Email
					}
					assert.ElementsMatch(t, tt.expectedEmails, actualEmails)
				}
			}
		})
	}
}

func TestNotificationCreator_checkUniqueReceiver(t *testing.T) {
	creator := &NotificationCreator{}
	in := &Request{
		Emails: []string{"a@example.com", "b@example.com", "a@example.com", "c@example.com", "b@example.com"},
	}
	out := creator.checkUniqueReceiver(in)
	assert.ElementsMatch(t, []string{"a@example.com", "b@example.com", "c@example.com"}, out.Emails)
}
