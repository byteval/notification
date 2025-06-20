package update

import (
	"context"

	"notification/internal/domain/layout/ports"
	"notification/pkg/logger"
)

type Updater struct {
	repo   ports.LayoutRepository
	logger logger.Logger
}

func NewUpdater(repo ports.LayoutRepository, logger logger.Logger) *Updater {
	return &Updater{repo: repo, logger: logger}
}

func (uc *Updater) Execute(ctx context.Context, req Request) (*Response, error) {
	existing, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Subject != "" {
		existing.Subject = req.Subject
	}
	if req.Body != "" {
		existing.Body = req.Body
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.Version = req.Version

	updated, err := uc.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return ToResponse(updated), nil
}
