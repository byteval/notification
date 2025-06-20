// internal/app/layout/list/dto.go
package list

import "time"

type Request struct {
	Limit  int
	Offset int
}

// Response определяет элемент списка шаблонов
type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
