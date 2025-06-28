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
	logger := logger.New(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Configuration loading error:", "error", err)
	}

	ctn, err := container.Build(*cfg)
	if err != nil {
		logger.Fatal("Container initialization error:", "error", err)
	}

	srv := initHttpServer(ctn, *cfg)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("Starting the server", "port", cfg.HTTP.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server startup error", "error", err)
			os.Exit(1)
		}
	}()

	ReadingIncomingMail(ctx, ctn, *cfg, logger)

	<-ctx.Done()
	logger.Info("A shutdown signal has been received")

	// Graceful shutdown с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful server shutdown error", "error", err)
	} else {
		logger.Info("The server has been successfully stopped")
	}
}

func initHttpServer(cnt *container.Container, cfg config.Config) *http.Server {
	r := api.SetupRouter(cnt)

	return &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  15 * time.Second,
	}
}

func ReadingIncomingMail(ctx context.Context, ctn *container.Container, cfg config.Config, logger logger.Logger) {
	go func() {
		if cfg.ReadInterval <= 0 {
			logger.Error("The ticker interval must be > 0")
			return
		}

		ticker := time.NewTicker(cfg.ReadInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := ctn.MailProcessor.Execute(ctx); err != nil {
					logger.Error("mail processing failed", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
