package service

import (
	"errors"
	"fitness-site/internal/models"
)

var mockUsers = []model.User{}

func Register(name, email, password string) error {
	for _, u := range mockUsers {
		if u.Email == email {
			return errors.New("user already exists")
		}
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: password, // позже — хешировать!
	}

	mockUsers = append(mockUsers, user)
	return nil
}

func Login(email, password string) (*model.User, error) {
	for _, u := range mockUsers {
		if u.Email == email && u.Password == password {
			return &u, nil
		}
	}
	return nil, errors.New("invalid credentials")
}
