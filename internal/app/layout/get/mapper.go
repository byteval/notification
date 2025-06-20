package get

import (
	"notification/internal/domain/layout"
)

// ToResponse преобразует доменную модель в Response
func ToResponse(l *layout.Layout) *Response {
	return &Response{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		Subject:     l.Subject,
		Body:        l.Body,
		Type:        l.Type,
		Variables:   l.Variables,
		IsActive:    l.IsActive,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}
