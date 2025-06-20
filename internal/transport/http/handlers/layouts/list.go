package layouts

import (
	"net/http"

	"notification/internal/app/layout/list"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @Summary Получить список шаблонов уведомлений
// @Description Возвращает список шаблонов уведомлений с пагинацией
// @Tags Layouts
// @Accept  json
// @Produce  json
// @Param type query string false "Тип шаблона"
// @Param active query boolean false "Фильтр по активности"
// @Param search query string false "Поисковая строка"
// @Param limit query int false "Лимит (default 10)"
// @Param offset query int false "Смещение (default 0)"
// @Success 200 {object} list.Response
// @Failure 500 {object} map[string]string
// @Router /layouts [get]
func NewListLayoutsHandler(lister *list.Lister, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req list.Request
		if err := c.ShouldBindQuery(&req); err != nil {
			log.Error("Invalid query parameters", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
			return
		}

		resp, total, err := lister.Execute(c.Request.Context())
		if err != nil {
			log.Error("Failed to list layouts", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": resp,
			"total": total,
		})
	}
}
