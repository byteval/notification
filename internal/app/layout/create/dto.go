package create

import "time"

// Request определяет входные данные для создания шаблона
// swagger:model CreateLayoutRequest
type Request struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"max=500"`
	Subject     string   `json:"subject" validate:"required,min=3,max=200"`
	Body        string   `json:"body" validate:"required,min=10"`
	Type        string   `json:"type" validate:"required"`
	Variables   []string `json:"variables" validate:"dive,required"`
	IsActive    bool     `json:"is_active"`
}

// Response определяет ответ после создания шаблона
// swagger:model CreateLayoutResponse
type Response struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
