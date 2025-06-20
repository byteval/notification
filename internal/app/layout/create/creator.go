package create

import (
	"context"
	"fmt"
	"time"

	"notification/internal/domain/layout/ports"
	"notification/pkg/logger"
)

// Creator обрабатывает логику создания шаблонов
type Creator struct {
	repo   ports.LayoutRepository
	logger logger.Logger
}

func NewCreator(
	repo ports.LayoutRepository,
	logger logger.Logger,
) *Creator {
	return &Creator{
		repo:   repo,
		logger: logger,
	}
}

// Создание нового шаблона
func (s *Creator) Create(ctx context.Context, req Request) (*Response, error) {
	n, err := req.ToDomain()
	if err != nil {
		s.logger.Error("Failed to convert request to domain", "error", err)
		return nil, fmt.Errorf("failed to convert request to domain: %w", err)
	}

	n.CreatedAt = time.Now()

	created, err := s.repo.Create(ctx, n)

	if err != nil {
		s.logger.Error("Failed to save layout", "error", err)
		return nil, fmt.Errorf("failed to save layout: %w", err)
	}

	return ToResponse(created), nil
}
