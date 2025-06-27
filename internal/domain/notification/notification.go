package notification

import (
	"encoding/json"
	"errors"
	"time"
)

type Notification struct {
	ID                    string                 `json:"id" db:"id"`
	LayoutID              string                 `json:"layout_id" db:"layout_id"`
	Title                 string                 `json:"title" db:"title"`
	Data                  JSONB                  `json:"data" db:"data"`
	NotificationReceivers []NotificationReceiver `json:"emails" db:"-"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
}

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
