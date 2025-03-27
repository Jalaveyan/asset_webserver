package main

import (
	"context"
	"go-asset-service/internal/config"   // Чтение конфигурации из переменных окружения или .env файла
	"go-asset-service/internal/db"       // Подключение к базе данных через pgx
	"go-asset-service/internal/handlers" // Регистрация HTTP-обработчиков (роутов)
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Подключаемся к базе данных
	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)
	}
	defer pool.Close()

	// Создаем HTTP-маршрутизатор и регистрируем маршруты API
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, pool)

	// Настраиваем HTTP-сервер с таймаутами
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Запускаем HTTPS-сервер в отдельной горутине
	go func() {
		log.Printf("Starting HTTPS server on port %s", cfg.AppPort)
		if err := server.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServeTLS error: %v", err)
		}
	}()

	// Создаем канал для получения сигналов завершения (например, Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown: даем 10 секунд на завершение активных запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
