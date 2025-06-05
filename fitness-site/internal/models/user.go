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
	IsAdmin      bool

	Age       sql.NullInt64
	HeightCM  sql.NullInt64
	WeightKG  sql.NullFloat64
	Goals     sql.NullString
	AvatarURL sql.NullString
}

var ErrUserNotFound = errors.New("user not found")

func GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
func GetByID(ctx context.Context, id int) (*User, error) {
	return nil, ErrUserNotFound
}

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
