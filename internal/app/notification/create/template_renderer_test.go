package create

import (
	"context"
	"testing"

	"notification/internal/domain/layout"
	"notification/internal/domain/notification"
	"notification/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTemplateRenderer_Render(t *testing.T) {
	tests := []struct {
		name          string
		notification  *notification.Notification
		layout        *layout.Layout
		layoutError   error
		expectedTitle string
		expectedBody  string
		expectedError string
	}{
		{
			name: "успешный рендеринг шаблона",
			notification: &notification.Notification{
				LayoutID: "test-layout-id",
				Data: notification.JSONB{
					"name": "John",
					"age":  30,
					"city": "Moscow",
				},
			},
			layout: &layout.Layout{
				ID:       "test-layout-id",
				Subject:  "Привет, {{.name}}!",
				Body:     "Вам {{.age}} лет и вы живете в {{.city}}.",
				IsActive: true,
			},
			expectedTitle: "Привет, John!",
			expectedBody:  "Вам 30 лет и вы живете в Moscow.",
		},
		{
			name: "рендеринг с числовыми данными",
			notification: &notification.Notification{
				LayoutID: "test-layout-id",
				Data: notification.JSONB{
					"count": 5,
					"price": 99.99,
				},
			},
			layout: &layout.Layout{
				ID:       "test-layout-id",
				Subject:  "У вас {{.count}} новых сообщений",
				Body:     "Общая стоимость: ${{.price}}",
				IsActive: true,
			},
			expectedTitle: "У вас 5 новых сообщений",
			expectedBody:  "Общая стоимость: $99.99",
		},
		{
			name: "рендеринг с пустыми данными",
			notification: &notification.Notification{
				LayoutID: "test-layout-id",
				Data:     notification.JSONB{},
			},
			layout: &layout.Layout{
				ID:       "test-layout-id",
				Subject:  "Статический заголовок",
				Body:     "Статическое содержимое",
				IsActive: true,
			},
			expectedTitle: "Статический заголовок",
			expectedBody:  "Статическое содержимое",
		},
		{
			name: "ошибка получения шаблона",
			notification: &notification.Notification{
				LayoutID: "invalid-layout-id",
				Data:     notification.JSONB{},
			},
			layoutError:   assert.AnError,
			expectedError: "failed to get layout: assert.AnError general error for testing",
		},
		{
			name: "неактивный шаблон",
			notification: &notification.Notification{
				LayoutID: "inactive-layout-id",
				Data:     notification.JSONB{},
			},
			layout: &layout.Layout{
				ID:       "inactive-layout-id",
				Subject:  "Заголовок",
				Body:     "Содержимое",
				IsActive: false,
			},
			expectedError: "layout inactive-layout-id is not active",
		},
		{
			name: "рендеринг с отсутствующими полями",
			notification: &notification.Notification{
				LayoutID: "test-layout-id",
				Data: notification.JSONB{
					"name": "John",
				},
			},
			layout: &layout.Layout{
				ID:       "test-layout-id",
				Subject:  "Привет, {{.name}}! {{.invalid_field}}",
				Body:     "Валидное содержимое {{.missing_field}}",
				IsActive: true,
			},
			expectedTitle: "Привет, John! <no value>",
			expectedBody:  "Валидное содержимое <no value>",
		},
	}

	// Создаем моки один раз для всех тестов
	mockLayoutRepo := new(MockLayoutRepository)
	mockLogger := logger.New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем мок перед каждым тестом
			mockLayoutRepo.ExpectedCalls = nil

			// Настраиваем ожидания мока
			if tt.layoutError != nil {
				mockLayoutRepo.On("GetByID", mock.Anything, tt.notification.LayoutID).Return((*layout.Layout)(nil), tt.layoutError)
			} else {
				mockLayoutRepo.On("GetByID", mock.Anything, tt.notification.LayoutID).Return(tt.layout, nil)
			}

			// Создаем TemplateRenderer
			renderer := NewTemplateRenderer(mockLayoutRepo, mockLogger)

			// Выполняем рендеринг
			title, body, err := renderer.Render(context.Background(), tt.notification)

			// Проверяем результаты
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, title)
				assert.Empty(t, body)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTitle, title)
				assert.Equal(t, tt.expectedBody, body)
			}

			// Проверяем, что все ожидания мока выполнены
			mockLayoutRepo.AssertExpectations(t)
		})
	}
}

func TestTemplateRenderer_renderTemplate(t *testing.T) {
	tests := []struct {
		name          string
		templateText  string
		data          map[string]interface{}
		expected      string
		expectedError string
	}{
		{
			name:         "простой шаблон",
			templateText: "Привет, {{.name}}!",
			data: map[string]interface{}{
				"name": "Мир",
			},
			expected: "Привет, Мир!",
		},
		{
			name:         "шаблон с условиями",
			templateText: "{{if .isActive}}Активен{{else}}Неактивен{{end}}",
			data: map[string]interface{}{
				"isActive": true,
			},
			expected: "Активен",
		},
		{
			name:         "шаблон с циклом",
			templateText: "{{range .items}}{{.}}{{end}}",
			data: map[string]interface{}{
				"items": []string{"a", "b", "c"},
			},
			expected: "abc",
		},
		{
			name:         "шаблон с функциями",
			templateText: "{{len .name}} символов",
			data: map[string]interface{}{
				"name": "Тест",
			},
			expected: "8 символов",
		},
		{
			name:         "некорректный синтаксис шаблона",
			templateText: "{{.name",
			data: map[string]interface{}{
				"name": "Тест",
			},
			expectedError: "template parsing error: template: notification:1: unclosed action",
		},
		{
			name:         "отсутствующее поле",
			templateText: "{{.missing_field}}",
			data: map[string]interface{}{
				"name": "Тест",
			},
			expected: "<no value>",
		},
		{
			name:         "пустой шаблон",
			templateText: "",
			data:         map[string]interface{}{},
			expected:     "",
		},
		{
			name:         "шаблон без переменных",
			templateText: "Статический текст",
			data:         map[string]interface{}{},
			expected:     "Статический текст",
		},
	}

	// Создаем моки один раз для всех тестов
	mockLayoutRepo := new(MockLayoutRepository)
	mockLogger := logger.New(nil)
	renderer := NewTemplateRenderer(mockLayoutRepo, mockLogger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Выполняем рендеринг шаблона
			result, err := renderer.renderTemplate(tt.templateText, tt.data)

			// Проверяем результаты
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNewTemplateRenderer(t *testing.T) {
	mockLayoutRepo := new(MockLayoutRepository)
	mockLogger := logger.New(nil)

	renderer := NewTemplateRenderer(mockLayoutRepo, mockLogger)

	assert.NotNil(t, renderer)
	assert.Equal(t, mockLayoutRepo, renderer.layoutRepo)
	assert.Equal(t, mockLogger, renderer.logger)
}
