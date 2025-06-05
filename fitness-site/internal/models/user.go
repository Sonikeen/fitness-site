// internal/models/user.go
package models

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string

	Age       sql.NullInt64
	HeightCM  sql.NullInt64
	WeightKG  sql.NullFloat64
	Goals     sql.NullString
	AvatarURL sql.NullString
}

var ErrUserNotFound = errors.New("user not found")

// GetByEmail читает пользователя без created_at/updated_at:
func GetByEmail(ctx context.Context, email string) (*User, error) {
	row := DB.QueryRowContext(ctx, `
SELECT id, username, email, password_hash,
       age, height_cm, weight_kg, goals, avatar_url
  FROM users WHERE email=$1`, email)

	var u User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Age, &u.HeightCM, &u.WeightKG, &u.Goals, &u.AvatarURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetByID без created_at/updated_at:
func GetByID(ctx context.Context, id int) (*User, error) {
	row := DB.QueryRowContext(ctx, `
SELECT id, username, email, password_hash,
       age, height_cm, weight_kg, goals, avatar_url
  FROM users WHERE id=$1`, id)

	var u User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Age, &u.HeightCM, &u.WeightKG, &u.Goals, &u.AvatarURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// UpdateProfile обновляет поля без updated_at:
func UpdateProfile(ctx context.Context, u *User) error {
	_, err := DB.ExecContext(ctx, `
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
	return err
}
