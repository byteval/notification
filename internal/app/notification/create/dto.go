package create

// Request - модель запроса для создания уведомления
// @Description Модель запроса для создания уведомления
type Request struct {
	LayoutID    string                 `form:"layout_id" validate:"required,uuid4" example:"5471ced2-46ce-4a9f-955e-363b3b11db87"`
	Title       string                 `form:"title" validate:"required"`
	Data        map[string]interface{} `form:"data"` // JSON-массив
	Emails      []string               `form:"receiver" validate:"required,min=1,dive,email" example:"[\"mail@example.com\"]"`
	Attachments []Attachment           `form:"-" swaggerignore:"true"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	FilePath    string `json:"-"`
}
