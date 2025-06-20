package create

import (
	"notification/internal/domain/layout"
	"time"
)

// ToDomain преобразует Request в доменную модель Layout
func (r Request) ToDomain() (*layout.Layout, error) {
	return &layout.Layout{
		Name:        r.Name,
		Description: r.Description,
		Subject:     r.Subject,
		Body:        r.Body,
		Type:        r.Type,
		Variables:   r.Variables,
		IsActive:    r.IsActive,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Version:     1,
	}, nil
}

// ToResponse преобразует доменную модель в Response
func ToResponse(l *layout.Layout) *Response {
	return &Response{
		ID:        l.ID,
		CreatedAt: l.CreatedAt,
	}
}
