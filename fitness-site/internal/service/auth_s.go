package service

import (
	"errors"
	"fitness-site/internal/models"
)

var mockUsers = []models.User{}

func Register(name, email, password string) error {
	for _, u := range mockUsers {
		if u.Email == email {
			return errors.New("user already exists")
		}
	}

	user := models.User{
		Name:     name,
		Email:    email,
		PasswordHash: password, // позже — хешировать!
	}

	mockUsers = append(mockUsers, user)
	return nil
}

func Login(email, password string) (*models.User, error) {
	for _, u := range mockUsers {
		if u.Email == email && u.PasswordHash == password {
			return &u, nil
		}
	}
	return nil, errors.New("invalid credentials")
}
