package entity

import "time"

type Transaction struct {
	UserID    string     `json:"user_id"`
	TargetID  string     `json:"target_id"`
	Amount    string     `json:"amount"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
