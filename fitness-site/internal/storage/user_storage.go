package storage

import (
	"context"
	"fmt"
    "github.com/jackc/pgx/v5/pgxpool"
	"fitness-site/internal/models"
)

// UserStorage описывает операции с таблицей users.
type UserStorage interface {
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// UserPGStorage — реализация UserStorage поверх db.Pool.
type UserPGStorage struct {
	pool *pgxpool.Pool
}

// NewUserStorage создаёт новое хранилище пользователей.
func NewUserStorage(pool *pgxpool.Pool) *UserPGStorage {
    return &UserPGStorage{pool: pool}
}


// Create сохраняет пользователя u (PasswordHash должен быть уже хешем).
func (s *UserPGStorage) Create(ctx context.Context, u *models.User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`,
		u.Name, u.Email, u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("UserPGStorage.Create: %w", err)
	}
	return nil
}

// GetByEmail возвращает пользователя по email.
func (s *UserPGStorage) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash FROM users WHERE email = $1`, email)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash); err != nil {
		return nil, fmt.Errorf("UserPGStorage.GetByEmail: %w", err)
	}
	return &u, nil
}
