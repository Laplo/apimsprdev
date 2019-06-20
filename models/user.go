package models

import "github.com/satori/go.uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password", db:"password"`
	Email    string    `json:"email", db:"email"`
	IsAdmin  bool      `json:"isAdmin"`
}
