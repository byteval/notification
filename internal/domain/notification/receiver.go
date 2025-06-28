package notification

import (
	"time"
)

// Email представляет получателя и статус отправки для конкретного email
type NotificationReceiver struct {
	ID             string     `json:"id" db:"id"`
	NotificationID string     `json:"notification_id" db:"notification_id"`
	Email          string     `json:"email" db:"email"`
	Status         Status     `json:"status" db:"status"`
	Error          *string    `json:"error,omitempty" db:"error"`
	SentAt         *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
}

// Status определяет статус уведомления
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)
