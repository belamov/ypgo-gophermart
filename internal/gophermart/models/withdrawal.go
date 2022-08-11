package models

import "time"

type Withdrawal struct {
	OrderID          int
	UserID           int
	WithdrawalAmount float64
	CreatedAt        time.Time
}
