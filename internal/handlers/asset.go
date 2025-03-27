package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go-asset-service/internal/models"
	"go-asset-service/internal/repository"
	"go-asset-service/internal/service"
)

// AssetHandler реализует HTTP-обработчики для работы с файлами (assets)
type AssetHandler struct {
	assetRepo   *repository.AssetRepository // Репозиторий для работы с данными файлов
	authService *service.AuthService        // Сервис авторизации для проверки токена
}

// NewAssetHandler создает новый экземпляр AssetHandler
func NewAssetHandler(assetRepo *repository.AssetRepository, auth *service.AuthService) *AssetHandler {
	return &AssetHandler{
		assetRepo:   assetRepo,
		authService: auth,
	}
}

// UploadAsset обрабатывает запрос POST /api/upload-asset/{assetName}.
// Он проверяет авторизацию, извлекает имя файла из URL, читает тело запроса
// и сохраняет данные в базе данных.
func (h *AssetHandler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации через заголовок Authorization: Bearer <token>
	userSession, err := h.checkAuth(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized upload attempt from ip=%s err=%v", r.RemoteAddr, err)
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Извлечение имени файла из URL (последний сегмент пути)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("[ERROR] Bad request (missing asset name), user=%d ip=%s", userSession.UID, r.RemoteAddr)
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	assetName := parts[len(parts)-1]

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read upload body: user=%d ip=%s err=%v", userSession.UID, r.RemoteAddr, err)
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}

	// Формирование объекта Asset для сохранения в БД
	asset := &models.Asset{
		Name:      assetName,
		UID:       userSession.UID,
		Data:      data,
		CreatedAt: time.Now(),
	}

	// Сохранение файла (assets) в базе данных через репозиторий
	err = h.assetRepo.CreateAsset(context.Background(), asset)
	if err != nil {
		log.Printf("[ERROR] DB error in CreateAsset: user=%d name=%s ip=%s err=%v", userSession.UID, assetName, r.RemoteAddr, err)
		http.Error(w, `{"error":"failed to save asset"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Asset uploaded successfully: name=%s user=%d ip=%s", assetName, userSession.UID, r.RemoteAddr)
	// Возвращаем успешный ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// GetAsset обрабатывает запрос GET /api/asset/{assetName}.
// Проверяет авторизацию, извлекает имя файла из URL и возвращает содержимое файла.
func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	userSession, err := h.checkAuth(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized get-asset attempt from ip=%s err=%v", r.RemoteAddr, err)
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Извлекаем имя файла из URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("[ERROR] Bad request (missing asset name), user=%d ip=%s", userSession.UID, r.RemoteAddr)
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	assetName := parts[len(parts)-1]

	// Получаем файл из базы данных
	asset, err := h.assetRepo.GetAsset(context.Background(), assetName, userSession.UID)
	if err != nil {
		log.Printf("[WARN] Asset not found: name=%s user=%d ip=%s err=%v", assetName, userSession.UID, r.RemoteAddr, err)
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	log.Printf("[INFO] Asset retrieved: name=%s user=%d ip=%s", assetName, userSession.UID, r.RemoteAddr)
	// Отдаем содержимое файла (raw data)
	w.WriteHeader(http.StatusOK)
	w.Write(asset.Data)
}

// ListAssets обрабатывает запрос GET /api/assets.
// Возвращает список файлов, загруженных текущим пользователем.
func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	userSession, err := h.checkAuth(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized list-assets attempt from ip=%s err=%v", r.RemoteAddr, err)
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Получаем список файлов из базы
	assets, err := h.assetRepo.ListAssets(context.Background(), userSession.UID)
	if err != nil {
		log.Printf("[ERROR] Failed to list assets for user=%d: %v", userSession.UID, err)
		http.Error(w, `{"error":"failed to list assets"}`, http.StatusInternalServerError)
		return
	}

	// Сериализуем список в JSON
	resp, err := json.Marshal(map[string]interface{}{
		"assets": assets,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to marshal assets: %v", err)
		http.Error(w, `{"error":"failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// DeleteAsset обрабатывает запрос DELETE /api/asset/{assetName}.
// Удаляет файл, принадлежащий текущему пользователю.
func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	// Допустим, данный обработчик вызывается только для DELETE-запросов
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Проверка авторизации
	userSession, err := h.checkAuth(r)
	if err != nil {
		log.Printf("[WARN] Unauthorized delete-asset attempt from ip=%s err=%v", r.RemoteAddr, err)
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Извлечение имени файла из URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	assetName := parts[len(parts)-1]

	// Удаляем файл из базы данных
	err = h.assetRepo.DeleteAsset(context.Background(), assetName, userSession.UID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete asset: name=%s user=%d ip=%s err=%v", assetName, userSession.UID, r.RemoteAddr, err)
		http.Error(w, `{"error":"failed to delete asset"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Asset deleted: name=%s user=%d ip=%s", assetName, userSession.UID, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// checkAuth проверяет наличие и валидность Bearer-токена в заголовке Authorization.
// Если токен отсутствует или недействителен, возвращает ошибку.
func (h *AssetHandler) checkAuth(r *http.Request) (*models.Session, error) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return nil, http.ErrNoCookie
	}
	token := strings.TrimPrefix(auth, prefix)
	return h.authService.ValidateToken(context.Background(), token)
}
