package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"go-asset-service/internal/repository"
	"go-asset-service/internal/service"
)

// AuthHandler отвечает за обработку запросов к эндпоинту аутентификации (/api/auth)
type AuthHandler struct {
	authService *service.AuthService // Сервис для авторизации пользователей
}

// NewAuthHandler создаёт новый экземпляр AuthHandler, инициализируя сервис авторизации.
func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(userRepo, sessionRepo),
	}
}

// loginRequest описывает структуру JSON-запроса для аутентификации.
type loginRequest struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// loginResponse описывает структуру JSON-ответа при успешной аутентификации.
type loginResponse struct {
	Token string `json:"token"` // Авторизационный токен (session-id)
}

// Login обрабатывает POST /api/auth.
// Он читает JSON-запрос с логином и паролем, получает IP-адрес клиента,
// вызывает сервис авторизации и возвращает токен, либо ошибку.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Логирование входящего запроса для отладки
	log.Printf("[INFO] /api/auth called from %s", r.RemoteAddr)

	var req loginRequest

	// Попытка декодировать JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Failed to decode JSON in /api/auth: %v", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Получаем IP-адрес клиента (используется для записи в сессию)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Вызываем сервис авторизации: передаём логин, пароль и IP-адрес
	token, err := h.authService.Login(context.Background(), req.Login, req.Password, ip)
	if err != nil {
		// Логирование ошибки авторизации (например, неверный логин/пароль)
		log.Printf("[WARN] Failed login for user=%s ip=%s err=%v", req.Login, ip, err)
		http.Error(w, `{"error":"invalid login/password"}`, http.StatusUnauthorized)
		return
	}

	// Логирование успешной авторизации
	log.Printf("[INFO] User logged in: login=%s ip=%s token=%s", req.Login, ip, token)

	// Формирование ответа с токеном в формате JSON
	resp := loginResponse{Token: token}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
