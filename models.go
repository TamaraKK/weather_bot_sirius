package main

type User struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chat_id"`
	City      string `json:"city"`
	Frequency string `json:"frequency"`
}
