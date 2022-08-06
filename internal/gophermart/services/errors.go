package services

import (
	"fmt"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
)

type RegistrationError struct {
	Err         error
	Credentials models.Credentials
}

func (err *RegistrationError) Error() string {
	return fmt.Sprintf("can't register user with credentials: %+v\n%s", err.Credentials, err.Err.Error())
}

func (err *RegistrationError) Unwrap() error {
	return err.Err
}

func NewRegistrationError(credentials models.Credentials, err error) error {
	return &RegistrationError{
		Err:         err,
		Credentials: credentials,
	}
}

type LoginTakenError struct {
	Err   error
	Login string
}

func (err *LoginTakenError) Error() string {
	return fmt.Sprintf("login is already taken: %s", err.Login)
}

func (err *LoginTakenError) Unwrap() error {
	return err.Err
}

func NewLoginTakenError(login string, err error) error {
	return &LoginTakenError{
		Err:   err,
		Login: login,
	}
}
