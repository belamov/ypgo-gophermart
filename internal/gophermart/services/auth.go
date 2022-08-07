package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(credentials models.Credentials) (models.User, error)
	GenerateToken(user models.User) (string, error)
	Login(credentials models.Credentials) (models.User, error)
	AuthMiddleware() func(h http.Handler) http.Handler
	GetUserID(r *http.Request) int
}

const UserIdClaim = "user_id"

type JWTAuth struct {
	UserRepo  storage.Users
	tokenAuth *jwtauth.JWTAuth
}

func (a *JWTAuth) GetUserID(r *http.Request) int {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims[UserIdClaim].(int)
	return userID
}

func NewAuth(repo storage.Users, secret string) *JWTAuth {
	jwtAuth := jwtauth.New("HS256", []byte(secret), nil)
	return &JWTAuth{
		UserRepo:  repo,
		tokenAuth: jwtAuth,
	}
}

func (a *JWTAuth) Register(credentials models.Credentials) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, NewRegistrationError(credentials, err)
	}

	user, err := a.UserRepo.CreateNew(credentials.Login, string(hash))

	var notUniqueError *storage.NotUniqueError
	if errors.As(err, &notUniqueError) {
		return models.User{}, NewLoginTakenError(credentials.Login, err)
	}

	if err != nil {
		return models.User{}, NewRegistrationError(credentials, err)
	}

	return user, nil
}

func (a *JWTAuth) Login(credentials models.Credentials) (models.User, error) {
	user, err := a.UserRepo.FindByLogin(credentials.Login)
	if err != nil {
		return models.User{}, err
	}

	if user.ID == 0 {
		return models.User{}, NewInvalidCredentialsError(credentials, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(credentials.Password)); err != nil {
		return models.User{}, NewInvalidCredentialsError(credentials, err)
	}

	return user, nil
}

func (a *JWTAuth) GenerateToken(user models.User) (string, error) {
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

func (a *JWTAuth) AuthMiddleware() func(h http.Handler) http.Handler {
	verifier := jwtauth.Verifier(a.tokenAuth)
	authenticator := jwtauth.Authenticator

	return func(h http.Handler) http.Handler {
		return authenticator(verifier(h))
	}
}

func (a *JWTAuth) getTokenClaims(user models.User) (map[string]interface{}, error) {
	claims := map[string]interface{}{}

	jwtauth.SetIssuedNow(claims)

	duration, err := time.ParseDuration("10h")
	if err != nil {
		return nil, err
	}
	jwtauth.SetExpiryIn(claims, duration)

	if user.ID == 0 {
		return nil, errors.New("user id is required")
	}
	claims[UserIdClaim] = user.ID

	return claims, nil
}
