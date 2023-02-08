package entity

import "time"

type Deposit struct {
	UserID string `json:"user_id"`
	Amount string `json:"amount"`
}

type DepositView struct {
	ID string `json:"id"`
	Deposit
	CreatedAt time.Time `json:"created_at"`
}
