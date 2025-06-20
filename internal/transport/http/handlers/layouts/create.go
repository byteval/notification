package layouts

import (
	"net/http"

	"notification/internal/app/layout/create"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CreateLayoutHandler обработчик для создания шаблона
// @Summary Создать новый шаблон уведомления
// @Description Создает новый шаблон уведомления
// @Tags Layouts
// @Accept  json
// @Produce  json
// @Param input body create.Request true "Данные шаблона"
// @Success 201 {object} create.Response
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /layouts [post]
func NewCreateLayoutHandler(creator *create.Creator, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req create.Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		resp, err := creator.Create(c.Request.Context(), req)
		if err != nil {
			log.Error("Failed to create layout", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}
