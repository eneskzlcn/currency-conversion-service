package model

import "time"

type UserWallet struct {
	UserID    int
	Currency  string
	Balance   float32
	CreatedAt time.Time
	UpdatedAt time.Time
}
