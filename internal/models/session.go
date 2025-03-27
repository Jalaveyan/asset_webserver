package models

import "time"

// Session представляет пользовательскую сессию.
// Поле ID — уникальный идентификатор сессии (например, session token).
// Поле UID — идентификатор пользователя, которому принадлежит сессия.
// Поле IPAddress содержит IP-адрес, с которого пользователь прошёл авторизацию.
// Поле CreatedAt фиксирует время создания сессии.
type Session struct {
	ID        string    `json:"id"`         // Уникальный идентификатор сессии
	UID       int64     `json:"uid"`        // Идентификатор пользователя
	IPAddress string    `json:"ip_address"` // IP-адрес клиента
	CreatedAt time.Time `json:"created_at"` // Время создания сессии
}
