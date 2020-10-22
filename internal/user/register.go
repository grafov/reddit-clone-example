package user

import (
	"context"
	"database/sql"
	"time"

	"reddit-clone-example/internal/config"
	"reddit-clone-example/internal/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx context.Context, name, pass string) (token, message string) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "register")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return "", "can't register"
	}
	defer tx.Rollback()

	// Check for the same login existence.
	{
		const q = `SELECT id, login FROM account WHERE login = $1`
		var id uuid.UUID
		if err = tx.QueryRowxContext(ctx, q, name).Scan(&id); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return "", "internal error"
		}
		if id != uuid.Nil {
			return "", "user exists"
		}
	}

	// Create a new user.
	var id = uuid.New()
	{
		const q = `INSERT INTO account (id, login, pass) VALUES ($1, $2, $3)`
		// cost choice refs to
		// https://security.stackexchange.com/questions/17207/recommended-of-rounds-for-bcrypt/83382#83382
		var h []byte
		if h, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost); err != nil {
			l.Log("err", err, "sql", q, "desc", "password hash generation failed")
			return "", "internal error"
		}
		if _, err = tx.ExecContext(ctx, q, id, name, h); err != nil {
			l.Log("err", err, "sql", q, "desc", "new user creation failed")
			return "", "internal error"
		}
	}

	// Generate JWT token.
	if token, err = createToken(id, name); err != nil {
		log.Log("err", err, "desc", "can't create token")
		return "", "can't register"
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return "", "can't register"
	}

	return token, ""
}

func createToken(id uuid.UUID, name string) (string, error) {
	n := time.Now().Unix()
	a := Authbox{
		jwt.StandardClaims{
			IssuedAt:  n,
			ExpiresAt: n + int64(config.App.TokenDuration*60),
		},
		User{id, name, time.Now()},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, a.claims())
	return t.SignedString([]byte(config.App.TokenSecret))
}
