package validation

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var ErrorMessages = map[string]string{
	"required":             "Поле обязательно для заполнения",
	"email":                "Некорректный email",
	"max":                  "Превышена максимальная длина",
	"notification_channel": "Недопустимый канал уведомления",
}

func HandleValidationErrors(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	errors := make([]FieldError, 0, len(errs))
	for _, fe := range errs {
		errors = append(errors, FieldError{
			Field:   fe.Field(),
			Message: getErrorMessage(fe),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
}

func getErrorMessage(fe validator.FieldError) string {
	if msg, exists := ErrorMessages[fe.Tag()]; exists {
		return msg
	}
	return fe.Error()
}
