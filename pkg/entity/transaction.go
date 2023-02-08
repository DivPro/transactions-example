package entity

import "time"

type Transaction struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
	Amount   string `json:"amount"`
}

type TransactionView struct {
	Transaction
	CreatedAt time.Time `json:"created_at"`
}
