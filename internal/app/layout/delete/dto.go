package delete

// Request определяет входные данные для удаления шаблона
// swagger:parameters deleteLayout
type Request struct {
	// ID шаблона
	// in: path
	// required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" validate:"required,uuid4"`
}

// Response определяет ответ после удаления шаблона
// swagger:model DeleteLayoutResponse
type Response struct {
	// Статус операции
	// example: success
	Status string `json:"status"`
}
