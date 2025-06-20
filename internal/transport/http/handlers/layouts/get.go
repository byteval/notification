package layouts

import (
	"net/http"

	"notification/internal/app/layout/get"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetLayoutHandler обработчик для получения шаблона
// @Summary Получить шаблон уведомления
// @Description Возвращает шаблон уведомления по ID
// @Tags Layouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID шаблона"
// @Success 200 {object} get.Response
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /layouts/{id} [get]
func NewGetLayoutHandler(getter *get.Getter, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		resp, err := getter.Execute(c.Request.Context(), id)
		if err != nil {
			log.Error("Failed to get layout", "error", err, "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "layout not found"})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
