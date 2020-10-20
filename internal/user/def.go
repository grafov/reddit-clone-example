package user

import "github.com/google/uuid"

type Credentials struct {
	User User  `json:"user"`
	Iat  int64 `json:"iat"`
	Exp  int64 `json:"exp"`
}

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
