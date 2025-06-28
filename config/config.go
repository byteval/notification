package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	HTTP         HTTP
	WorkerPool   WorkerPool
	Database     Database
	SMTP         SMTP
	IMAP         IMAP
	Websocket    Websocket
	WebHook      WebHook
	ReadInterval time.Duration `env:"READ_INTERVAL" envDefault:"5m"`
}

type HTTP struct {
	Port          int           `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout   time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout  time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
	UploadDir     string        `env:"UPLOAD_DIR" envDefault:"upload"`
	MaxUploadSize int64         `env:"MAX_UPLOAD_SIZE" envDefault:"10485760"`
}

type WorkerPool struct {
	Workers int `env:"WORKER_POOL_WORKERS" envDefault:"50"`
}

type Database struct {
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"postgres"`
	Password        string        `env:"DB_PASSWORD" envDefault:"postgres"`
	Name            string        `env:"DB_NAME" envDefault:"notifications"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"60m"`
	ConnectTimeout  time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"5s"`
}

type SMTP struct {
	Host     string `env:"SMTP_HOST" envDefault:"smtp.example.com"`
	Port     int    `env:"SMTP_PORT" envDefault:"1025"`
	Username string `env:"SMTP_USERNAME" envDefault:""`
	Password string `env:"SMTP_PASSWORD" envDefault:""`
	From     string `env:"SMTP_FROM" envDefault:"noreply@example.com"`
	BCC      string `env:"SMTP_BCC" envDefault:""`
}

type IMAP struct {
	Host     string        `env:"IMAP_HOST" envDefault:"imap.example.com"`
	Port     int           `env:"IMAP_PORT" envDefault:"993"`
	UseTLS   bool          `env:"IMAP_USE_TLS" envDefault:"true"`
	Username string        `env:"IMAP_USERNAME" envDefault:""`
	Password string        `env:"IMAP_PASSWORD" envDefault:""`
	Timeout  time.Duration `env:"IMAP_TIMEOUT" envDefault:"5m"`
}

type Websocket struct {
	Endpoint string        `env:"WEBSOCKET_ENDPOINT" envDefault:"/ws"`
	Timeout  time.Duration `env:"WEBSOCKET_TIMEOUT" envDefault:"30s"`
}

type WebHook struct {
	Endpoint string        `env:"WEBHOOK_ENDPOINT" envDefault:""`
	Timeout  time.Duration `env:"WEBHOOK_TIMEOUT" envDefault:"30s"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// TODO Валидация конфига

	return &cfg, nil
}
