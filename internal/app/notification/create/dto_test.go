package create

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     Request
		expectError bool
		errorField  string // поле, которое должно вызвать ошибку
	}{
		{
			name: "валидный запрос",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com", "user@example.com"},
			},
			expectError: false,
		},
		{
			name: "отсутствует layout_id",
			request: Request{
				Title:  "Test Notification",
				Data:   map[string]interface{}{"key": "value"},
				Emails: []string{"test@example.com"},
			},
			expectError: true,
			errorField:  "LayoutID",
		},
		{
			name: "некорректный UUID layout_id",
			request: Request{
				LayoutID: "invalid-uuid",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com"},
			},
			expectError: true,
			errorField:  "LayoutID",
		},
		{
			name: "отсутствует title",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com"},
			},
			expectError: true,
			errorField:  "Title",
		},
		{
			name: "пустой title",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com"},
			},
			expectError: true,
			errorField:  "Title",
		},
		{
			name: "отсутствуют emails",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
			},
			expectError: true,
			errorField:  "Emails",
		},
		{
			name: "пустой список emails",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{},
			},
			expectError: true,
			errorField:  "Emails",
		},
		{
			name: "некорректный email",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"invalid-email"},
			},
			expectError: true,
			errorField:  "Emails",
		},
		{
			name: "один корректный, один некорректный email",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com", "invalid-email"},
			},
			expectError: true,
			errorField:  "Emails",
		},
		{
			name: "валидный запрос с пустыми данными",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{},
				Emails:   []string{"test@example.com"},
			},
			expectError: false,
		},
		{
			name: "валидный запрос с nil данными",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     nil,
				Emails:   []string{"test@example.com"},
			},
			expectError: false,
		},
		{
			name: "валидный запрос с вложениями",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"test@example.com"},
				Attachments: []Attachment{
					{
						Filename:    "test.pdf",
						ContentType: "application/pdf",
						Size:        1024,
					},
				},
			},
			expectError: false,
		},
		{
			name: "несколько корректных email адресов",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"user1@example.com", "user2@test.org", "admin@company.co.uk"},
			},
			expectError: false,
		},
		{
			name: "email с поддоменами",
			request: Request{
				LayoutID: "550e8400-e29b-41d4-a716-446655440000",
				Title:    "Test Notification",
				Data:     map[string]interface{}{"key": "value"},
				Emails:   []string{"user@subdomain.example.com"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorField != "" {
					// Проверяем, что ошибка связана с ожидаемым полем
					validationErrors, ok := err.(validator.ValidationErrors)
					if ok {
						found := false
						for _, fieldError := range validationErrors {
							// Для массивов валидатор возвращает ошибки в формате "Emails[0]", "Emails[1]" и т.д.
							// Проверяем, начинается ли поле с ожидаемого имени
							if fieldError.Field() == tt.errorField ||
								(len(fieldError.Field()) > len(tt.errorField) &&
									fieldError.Field()[:len(tt.errorField)] == tt.errorField &&
									fieldError.Field()[len(tt.errorField)] == '[') {
								found = true
								break
							}
						}
						assert.True(t, found, "Ошибка должна быть связана с полем %s", tt.errorField)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAttachment_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		attachment  Attachment
		expectError bool
	}{
		{
			name: "валидное вложение",
			attachment: Attachment{
				Filename:    "test.pdf",
				ContentType: "application/pdf",
				Size:        1024,
				FilePath:    "/path/to/file",
			},
			expectError: false,
		},
		{
			name: "вложение с пустым именем файла",
			attachment: Attachment{
				Filename:    "",
				ContentType: "application/pdf",
				Size:        1024,
			},
			expectError: false, // Attachment не имеет валидационных тегов
		},
		{
			name: "вложение с нулевым размером",
			attachment: Attachment{
				Filename:    "test.pdf",
				ContentType: "application/pdf",
				Size:        0,
			},
			expectError: false, // Attachment не имеет валидационных тегов
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.attachment)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
