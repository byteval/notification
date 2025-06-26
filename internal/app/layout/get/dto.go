package get

// Request определяет входные данные для получения шаблона
// swagger:parameters getLayout
type Request struct {
	// ID шаблона
	// in: path
	// required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" validate:"required,uuid4"`
}
