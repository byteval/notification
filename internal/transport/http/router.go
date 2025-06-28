package http

import (
	"notification/internal/container"
	layoutHandlers "notification/internal/transport/http/handlers/layouts"
	handlers "notification/internal/transport/http/handlers/notifications"
	"notification/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"

	_ "notification/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(ctn *container.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Глобальные middleware
	router.Use(
		middleware.Logging(ctn.Logger),
		middleware.Recovery(ctn.Logger),
		middleware.CORS(),
	)

	// Swagger handler
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа API v1
	apiV1 := router.Group("/api/v1")
	{
		// Группа уведомлений
		notificationsGroup := apiV1.Group("/notifications")
		{
			notificationsGroup.POST(
				"",
				handlers.NewCreateNotificationHandler(
					ctn.NotificationCreator,
					ctn.Logger,
					ctn.Config.HTTP.UploadDir,
					ctn.Config.HTTP.MaxUploadSize,
				),
			)
			notificationsGroup.GET("/:id", handlers.NewGetNotificationHandler(ctn.NotificationGetter, ctn.Logger))
		}

		// Группа шаблонов уведомлений
		layoutsGroup := apiV1.Group("/layouts")
		{
			layoutsGroup.POST("", layoutHandlers.NewCreateLayoutHandler(ctn.LayoutCreator, ctn.Logger))
			layoutsGroup.GET("/:id", layoutHandlers.NewGetLayoutHandler(ctn.LayoutGetter, ctn.Logger))
			layoutsGroup.PUT("/:id", layoutHandlers.NewUpdateLayoutHandler(ctn.LayoutUpdater, ctn.Logger))
			layoutsGroup.DELETE("/:id", layoutHandlers.NewDeleteLayoutHandler(ctn.LayoutDeleter, ctn.Logger))
			layoutsGroup.GET("", layoutHandlers.NewListLayoutsHandler(ctn.LayoutLister, ctn.Logger))
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
