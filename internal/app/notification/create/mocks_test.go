package create

import (
	"context"
	"notification/internal/domain/layout"
	"notification/internal/domain/notification"

	"github.com/stretchr/testify/mock"
)

// MockNotificationRepository - мок для NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, n *notification.Notification) (*notification.Notification, error) {
	args := m.Called(ctx, n)
	return args.Get(0).(*notification.Notification), args.Error(1)
}

func (m *MockNotificationRepository) GetByID(ctx context.Context, id string) (*notification.Notification, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*notification.Notification), args.Error(1)
}

func (m *MockNotificationRepository) UpdateStatus(ctx context.Context, id string, status notification.Status) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockLayoutRepository - мок для LayoutRepository
type MockLayoutRepository struct {
	mock.Mock
}

func (m *MockLayoutRepository) Create(ctx context.Context, l *layout.Layout) (*layout.Layout, error) {
	args := m.Called(ctx, l)
	return args.Get(0).(*layout.Layout), args.Error(1)
}

func (m *MockLayoutRepository) GetByID(ctx context.Context, id string) (*layout.Layout, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*layout.Layout), args.Error(1)
}

func (m *MockLayoutRepository) List(ctx context.Context) ([]*layout.Layout, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*layout.Layout), args.Error(1)
}

func (m *MockLayoutRepository) Update(ctx context.Context, l *layout.Layout) (*layout.Layout, error) {
	args := m.Called(ctx, l)
	return args.Get(0).(*layout.Layout), args.Error(1)
}

func (m *MockLayoutRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLayoutRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// MockSMTPSender - мок для SMTPSender
type MockSMTPSender struct {
	mock.Mock
}

func (m *MockSMTPSender) Send(id string, to string, subject string, message string, attachments map[string]string) error {
	args := m.Called(id, to, subject, message, attachments)
	return args.Error(0)
}

// createMockNotification - вспомогательная функция для создания тестового уведомления
func createMockNotification() *notification.Notification {
	return &notification.Notification{
		ID:       "test-notification-id",
		LayoutID: "test-layout-id",
		Title:    "Test Notification",
		Data: notification.JSONB{
			"name": "John",
			"age":  30,
		},
		NotificationReceivers: []notification.NotificationReceiver{
			{
				ID:     "receiver-1",
				Email:  "test1@example.com",
				Status: notification.StatusPending,
			},
			{
				ID:     "receiver-2",
				Email:  "test2@example.com",
				Status: notification.StatusPending,
			},
		},
		Attachments: []notification.Attachment{
			{
				FileName:     "test.pdf",
				OriginalName: "test.pdf",
				ContentType:  "application/pdf",
				Size:         1024,
				FilePath:     "/tmp/test.pdf",
			},
		},
	}
}
