package main

import (
	"context"
	"net/http"
	"notification/config"
	"notification/internal/container"
	api "notification/internal/transport/http"
	"notification/pkg/logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// @title Notification Service API
// @version 1.0
// @description API сервиса уведомлений
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// Инициализация логгера
	log := logger.New(os.Stdout)

	// Инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", "error", err)
	}

	// Инициализация DI-контейнера
	ctn, err := container.Build(*cfg)
	if err != nil {
		log.Fatal("Ошибка инициализации контейнера:", "error", err)
	}

	// Инициализация HTTP сервера
	srv := initHttpServer(ctn, *cfg)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск сервера в goroutine
	log.Info("Запуск сервера", "port", cfg.HTTP.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Ошибка запуска сервера", "error", err)
			os.Exit(1) // Выход при ошибке запуска
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Info("Получен сигнал завершения работы")

	// Graceful shutdown с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Ошибка graceful shutdown сервера", "error", err)
	} else {
		log.Info("Сервер успешно остановлен")
	}
}

func initHttpServer(cnt *container.Container, cfg config.Config) *http.Server {
	// Инициализация роутеров
	r := api.SetupRouter(cnt)

	return &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: r,
		// Дополнительные настройки сервера
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  15 * time.Second,
	}
}
