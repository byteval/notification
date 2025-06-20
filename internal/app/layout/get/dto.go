package get

import "time"

// Request определяет входные данные для получения шаблона
// swagger:parameters getLayout
type Request struct {
	// ID шаблона
	// in: path
	// required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" validate:"required,uuid4"`
}

// Response определяет ответ с данными шаблона
// swagger:model GetLayoutResponse
type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	Type        string    `json:"type"`
	Variables   []string  `json:"variables"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"`
}
