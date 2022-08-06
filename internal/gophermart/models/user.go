package models

type User struct {
	ID             int
	Login          string
	HashedPassword string
}
