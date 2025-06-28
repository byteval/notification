package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"notification/config"

	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер PostgreSQL
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

// Cоздаем новое подключение к PostgreSQL
func New(cfg config.Database) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)       // Максимум открытых соединений
	db.SetMaxIdleConns(cfg.MaxIdleConns)       // Максимум бездействующих соединений
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime) // Максимальное время жизни соединения

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

// GetDB возвращает объект *sql.DB для прямого доступа
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}

// Close закрывает соединение с БД
func (d *Database) Close() error {
	return d.db.Close()
}

// WithTransaction выполняет операцию в транзакции
func (d *Database) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// IsNotFound проверяет, является ли ошибка "не найдено"
func (d *Database) IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
