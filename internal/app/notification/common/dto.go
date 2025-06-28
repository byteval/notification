package common

import (
	"time"
)

// NotificationResponse - модель успешного ответа
// @name Response
// @Description Стандартный ответ API
type Response struct {
	ID          string               `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	Receivers   []ReceiverResponse   `json:"receivers,omitempty"`
	Attachments []AttachmentResponse `json:"attachments,omitempty"`
}

// ReceiverResponse - модель ответа для получателей
type ReceiverResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type AttachmentResponse struct {
	ID          string `json:"id"`
	FileName    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}
