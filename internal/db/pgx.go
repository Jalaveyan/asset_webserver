package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-asset-service/internal/config"
)

// Connect устанавливает подключение к базе данных PostgreSQL,
// используя параметры из конфигурации, и возвращает пул соединений.
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	// Формируем DSN для подключения к базе данных
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	// Создаем пул соединений
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to create pgxpool: %v", err)
		return nil, err
	}

	// Проверяем соединение с базой данных
	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("[ERROR] Cannot ping DB: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("[INFO] Successfully connected to DB")
	return pool, nil
}
