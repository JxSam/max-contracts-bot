package model

type User struct {
	ChatID       int64   `json:"chat_id"`
	Name         *string `json:"name"`
	Verification bool    `json:"verification"`
}

type UserDefault struct {
	ChatID int64 `json:"chat_id"`
}
