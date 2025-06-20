package ports

import (
	"context"
	"notification/internal/domain/layout"
)

// Интерфейс для работы с хранилищем шаблонов уведомлений
type LayoutRepository interface {
	// Create создает новый шаблон уведомления
	Create(ctx context.Context, layout *layout.Layout) (*layout.Layout, error)

	// GetByID возвращает шаблон по ID
	GetByID(ctx context.Context, id string) (*layout.Layout, error)

	// List возвращает список шаблонов с пагинацией
	List(ctx context.Context) ([]*layout.Layout, error)

	// Update обновляет существующий шаблон
	Update(ctx context.Context, layout *layout.Layout) (*layout.Layout, error)

	// Delete удаляет шаблон по ID
	Delete(ctx context.Context, id string) error

	// Количество записей
	Count(ctx context.Context) (int, error)
}
