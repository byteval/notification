package postgres

import (
	"context"
	"fmt"
	"notification/internal/domain/notification"

	"github.com/jmoiron/sqlx"
)

// @see internal\domain\notification\repository.go
type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) (*notification.Notification, error) {
	query := `
        INSERT INTO notifications (
            layout_id, status, title, content, data, channels, receiver
        ) VALUES (
            :layout_id, :status, :title, :content, :data, :channels, :receiver
        ) RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, n)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after insert")
	}

	if err := rows.StructScan(n); err != nil {
		return nil, fmt.Errorf("failed to scan created notification: %w", err)
	}

	return n, nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*notification.Notification, error) {
	query := `
		SELECT *
		FROM notifications 
		WHERE id = $1
	`

	var n notification.Notification
	err := r.db.GetContext(ctx, &n, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}

	return &n, nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status notification.Status) error {
	query := `
		UPDATE notifications
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
