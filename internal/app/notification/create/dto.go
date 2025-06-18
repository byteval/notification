package create

import "time"

// Request - модель запроса для создания уведомления
// @Description Модель запроса для создания уведомления
type Request struct {
	LayoutID string                 `json:"layout_id" validate:"required"`
	Title    string                 `json:"title" validate:"required"`
	Content  string                 `json:"content" validate:"required"`
	Data     map[string]interface{} `json:"data"`
	Channels []string               `json:"channels" validate:"required,min=1"`
	Receiver string                 `json:"receiver" validate:"required"`
}

// Response - модель успешного ответа
// @name Response
// @Description Стандартный ответ API
type Response struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
