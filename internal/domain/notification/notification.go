package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Notification struct {
	ID string `json:"id" db:"id"`
	// @Schema(example="a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", format="uuid", description="ID макета уведомления")
	LayoutID  string         `json:"layout_id" db:"layout_id"`
	Status    Status         `json:"status" db:"status"`
	Title     string         `json:"title" db:"title"`
	Content   string         `json:"content" db:"content"`
	Data      JSONB          `json:"data" db:"data"`         // тип для JSON
	Channels  pq.StringArray `json:"channels" db:"channels"` // тип для массива
	SentAt    sql.NullTime   `json:"sent_at" db:"sent_at"`   // Nullable поле
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	Receiver  string         `json:"receiver" db:"receiver"`
}

// Для работы с JSONB в Postgres
type JSONB map[string]interface{}

// Scan реализация интерфейса sql.Scanner
func (m *JSONB) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, m)
}

// Status определяет статус уведомления
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)
