package services

import (
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
)

type Authenticator interface {
	Register(credentials models.Credentials) (models.User, error)
	GenerateToken(user models.User) (string, error)
}

type Auth struct{}

func (a *Auth) Register(credentials models.Credentials) (models.User, error) {
	// TODO: implement
	return models.User{}, nil
}

func (a *Auth) GenerateToken(user models.User) (string, error) {
	// TODO: implement
	return "", nil
}
