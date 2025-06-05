package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"fitness-site/internal/models"
)

type UserStorage interface {
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
}

type UserPGStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserPGStorage {
	return &UserPGStorage{pool: pool}
}

func (s *UserPGStorage) Create(ctx context.Context, u *models.User) error {
	_, err := s.pool.Exec(ctx, `
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)`,
		u.Name, u.Email, u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}
	return nil // <--- вот так!
}

func (s *UserPGStorage) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `
SELECT id, username, email, password_hash,
       age, height_cm, weight_kg, goals, avatar_url
  FROM users WHERE email=$1`, email)

	var u models.User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Age, &u.HeightCM, &u.WeightKG, &u.Goals, &u.AvatarURL,
	); err != nil {
		return nil, fmt.Errorf("GetByEmail: %w", err)
	}
	return &u, nil
}

func (s *UserPGStorage) GetByID(ctx context.Context, id int) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `
SELECT id, username, email, password_hash,
       age, height_cm, weight_kg, goals, avatar_url
  FROM users WHERE id=$1`, id)

	var u models.User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Age, &u.HeightCM, &u.WeightKG, &u.Goals, &u.AvatarURL,
	); err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return &u, nil
}

func (s *UserPGStorage) Update(ctx context.Context, u *models.User) error {
	_, err := s.pool.Exec(ctx, `
UPDATE users SET
  username   = $1,
  age        = $2,
  height_cm  = $3,
  weight_kg  = $4,
  goals      = $5,
  avatar_url = $6
WHERE id = $7`,
		u.Name, u.Age, u.HeightCM, u.WeightKG, u.Goals, u.AvatarURL, u.ID,
	)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	return nil
}
