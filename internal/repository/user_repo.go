package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-asset-service/internal/models"
)

// UserRepository отвечает за операции с таблицей пользователей.
type UserRepository struct {
	db *pgxpool.Pool // Пул соединений с базой данных
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// FindByLogin находит пользователя по логину.
// Если пользователь найден, возвращает указатель на объект модели User, иначе — ошибку.
func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, created_at 
		FROM users 
		WHERE login = $1`, login)

	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser создает нового пользователя в базе данных.
// Пример реализации, которую можно доработать в зависимости от требований.
func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (login, password_hash, created_at)
		VALUES ($1, $2, $3)`,
		u.Login, u.PasswordHash, u.CreatedAt)
	return err
}

// GetUserByID возвращает пользователя по его ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE id = $1`, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
