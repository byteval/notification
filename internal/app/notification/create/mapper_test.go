package create

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"notification/internal/domain/notification"
)

func TestToDomain(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		expected *notification.Notification
	}{
		{
			name: "успешная конвертация с полными данными",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Тестовое уведомление",
				Data: map[string]interface{}{
					"name": "John",
					"age":  30,
				},
				Emails: []string{"test1@example.com", "test2@example.com"},
				Attachments: []Attachment{
					{
						Filename:    "file1.pdf",
						ContentType: "application/pdf",
						Size:        1024,
						FilePath:    "/tmp/file1.pdf",
					},
					{
						Filename:    "file2.jpg",
						ContentType: "image/jpeg",
						Size:        2048,
						FilePath:    "/tmp/file2.jpg",
					},
				},
			},
			expected: &notification.Notification{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Тестовое уведомление",
				Data: notification.JSONB{
					"name": "John",
					"age":  30,
				},
				NotificationReceivers: []notification.NotificationReceiver{
					{Email: "test1@example.com", Status: notification.StatusPending},
					{Email: "test2@example.com", Status: notification.StatusPending},
				},
				Attachments: []notification.Attachment{
					{FileName: "file1.pdf", OriginalName: "file1.pdf", ContentType: "application/pdf", Size: 1024, FilePath: "/tmp/file1.pdf"},
					{FileName: "file2.jpg", OriginalName: "file2.jpg", ContentType: "image/jpeg", Size: 2048, FilePath: "/tmp/file2.jpg"},
				},
			},
		},
		{
			name: "конвертация с минимальными данными",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440001",
				Title:    "Минимальное уведомление",
				Data:     map[string]interface{}{},
				Emails:   []string{"minimal@example.com"},
			},
			expected: &notification.Notification{
				LayoutID: "550e8400-e29b-41d4-a716-446655440001",
				Title:    "Минимальное уведомление",
				Data:     notification.JSONB{},
				NotificationReceivers: []notification.NotificationReceiver{
					{Email: "minimal@example.com", Status: notification.StatusPending},
				},
			},
		},
		{
			name: "конвертация с пустыми вложениями",
			request: Request{
				LayoutID:    "550e8400-e29b-41d4-a716-446655440002",
				Title:       "Уведомление без вложений",
				Data:        map[string]interface{}{},
				Emails:      []string{"noattach@example.com"},
				Attachments: []Attachment{},
			},
			expected: &notification.Notification{
				LayoutID: "550e8400-e29b-41d4-a716-446655440002",
				Title:    "Уведомление без вложений",
				Data:     notification.JSONB{},
				NotificationReceivers: []notification.NotificationReceiver{
					{Email: "noattach@example.com", Status: notification.StatusPending},
				},
				Attachments: []notification.Attachment{},
			},
		},
		{
			name: "конвертация с пустой строкой в Title",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440003",
				Title:    "",
				Data:     map[string]interface{}{},
				Emails:   []string{"empty@example.com"},
			},
			expected: &notification.Notification{
				LayoutID: "550e8400-e29b-41d4-a716-446655440003",
				Title:    "",
				Data:     notification.JSONB{},
				NotificationReceivers: []notification.NotificationReceiver{
					{Email: "empty@example.com", Status: notification.StatusPending},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomain(tt.request)

			// Проверяем, что нет ошибок
			require.NoError(t, err)
			require.NotNil(t, result)

			// Проверяем основные поля
			assert.Equal(t, tt.expected.LayoutID, result.LayoutID)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Data, result.Data)

			// Проверяем получателей
			require.Len(t, result.NotificationReceivers, len(tt.expected.NotificationReceivers))
			for i, expectedReceiver := range tt.expected.NotificationReceivers {
				assert.Equal(t, expectedReceiver.Email, result.NotificationReceivers[i].Email)
				assert.Equal(t, expectedReceiver.Status, result.NotificationReceivers[i].Status)
			}

			// Проверяем вложения
			require.Len(t, result.Attachments, len(tt.expected.Attachments))
			for i, expectedAttachment := range tt.expected.Attachments {
				assert.Equal(t, expectedAttachment.FileName, result.Attachments[i].FileName)
				assert.Equal(t, expectedAttachment.OriginalName, result.Attachments[i].OriginalName)
				assert.Equal(t, expectedAttachment.ContentType, result.Attachments[i].ContentType)
				assert.Equal(t, expectedAttachment.Size, result.Attachments[i].Size)
				assert.Equal(t, expectedAttachment.FilePath, result.Attachments[i].FilePath)
			}

			// Проверяем, что ID пустой (не установлен)
			assert.Empty(t, result.ID)

			// Проверяем, что CreatedAt установлен
			assert.False(t, result.CreatedAt.IsZero())

			// Проверяем, что CreatedAt примерно равен текущему времени
			now := time.Now()
			assert.WithinDuration(t, now, result.CreatedAt, 2*time.Second)
		})
	}
}
