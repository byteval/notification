package notifications

import (
	"net/http"

	getUseCase "notification/internal/app/notification/get"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

type GetNotificationHandler struct {
	getter *getUseCase.Getter
	logger logger.Logger
}

func NewGetNotificationHandler(getter *getUseCase.Getter, log logger.Logger) gin.HandlerFunc {
	handler := &GetNotificationHandler{
		getter: getter,
		logger: log,
	}
	return handler.GetNotification
}

// GetNotification обрабатывает GET /notifications/:id
// @Summary Получить уведомление
// @Description Возвращает уведомление по ID
// @Tags Notifications
// @Produce json
// @Param id path string true "ID уведомления"
// @Success 200 {object} getUseCase.Response
// @Router /notifications/{id} [get]
func (h *GetNotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	response, err := h.getter.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
