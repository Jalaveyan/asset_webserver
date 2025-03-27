package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-asset-service/internal/models"
)

// SessionRepository отвечает за выполнение операций с таблицей sessions в базе данных.
type SessionRepository struct {
	db *pgxpool.Pool // Пул соединений с базой данных
}

// NewSessionRepository создает новый экземпляр SessionRepository.
func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create создает новую сессию и сохраняет ее в таблице sessions.
func (r *SessionRepository) Create(ctx context.Context, s *models.Session) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (id, uid, ip_address, created_at)
		 VALUES ($1, $2, $3, $4)`,
		s.ID, s.UID, s.IPAddress, s.CreatedAt,
	)
	return err
}

// FindByID ищет и возвращает сессию по ее уникальному идентификатору (ID).
func (r *SessionRepository) FindByID(ctx context.Context, sessionID string) (*models.Session, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, uid, ip_address, created_at
		 FROM sessions
		 WHERE id = $1`,
		sessionID,
	)

	var s models.Session
	err := row.Scan(&s.ID, &s.UID, &s.IPAddress, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DeleteByUID удаляет все сессии для указанного пользователя (UID).
// Это используется для реализации механизма "единственной активной сессии".
func (r *SessionRepository) DeleteByUID(ctx context.Context, uid int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM sessions WHERE uid = $1`,
		uid,
	)
	return err
}

// DeleteExpired удаляет все сессии, созданные до времени cutoff.
// Этот метод можно использовать для очистки просроченных сессий.
func (r *SessionRepository) DeleteExpired(ctx context.Context, cutoff time.Time) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE created_at < $1`, cutoff)
	return err
}
