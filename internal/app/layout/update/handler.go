package update

import (
	"context"

	"notification/internal/app/layout/common"
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

func (uc *Updater) Execute(ctx context.Context, req common.Request) (*common.Response, error) {
	existing, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	l, err := ToDomain(req, existing)
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.Update(ctx, l)
	if err != nil {
		return nil, err
	}

	return ToResponse(updated), nil
}
