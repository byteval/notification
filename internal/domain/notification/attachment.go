package notification

import (
	"time"
)

type Attachment struct {
	ID             string    `json:"id" db:"id"`
	NotificationID string    `json:"notification_id" db:"notification_id"`
	FileName       string    `json:"file_name" db:"file_name"`
	OriginalName   string    `json:"original_name" db:"original_name"`
	ContentType    string    `json:"content_type" db:"content_type"`
	Size           int64     `json:"size" db:"size"`
	FilePath       string    `json:"file_path" db:"file_path"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
