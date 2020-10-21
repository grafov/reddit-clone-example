package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type (
	Authbox struct {
		jwt.StandardClaims
		User User `json:"user"`
	}
	User struct {
		ID    uuid.UUID `json:"id"`
		Login string    `json:"username"`
	}
)

// Claims makes the structure for JWT signing.
func (a *Authbox) claims() *jwt.MapClaims {
	return &jwt.MapClaims{"user": a.User, "iat": a.IssuedAt, "exp": a.ExpiresAt}
}
