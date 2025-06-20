package delete

import (
	"context"

	"notification/internal/domain/layout/ports"
	"notification/pkg/logger"
)

type Deleter struct {
	repo   ports.LayoutRepository
	logger logger.Logger
}

func NewDeleter(repo ports.LayoutRepository, logger logger.Logger) *Deleter {
	return &Deleter{repo: repo, logger: logger}
}

func (uc *Deleter) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
