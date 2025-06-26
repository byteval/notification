package update

import (
	"notification/internal/app/layout/common"
	"notification/internal/domain/layout"
)

// ToDomain применяет изменения из Request к существующей доменной модели
func ToDomain(r common.Request, existing *layout.Layout) (*layout.Layout, error) {
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
func ToResponse(l *layout.Layout) *common.Response {
	return &common.Response{
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
		Version:     l.Version,
	}
}
