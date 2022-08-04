package models

import (
	"errors"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c Credentials) Validate() error {
	if c.Login == "" {
		return errors.New("login required")
	}

	if c.Password == "" {
		return errors.New("password required")
	}

	return nil
}
