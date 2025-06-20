package update

import "time"

// Request определяет входные данные для обновления шаблона
// swagger:model UpdateLayoutRequest
type Request struct {
	// ID шаблона
	// required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" validate:"required,uuid4"`

	// Версия для оптимистичной блокировки
	// required: true
	// example: 1
	Version int `json:"version" validate:"required,min=1"`

	Name        string   `json:"name" validate:"omitempty,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty,max=500"`
	Subject     string   `json:"subject" validate:"omitempty,min=3,max=200"`
	Body        string   `json:"body" validate:"omitempty,min=10"`
	Type        string   `json:"type" validate:"omitempty"`
	Variables   []string `json:"variables" validate:"omitempty,dive,required"`
	IsActive    *bool    `json:"is_active"`
}

// Response определяет ответ после обновления шаблона
// swagger:model UpdateLayoutResponse
type Response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}
