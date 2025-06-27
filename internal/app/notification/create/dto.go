package create

import "time"

// Request - модель запроса для создания уведомления
// @Description Модель запроса для создания уведомления
type Request struct {
	LayoutID string                 `json:"layout_id" validate:"required"`
	Title    string                 `json:"title" validate:"required"`
	Data     map[string]interface{} `json:"data"`
	Emails   []string               `json:"receiver" validate:"required"`
}

// Response - модель успешного ответа
// @name Response
// @Description Стандартный ответ API
type Response struct {
	ID        string             `json:"id"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	Receivers []ReceiverResponse `json:"receivers,omitempty"`
}

// ReceiverResponse - модель ответа для получателей
type ReceiverResponse struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}
