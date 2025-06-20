package list

import (
	"notification/internal/domain/layout"
)

// ToResponseList преобразует список доменных моделей в список Response
func ToResponseList(layouts []*layout.Layout) []*Response {
	responses := make([]*Response, 0, len(layouts))
	for _, l := range layouts {
		responses = append(responses, ToResponse(l))
	}
	return responses
}

// ToResponse преобразует одну доменную модель в Response
func ToResponse(l *layout.Layout) *Response {
	return &Response{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		Type:        l.Type,
		IsActive:    l.IsActive,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}
