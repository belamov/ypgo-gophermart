package models

import "time"

type Withdrawal struct {
	CreatedAt        time.Time
	OrderID          int
	UserID           int
	WithdrawalAmount float64
}
