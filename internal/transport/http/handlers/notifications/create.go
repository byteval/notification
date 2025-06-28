package notifications

import (
	"mime/multipart"
	"net/http"
	"path/filepath"

	createUseCase "notification/internal/app/notification/create"
	"notification/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateNotificationHandler struct {
	service       *createUseCase.Creator
	logger        logger.Logger
	validate      *validator.Validate
	uploadDir     string
	maxUploadSize int64
}

func NewCreateNotificationHandler(
	creator *createUseCase.Creator,
	log logger.Logger,
	uploadDir string,
	maxUploadSize int64,
) gin.HandlerFunc {
	handler := &CreateNotificationHandler{
		service:       creator,
		logger:        log,
		validate:      validator.New(),
		uploadDir:     uploadDir,
		maxUploadSize: maxUploadSize,
	}
	return handler.CreateNotification
}

// @Summary Создание уведомления
// @Description Создает новое уведомление
// @Tags Notifications
// @Accept multipart/form-data
// @Produce json
// @Param layout_id formData string true "ID шаблона (UUID)"
// @Param title formData string true "Заголовок уведомления"
// @Param data formData string false "Дополнительные данные (JSON строка)"
// @Param emails formData []string true "Список email получателей" collectionFormat(multi)
// @Param attachments formData []file false "Файлы вложений"
// @Success 201 {object} createUseCase.Response
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications [post]
func (h *CreateNotificationHandler) CreateNotification(c *gin.Context) {
	var req createUseCase.Request

	// Парсинг основных полей формы
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация запроса
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	// Обработка файловых вложений
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	h.getAttachments(&req, form, c)

	response, err := h.service.CreateNotification(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create notification", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CreateNotificationHandler) getAttachments(req *createUseCase.Request, form *multipart.Form, c *gin.Context) {
	req.Attachments = make([]createUseCase.Attachment, 0)
	files := form.File["attachments"]

	for _, fileHeader := range files {
		if fileHeader.Size > h.maxUploadSize {
			continue
		}

		filePath := filepath.Join(h.uploadDir, uuid.New().String())
		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			h.logger.Error("File save error", "error", err)
			continue
		}

		req.Attachments = append(req.Attachments, createUseCase.Attachment{
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
			FilePath:    filePath,
		})
	}
}
