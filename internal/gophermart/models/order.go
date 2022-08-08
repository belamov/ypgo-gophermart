package models

type Order struct {
	ID        int `json:"id"`
	CreatedBy int `json:"user_id"`
}
