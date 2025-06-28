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

	// 2. Сохраняем получателей
	for i := range n.NotificationReceivers {
		receiver := &n.NotificationReceivers[i]
		receiver.NotificationID = n.ID // Устанавливаем связь

		err = tx.QueryRowContext(ctx,
			`INSERT INTO notification_receivers 
             (notification_id, email, status)
             VALUES ($1, $2, $3) RETURNING id`,
			receiver.NotificationID,
			receiver.Email,
			receiver.Status,
		).Scan(&receiver.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create receiver: %w", err)
		}
	}

	// 3. Сохраняем вложения
	for i := range n.Attachments {
		attachment := &n.Attachments[i]
		attachment.NotificationID = n.ID // Устанавливаем связь

		err = tx.QueryRowContext(ctx,
			`INSERT INTO notification_attachments 
             (notification_id, file_name, original_name, content_type, size, file_path)
             VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
			attachment.NotificationID,
			attachment.FileName,
			attachment.OriginalName,
			attachment.ContentType,
			attachment.Size,
			attachment.FilePath,
		).Scan(&attachment.ID, &attachment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
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
		UPDATE notification_receivers
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
