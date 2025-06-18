package container

import (
	"fmt"
	"os"

	"notification/config"
	createUseCase "notification/internal/app/notification/create"
	"notification/internal/infrastructure/notifiers/smtp"
	"notification/internal/infrastructure/postgres"
	"notification/pkg/logger"
)

type Container struct {
	Config  config.Config
	Logger  logger.Logger
	DB      *postgres.Database
	Service *createUseCase.Creator
}

func Build(cfg config.Config) (*Container, error) {
	log := logger.New(os.Stdout)

	db, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	repo := postgres.NewNotificationRepository(db.GetDB())
	smtpSender := smtp.NewSender(cfg, log)
	notifier := createUseCase.NewNotifier(repo, smtpSender, log)

	notificationService := createUseCase.NewNotificationService(
		repo,
		notifier,
		log,
	)

	return &Container{
		Config:  cfg,
		Logger:  log,
		DB:      db,
		Service: notificationService,
	}, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
