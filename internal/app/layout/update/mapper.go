package update

import (
	"notification/internal/domain/layout"
)

// ToDomain применяет изменения из Request к существующей доменной модели
func (r Request) ToDomain(existing *layout.Layout) (*layout.Layout, error) {
	if r.Name != "" {
		existing.Name = r.Name
	}
	if r.Description != "" {
		existing.Description = r.Description
	}
	if r.Subject != "" {
		existing.Subject = r.Subject
	}
	if r.Body != "" {
		existing.Body = r.Body
	}
	if r.IsActive != nil {
		existing.IsActive = *r.IsActive
	}
	existing.Version = r.Version

	return existing, nil
}

// ToResponse преобразует доменную модель в Response
func ToResponse(l *layout.Layout) *Response {
	return &Response{
		UpdatedAt: l.UpdatedAt,
		Version:   l.Version,
	}
}
