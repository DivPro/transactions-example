package entity

import "time"

type User struct {
	UserID     string    `json:"user_id"`
	FistName   string    `json:"fist_name"`
	LastName   string    `json:"last_name"`
	SecondName string    `json:"second_name"`
	CreatedAt  time.Time `json:"created_at"`
	Balance    string    `json:"balance"`
}
