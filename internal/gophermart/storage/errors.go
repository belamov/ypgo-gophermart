package storage

import "fmt"

type NotUniqueError struct {
	Err   error
	Field string
}

func (err *NotUniqueError) Error() string {
	return fmt.Sprintf("not unique field: %s", err.Field)
}

func (err *NotUniqueError) Unwrap() error {
	return err.Err
}

func NewNotUniqueError(field string, err error) error {
	return &NotUniqueError{
		Err:   err,
		Field: field,
	}
}
