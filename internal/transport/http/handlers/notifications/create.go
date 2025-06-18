package notifications

import (
	"net/http"

	createUseCase "notification/internal/app/notification/create"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

type CreateNotificationHandler struct {
	service *createUseCase.Creator
	logger  logger.Logger
}

func NewCreateNotificationHandler(creator *createUseCase.Creator, log logger.Logger) gin.HandlerFunc {
	handler := &CreateNotificationHandler{
		service: creator,
		logger:  log,
	}
	return handler.CreateNotification
}

// CreateNotification обрабатывает POST /notifications
// @Summary Создать уведомление
// @Description Создает новое уведомление
// @Tags Notifications
// @Accept json
// @Produce json
// @Param input body createUseCase.Request true "Данные уведомления"
// @Success 201 {object} createUseCase.Response
// @Router /notifications [post]
func (h *CreateNotificationHandler) CreateNotification(c *gin.Context) {
	var req createUseCase.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание уведомления через сервис
	response, err := h.service.CreateNotification(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create notification", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
