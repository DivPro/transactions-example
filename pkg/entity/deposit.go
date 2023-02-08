package entity

import "time"

type Deposit struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Amount string `json:"amount"`
}

type DepositView struct {
	Deposit
	CreatedAt time.Time `json:"created_at"`
}
