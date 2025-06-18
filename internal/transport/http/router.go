package http

import (
	"notification/internal/container"
	handlers "notification/internal/transport/http/handlers/notifications"
	"notification/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"

	//_ "notification/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Notification Service API
// @version 1.0
// @description API сервиса уведомлений
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func SetupRouter(cnt *container.Container) *gin.Engine {
	router := gin.Default()

	// Глобальные middleware
	router.Use(
		middleware.Logging(cnt.Logger),
		middleware.Recovery(cnt.Logger),
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
			notificationsGroup.POST("", handlers.NewCreateNotificationHandler(cnt.Service, cnt.Logger))
			// notificationsGroup.GET("/:id", api.NewGetNotificationHandler(cnt.Getter, cnt.Logger))
		}
	}

	// Health check вне API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
