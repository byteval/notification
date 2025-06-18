package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Notification struct {
	ID        string             `json:"id" db:"id"`
	LayoutID  string             `json:"layout_id" db:"layout_id"`
	Status    Status             `json:"status" db:"status"`
	Title     string             `json:"title" db:"title"`
	Content   string             `json:"content" db:"content"`
	Data      MapStringInterface `json:"data" db:"data"`         // тип для JSON
	Channels  StringArray        `json:"channels" db:"channels"` // тип для массива
	SentAt    sql.NullTime       `json:"sent_at" db:"sent_at"`   // Nullable поле
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	Receiver  string             `json:"receiver" db:"receiver"`
}

// MapStringInterface для работы с JSONB в Postgres
type MapStringInterface map[string]interface{}

// Scan реализация интерфейса sql.Scanner
func (m *MapStringInterface) Scan(value interface{}) error {
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

// StringArray для работы с массивами в Postgres
type StringArray []string

// Status определяет статус уведомления
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)

func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusSent, StatusDelivered, StatusFailed:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

// Scan реализация интерфейса sql.Scanner
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}
