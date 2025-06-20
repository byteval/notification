package container

import (
	"fmt"
	"os"

	"notification/config"
	createNotification "notification/internal/app/notification/create"
	getNotification "notification/internal/app/notification/get"

	createLayout "notification/internal/app/layout/create"
	deleteLayout "notification/internal/app/layout/delete"
	getLayout "notification/internal/app/layout/get"
	listLayout "notification/internal/app/layout/list"
	updateLayout "notification/internal/app/layout/update"

	"notification/internal/infrastructure/notifiers/smtp"
	"notification/internal/infrastructure/postgres"
	"notification/pkg/logger"
)

type Container struct {
	Config              config.Config
	Logger              logger.Logger
	DB                  *postgres.Database
	NotificationCreator *createNotification.Creator
	NotificationGetter  *getNotification.Getter
	LayoutGetter        *getLayout.Getter
	LayoutLister        *listLayout.Lister
	LayoutCreator       *createLayout.Creator
	LayoutUpdater       *updateLayout.Updater
	LayoutDeleter       *deleteLayout.Deleter
}

func Build(cfg config.Config) (*Container, error) {
	log := logger.New(os.Stdout)

	db, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	repo := postgres.NewNotificationRepository(db.GetDB())
	layoutRepo := postgres.NewLayoutRepository(db.GetDB())

	smtpSender := smtp.NewSender(cfg, log)
	notifier := createNotification.NewNotifier(repo, smtpSender, log)

	// Layout
	layoutGetter := getLayout.NewGetter(layoutRepo, log)
	layoutLister := listLayout.NewLister(layoutRepo, log)
	layoutCreator := createLayout.NewCreator(layoutRepo, log)
	layoutUpdater := updateLayout.NewUpdater(layoutRepo, log)
	layoutDeleter := deleteLayout.NewDeleter(layoutRepo, log)

	return &Container{
		Config:              cfg,
		Logger:              log,
		DB:                  db,
		NotificationCreator: createNotification.NewCreator(repo, notifier, log),
		NotificationGetter:  getNotification.NewGetter(repo),
		LayoutGetter:        layoutGetter,
		LayoutLister:        layoutLister,
		LayoutCreator:       layoutCreator,
		LayoutUpdater:       layoutUpdater,
		LayoutDeleter:       layoutDeleter,
	}, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
