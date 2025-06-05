package service

import (
	"context"
	"errors"
	"fitness-site/internal/models"
	"fitness-site/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store storage.UserStorage
}

func NewUserService(store storage.UserStorage) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Register(ctx context.Context, u *models.User) error {
	if u.Name == "" || u.Email == "" || u.PasswordHash == "" {
		return errors.New("name, email, password required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashed)
	return s.store.Create(ctx, u)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	u, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.store.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, u *models.User) error {
	return s.store.Update(ctx, u)
}
