package service

import (
	"context"
	"errors"
	"time"

	"go-asset-service/internal/models"
	"go-asset-service/internal/repository"
	"go-asset-service/pkg/utils"
)

// AuthService реализует бизнес-логику аутентификации пользователя.
type AuthService struct {
	userRepo    *repository.UserRepository    // Репозиторий для поиска пользователей
	sessionRepo *repository.SessionRepository // Репозиторий для работы с сессиями

	sessionTTL time.Duration // Максимальное время жизни сессии (например, 24 часа)
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(u *repository.UserRepository, s *repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    u,
		sessionRepo: s,
		sessionTTL:  24 * time.Hour, // Ограничение 24 часа для пользовательской сессии
	}
}

// Login осуществляет аутентификацию пользователя.
// Принимает логин, пароль и IP-адрес клиента. Если аутентификация успешна,
// удаляются предыдущие сессии пользователя (чтобы оставалась только одна активная),
// генерируется новый session ID, создается новая сессия и возвращается session ID.
func (as *AuthService) Login(ctx context.Context, login, password, ip string) (string, error) {
	// Поиск пользователя по логину
	user, err := as.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return "", errors.New("invalid login/password")
	}

	// Проверка пароля: хешируем входной пароль и сравниваем с сохраненным в базе
	hashed := utils.Md5Hash(password)
	if user.PasswordHash != hashed {
		return "", errors.New("invalid login/password")
	}

	// Удаляем предыдущие сессии пользователя, чтобы сохранить только одну активную сессию
	err = as.sessionRepo.DeleteByUID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// Генерируем новый session ID
	sessionID, err := utils.GenerateToken(16)
	if err != nil {
		return "", err
	}

	// Создаем новую сессию с текущим временем и IP-адресом клиента
	sess := &models.Session{
		ID:        sessionID,
		UID:       user.ID,
		IPAddress: ip,
		CreatedAt: time.Now(),
	}
	err = as.sessionRepo.Create(ctx, sess)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// ValidateToken проверяет, существует ли сессия с данным session ID,
// и не просрочена ли она (срок жизни не превышает sessionTTL).
func (as *AuthService) ValidateToken(ctx context.Context, sessionID string) (*models.Session, error) {
	sess, err := as.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Если сессия просрочена, возвращаем ошибку
	if time.Since(sess.CreatedAt) > as.sessionTTL {
		return nil, errors.New("session expired")
	}

	return sess, nil
}
