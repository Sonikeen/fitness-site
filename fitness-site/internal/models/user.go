package models

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User описывает структуру пользователя.
// PasswordHash хранит результат bcrypt.GenerateFromPassword.
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
}

var ErrUserNotFound = errors.New("user not found")

// CreateUser хеширует пароль и сохраняет нового пользователя.
func CreateUser(ctx context.Context, name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = DB.ExecContext(ctx,
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`,
		name, email, string(hashed),
	)
	return err
}

// GetUserByEmail возвращает пользователя по email, либо ErrUserNotFound.
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := DB.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash FROM users WHERE email = $1`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// AuthenticateUser проверяет email и пароль. Если успешно — возвращает *User.
func AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	u, err := GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
