package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type (
	Authbox struct {
		User     User  `json:"user"`
		IssuedAt int64 `json:"iat"`
		Expired  int64 `json:"exp"`
	}
	User struct {
		ID    uuid.UUID `json:"id"`
		Login string    `json:"username"`
	}
)

func (a *Authbox) claims() *jwt.MapClaims {

	return &jwt.MapClaims{"user": a.User, "iat": a.IssuedAt, "exp": a.Expired}
}
