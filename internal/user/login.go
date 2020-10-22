package user

import (
	"context"
	"database/sql"

	"reddit-clone-example/storage"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login sign in user into system by his name and pass. It returns
// token on success or message on failure.
func Login(ctx context.Context, name, pass string) (token, message string) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "login")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return "", "can't login"
	}
	defer tx.Rollback()

	// Check for the same login existence.
	var id uuid.UUID
	{
		const q = `SELECT id, pass FROM account WHERE login = $1`
		var h []byte
		if err = tx.QueryRowxContext(ctx, q, name).Scan(&id, &h); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return "", "internal error"
		}
		if err == sql.ErrNoRows {
			return "", "user not found"
		}
		if err = bcrypt.CompareHashAndPassword(h, []byte(pass)); err != nil {
			return "", "wrong password"
		}
	}

	// Generate JWT token.
	{
		if token, err = createToken(id, name); err != nil {
			log.Log("err", err, "desc", "can't create token")
			return "", "can't register"
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return "", "can't register"
	}

	return token, ""
}
