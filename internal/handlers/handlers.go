package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-asset-service/internal/repository"
	"go-asset-service/internal/service"
)

// RegisterRoutes регистрирует все HTTP-маршруты API.
func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	// Создаем репозитории для работы с пользователями, сессиями и файлами.
	userRepo := repository.NewUserRepository(pool)
	sessionRepo := repository.NewSessionRepository(pool)
	assetRepo := repository.NewAssetRepository(pool)

	// Инициализируем сервис авторизации.
	authSrv := service.NewAuthService(userRepo, sessionRepo)

	// Создаем хендлеры для авторизации и работы с файлами.
	authHandler := NewAuthHandler(userRepo, sessionRepo)
	assetHandler := NewAssetHandler(assetRepo, authSrv)

	// Эндпоинт авторизации: POST /api/auth.
	mux.HandleFunc("/api/auth", authHandler.Login)

	// Эндпоинт загрузки файла: POST /api/upload-asset/{assetName}.
	mux.HandleFunc("/api/upload-asset/", assetHandler.UploadAsset)

	// Эндпоинт для получения (GET) и удаления (DELETE) файла: /api/asset/{assetName}.
	mux.HandleFunc("/api/asset/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			assetHandler.GetAsset(w, r)
		case http.MethodDelete:
			assetHandler.DeleteAsset(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// Эндпоинт для получения списка файлов: GET /api/assets.
	mux.HandleFunc("/api/assets", assetHandler.ListAssets)

	// Эндпоинт healthcheck: GET /health.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
