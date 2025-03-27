package models

import "time"

// User представляет пользователя системы.
// ID — уникальный идентификатор пользователя.
// Login — логин пользователя.
// PasswordHash — хеш пароля (не выводится в JSON, чтобы не раскрывать пароль).
// CreatedAt — время создания записи о пользователе.
type User struct {
	ID           int64     `json:"id"`         // Уникальный идентификатор пользователя
	Login        string    `json:"login"`      // Логин пользователя
	PasswordHash string    `json:"-"`          // Хеш пароля (не сериализуется в JSON)
	CreatedAt    time.Time `json:"created_at"` // Время создания пользователя
}
