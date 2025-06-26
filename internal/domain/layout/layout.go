package layout

import (
	"time"

	"github.com/lib/pq"
)

// Layout представляет шаблон уведомления
type Layout struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description,omitempty" db:"description"`
	Subject     string         `json:"subject" db:"subject"`
	Body        string         `json:"body" db:"body"`
	Type        string         `json:"type" db:"type"`
	Variables   pq.StringArray `json:"variables" db:"variables"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	IsActive    bool           `json:"is_active" db:"is_active"`
	Version     int            `json:"version" db:"version"`
}

// TemplateData содержит данные для рендеринга шаблона
type TemplateData struct {
	Variables map[string]interface{} `json:"variables"` // Данные для подстановки в шаблон
	Metadata  map[string]string      `json:"metadata"`  // Дополнительные метаданные
}

// ValidationErrors содержит ошибки валидации шаблона
type ValidationErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
