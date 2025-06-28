package create

import "time"

// Request - модель запроса для создания уведомления
// @Description Модель запроса для создания уведомления
type Request struct {
	LayoutID    string                 `form:"layout_id" validate:"required,uuid4" example:"5471ced2-46ce-4a9f-955e-363b3b11db87"`
	Title       string                 `form:"title" validate:"required"`
	Data        map[string]interface{} `form:"data"` // JSON-массив
	Emails      []string               `form:"receiver" validate:"required,dive,email" example:"[\"mail@example.com\"]"`
	Attachments []Attachment           `form:"-" swaggerignore:"true"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	FilePath    string `json:"-"`
}

// Response - модель успешного ответа
// @name Response
// @Description Стандартный ответ API
type Response struct {
	ID          string               `json:"id"`
	Status      string               `json:"status"`
	CreatedAt   time.Time            `json:"created_at"`
	Receivers   []ReceiverResponse   `json:"receivers,omitempty"`
	Attachments []AttachmentResponse `json:"attachments,omitempty"`
}

// ReceiverResponse - модель ответа для получателей
type ReceiverResponse struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

type AttachmentResponse struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
}
