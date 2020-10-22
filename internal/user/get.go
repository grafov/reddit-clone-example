package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Get gets user info from a storage by ID. For internal usage only,
// not for HTTP API.
func Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (User, error) {
	var (
		u   User
		l   = log.Fork().With("fn", "get", "user", id)
		err error
	)
	const q = `SELECT id, login, created_at FROM account WHERE id = $1`
	if err = tx.QueryRowxContext(ctx, q, id).StructScan(&u); err != nil && err != sql.ErrNoRows {
		l.Log("err", err, "desc", "db select failed")
		return User{}, err
	}
	if err == sql.ErrNoRows {
		return User{}, err
	}

	return u, nil
}
