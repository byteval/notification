package postgres

import (
	"context"
	"fmt"

	"notification/internal/domain/layout"

	"github.com/jmoiron/sqlx"
)

// @see internal\domain\layout\repository.go
type LayoutRepository struct {
	db *sqlx.DB
}

func NewLayoutRepository(db *sqlx.DB) *LayoutRepository {
	return &LayoutRepository{db: db}
}

func (r *LayoutRepository) Create(ctx context.Context, l *layout.Layout) (*layout.Layout, error) {
	query := `
		INSERT INTO layouts (
			name, 
			description, 
			subject, 
			body, 
			type, 
			variables, 
			is_active, 
			version
		) VALUES (
			:name, 
			:description, 
			:subject, 
			:body, 
			:type, 
			:variables, 
			:is_active, 
			1
		) RETURNING 
		    id,
			created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, l)
	if err != nil {
		return nil, fmt.Errorf("failed to create layout: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after insert")
	}

	if err := rows.StructScan(l); err != nil {
		return nil, fmt.Errorf("failed to scan created layout: %w", err)
	}

	return l, nil
}

func (r *LayoutRepository) GetByID(ctx context.Context, id string) (*layout.Layout, error) {
	var l layout.Layout
	query := `SELECT * FROM layouts WHERE id = $1`

	err := r.db.GetContext(ctx, &l, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get layout by id: %w", err)
	}

	return &l, nil
}

func (r *LayoutRepository) List(ctx context.Context) ([]*layout.Layout, error) {
	var layouts []*layout.Layout

	query := `SELECT * FROM layouts WHERE 1=1`

	err := r.db.SelectContext(ctx, &layouts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list layouts: %w", err)
	}

	return layouts, nil
}

func (r *LayoutRepository) Update(ctx context.Context, l *layout.Layout) (*layout.Layout, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Оптимистичная блокировка
	query := `
		UPDATE layouts SET
			name = :name,
			description = :description,
			subject = :subject,
			body = :body,
			type = :type,
			variables = :variables,
			is_active = :is_active,
			updated_at = NOW(),
			version = version + 1
		WHERE 
			id = :id
		RETURNING *`

	result, err := tx.NamedExecContext(ctx, query, l)
	if err != nil {
		return nil, fmt.Errorf("failed to update layout: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Обновляем версию в возвращаемом объекте
	l.Version++

	return l, nil
}

func (r *LayoutRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM layouts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete layout: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return layout.ErrLayoutNotFound
	}

	return nil
}

func (r *LayoutRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM layouts`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count layouts: %w", err)
	}

	return count, nil
}
