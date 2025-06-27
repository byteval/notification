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

func (r *NotificationRepository) CreateWithReceivers(
	ctx context.Context,
	n *notification.Notification,
	receivers []notification.NotificationReceiver,
) (*notification.Notification, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Создаем основное уведомление
	err = tx.QueryRowContext(ctx,
		`INSERT INTO notifications (layout_id, title, data, created_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		n.LayoutID, n.Title, n.Data, n.CreatedAt,
	).Scan(&n.ID)
	if err != nil {
		return nil, err
	}

	// Создаем получателей
	for i := range receivers {
		receivers[i].NotificationID = n.ID
		err = tx.QueryRowContext(ctx,
			`INSERT INTO notification_receivers 
			 (notification_id, email, status)
			 VALUES ($1, $2, $3) RETURNING id`,
			receivers[i].NotificationID,
			receivers[i].Email,
			receivers[i].Status,
		).Scan(&receivers[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	n.NotificationReceivers = receivers
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
		UPDATE notification_receivers
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
