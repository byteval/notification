package layouts

import (
	"net/http"

	"notification/internal/app/layout/delete"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// DeleteLayoutHandler обработчик для удаления шаблона
// @Summary Удалить шаблон уведомления
// @Description Удаляет шаблон уведомления по ID
// @Tags Layouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID шаблона"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /layouts/{id} [delete]
func NewDeleteLayoutHandler(deleter *delete.Deleter, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		if err := deleter.Execute(c.Request.Context(), id); err != nil {
			log.Error("Failed to delete layout", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
