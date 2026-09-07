package model

import "time"

type Contract struct {
	Id          int64     `json:"id"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	NameProduct string    `json:"name_product"`
	Default     string    `json:"default"`
}
