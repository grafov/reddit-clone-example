package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var log = kiwi.Fork().With("pkg", "user")

type (
	Authbox struct {
		jwt.StandardClaims
		User User `json:"user"`
	}

	User struct {
		ID        uuid.UUID `json:"id" db:"id"`
		Login     string    `json:"username" db:"login"`
		CreatedAt time.Time `json:"-" db:"created_at"`
	}
)

// Claims makes the structure for JWT signing.
func (a *Authbox) claims() *jwt.MapClaims {
	return &jwt.MapClaims{"user": a.User, "iat": a.IssuedAt, "exp": a.ExpiresAt}
}
