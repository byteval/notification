package layouts

import (
	"net/http"

	"notification/internal/app/layout/update"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateLayoutHandler обработчик для обновления шаблона
// @Summary Обновить шаблон уведомления
// @Description Обновляет существующий шаблон уведомления
// @Tags Layouts
// @Accept  json
// @Produce  json
// @Param id path string true "ID шаблона"
// @Param input body update.Request true "Данные для обновления"
// @Success 200 {object} update.Response
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /layouts/{id} [put]
func NewUpdateLayoutHandler(updater *update.Updater, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		var req update.Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		req.ID = id

		resp, err := updater.Execute(c.Request.Context(), req)
		if err != nil {
			log.Error("Failed to update layout", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
