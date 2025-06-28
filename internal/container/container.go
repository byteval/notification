package container

import (
	"fmt"
	"os"

	"notification/config"
	notificationCreate "notification/internal/app/notification/create"
	notificationGet "notification/internal/app/notification/get"
	"notification/internal/app/notification/mailprocessor"

	layoutCreate "notification/internal/app/layout/create"
	layoutDelete "notification/internal/app/layout/delete"
	layoutGet "notification/internal/app/layout/get"
	layoutList "notification/internal/app/layout/list"
	layoutUpdate "notification/internal/app/layout/update"

	"notification/internal/infrastructure/imap"
	"notification/internal/infrastructure/notifiers/smtp"
	"notification/internal/infrastructure/postgres"
	"notification/pkg/logger"
)

type Container struct {
	Config              config.Config
	Logger              logger.Logger
	DB                  *postgres.Database
	NotificationCreator *notificationCreate.NotificationCreator
	NotificationGetter  *notificationGet.NotificationGetter
	LayoutGetter        *layoutGet.Getter
	LayoutLister        *layoutList.Lister
	LayoutCreator       *layoutCreate.Creator
	LayoutUpdater       *layoutUpdate.Updater
	LayoutDeleter       *layoutDelete.Deleter
	MailProcessor       *mailprocessor.ProcessIncomingEmail
}

func Build(cfg config.Config) (*Container, error) {
	log := logger.New(os.Stdout)

	db, err := postgres.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	repo := postgres.NewNotificationRepository(db.GetDB())
	layoutRepo := postgres.NewLayoutRepository(db.GetDB())

	notifier := notificationCreate.NewNotifier(repo, layoutRepo, smtp.NewSender(cfg.SMTP, log), log, cfg.WorkerPool.Workers)

	return &Container{
		Config:              cfg,
		Logger:              log,
		DB:                  db,
		NotificationCreator: notificationCreate.NewCreator(repo, notifier, log),
		NotificationGetter:  notificationGet.NewGetter(repo),
		LayoutGetter:        layoutGet.NewGetter(layoutRepo, log),
		LayoutLister:        layoutList.NewLister(layoutRepo, log),
		LayoutCreator:       layoutCreate.NewCreator(layoutRepo, log),
		LayoutUpdater:       layoutUpdate.NewUpdater(layoutRepo, log),
		LayoutDeleter:       layoutDelete.NewDeleter(layoutRepo, log),
		MailProcessor:       mailprocessor.NewProcessIncomingEmail(imap.NewIMAPClient(cfg.IMAP, log)),
	}, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
