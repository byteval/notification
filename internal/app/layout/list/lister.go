// internal/app/layout/list/lister.go
package list

import (
	"context"

	"notification/internal/domain/layout/ports"
	"notification/pkg/logger"
)

type Lister struct {
	repo   ports.LayoutRepository
	logger logger.Logger
}

func NewLister(repo ports.LayoutRepository, logger logger.Logger) *Lister {
	return &Lister{repo: repo, logger: logger}
}

func (uc *Lister) Execute(ctx context.Context) ([]*Response, int, error) {
	layouts, err := uc.repo.List(ctx)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return ToResponseList(layouts), total, nil
}
