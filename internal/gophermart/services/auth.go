package services

import (
	"errors"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Register(credentials models.Credentials) (models.User, error)
	GenerateToken(user models.User) (string, error)
}

type Auth struct {
	userRepo  storage.Users
	tokenAuth *jwtauth.JWTAuth
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
	claims, err := a.getTokenClaims(user)
	if err != nil {
		return "", err
	}

	_, tokenString, err := a.tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *Auth) getTokenClaims(user models.User) (map[string]interface{}, error) {
	claims := map[string]interface{}{}

	jwtauth.SetIssuedNow(claims)

	duration, err := time.ParseDuration("10h")
	if err != nil {
		return nil, err
	}
	jwtauth.SetExpiryIn(claims, duration)

	if user.ID == "" {
		return nil, errors.New("user id is required")
	}
	claims["user_id"] = user.ID

	return claims, nil
}
