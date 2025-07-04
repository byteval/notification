package get

import (
	"context"

	"notification/internal/app/layout/common"
	"notification/internal/domain/layout/ports"
	"notification/pkg/logger"
)

type Getter struct {
	repo   ports.LayoutRepository
	logger logger.Logger
}

func NewGetter(repo ports.LayoutRepository, logger logger.Logger) *Getter {
	return &Getter{repo: repo, logger: logger}
}

func (uc *Getter) Execute(ctx context.Context, id string) (*common.Response, error) {
	l, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return common.ToResponse(l), nil
}
