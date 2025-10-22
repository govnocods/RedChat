package models

import "time"

type User struct {
	Id        int       `json:"Id"`
	Username  string    `json:"Username"`
	Password  string    `json:"Password"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type Message struct {
	Id        int    `json:"Id"`
	SenderId  int    `json:"SenderId"`
	Content   any    `json:"Content"`
	CreatedAt string `json:"CreatedAt"`
}
