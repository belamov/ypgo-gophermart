package services

import (
	"errors"

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
		return models.User{}, NewRegistrationError(credentials, err)
	}

	user, err := a.userRepo.CreateNew(credentials.Login, string(hash))

	var notUniqueError *storage.NotUniqueError
	if errors.As(err, &notUniqueError) {
		return models.User{}, NewLoginTakenError(credentials.Login, err)
	}

	if err != nil {
		return models.User{}, NewRegistrationError(credentials, err)
	}

	return user, nil
}

func (a *Auth) GenerateToken(user models.User) (string, error) {
	// TODO: implement
	return "", nil
}
