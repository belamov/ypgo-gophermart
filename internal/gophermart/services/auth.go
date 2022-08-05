package services

import (
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Register(credentials models.Credentials) (models.User, error)
	GenerateToken(user models.User) (string, error)
}

type Auth struct {
	userRepo storage.Users
}

func (a *Auth) Register(credentials models.Credentials) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user, err := a.userRepo.CreateNew(credentials.Login, string(hash))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *Auth) GenerateToken(user models.User) (string, error) {
	// TODO: implement
	return "", nil
}
