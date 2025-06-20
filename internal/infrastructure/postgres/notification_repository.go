package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"notification/internal/domain/notification"
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
			layout_id, status, title, content, data, channels, receiver, type
		) VALUES (
			:layout_id, :status, :title, :content, :data, :channels, :receiver, :type
		) RETURNING id, created_at`

	_, err := r.db.NamedExecContext(ctx, query, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*notification.Notification, error) {
	query := `
		SELECT id, user_id, title, message, type, status, created_at
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
